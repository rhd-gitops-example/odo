package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/out"
	"github.com/openshift/odo/pkg/pipelines/out/gitrepo"
	"github.com/openshift/odo/pkg/pipelines/roles"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/openshift/odo/pkg/pipelines/tasks"
	"github.com/openshift/odo/pkg/pipelines/triggers"
	v1rbac "k8s.io/api/rbac/v1"
)

// InitParameters is a struct that provides flags for initialise command
type InitParameters struct {
	GitOpsRepo          string
	GitOpsWebhookSecret string
	Output              string
	Prefix              string
	SkipChecks          bool
}

// PolicyRules to be bound to service account
var (
	rules = []v1rbac.PolicyRule{
		v1rbac.PolicyRule{
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"patch"},
		},
		v1rbac.PolicyRule{
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{"clusterroles"},
			Verbs:     []string{"bind", "patch"},
		},
		v1rbac.PolicyRule{
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{"rolebindings"},
			Verbs:     []string{"get", "patch"},
		},
	}
)

const (
	pipelineDir       = "pipelines"
	cicdDir           = "cicd-environment"
	envsDir           = "envs"
	baseDir           = "base"
	kustomize         = "kustomization.yaml"
	namespacesPath    = "01-namespaces/cicd-environment.yaml"
	rolesPath         = "02-rolebindings/pipeline-service-role.yaml"
	rolebindingsPath  = "02-rolebindings/pipeline-service-rolebinding.yaml"
	secretsPath       = "03-secrets/gitops-webhook-secret.yaml"
	tasksPath         = "04-tasks/deploy-from-source-task.yaml"
	ciPipelinesPath   = "05-pipelines/ci-dryrun-from-pr-pipeline.yaml"
	cdPipelinesPath   = "05-pipelines/cd-deploy-from-push-pipeline.yaml"
	prBindingPath     = "06-bindings/github-pr-binding.yaml"
	pushBindingPath   = "06-bindings/github-push-binding.yaml"
	prTemplatePath    = "07-templates/ci-dryrun-from-pr-template.yaml"
	pushTemplatePath  = "07-templates/cd-deploy-from-push-template.yaml"
	eventListenerPath = "08-eventlisteners/cicd-event-listener.yaml"
	routePath         = "09-routes/gitops-webhook-event-listener.yaml"
)

// Init function will initialise the gitops directory
func Init(o *InitParameters) error {

	if !o.SkipChecks {
		installed, err := checkTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	namespaces := namespaceNames(o.Prefix)

	output, err := gitrepo.New("https://github.com/wtam2018/gitopstest.git",
		"wtam2018", "gitopstest", "dev", "051d73dc54a72b2fdece19cf4dace45cf630e9d9")
	if err != nil {
		return fmt.Errorf("failed to create output : %w", err)
	}

	/*
			output, err := fs.New(o.GitOpsRepo, o.Output)
		if err != nil {
			return fmt.Errorf("failed to create output : %w", err)
		}
	*/

	pipelinesPath := getPipelinesDir(o.Prefix)

	if o.GitOpsWebhookSecret != "" {
		githubSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(namespaces["cicd"], eventlisteners.GitOpsWebhookSecret),
			o.GitOpsWebhookSecret, eventlisteners.WebhookSecretKey)
		if err != nil {
			return fmt.Errorf("failed to generate GitHub Webhook Secret: %w", err)
		}

		output.Add(filepath.Join(pipelinesPath, secretsPath), githubSecret)
	}

	// add gitops pipeline resource to output
	addPipelineResources(pipelinesPath, output, namespaces, o.GitOpsRepo, o.Prefix)

	paths := output.GetPaths()
	sort.Strings(paths)

	// kustomize file should refer all the pipeline resources
	addKustomize(output, "resources", paths, filepath.Join(pipelinesPath, kustomize))
	addKustomize(output, "bases", []string{"./pipelines"}, filepath.Join(getCICDDir(o.Prefix), kustomize))
	addKustomize(output, "bases", []string{}, filepath.Join(envsDir, baseDir, kustomize))

	return output.Write()
}

func getCICDDir(prefix string) string {
	return filepath.Join(envsDir, addPrefix(prefix, cicdDir))
}

func addPipelineResources(pipelinePath string, o out.Output, namespaces map[string]string, gitopsRepo, prefix string) {

	o.Add(filepath.Join(pipelinePath, namespacesPath), createNamespace(namespaces["cicd"]))

	o.Add(filepath.Join(pipelinePath, rolesPath), roles.CreateClusterRole(meta.NamespacedName("", roles.ClusterRoleName), rules))

	sa := roles.CreateServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))

	o.Add(filepath.Join(pipelinePath, rolebindingsPath), roles.CreateRoleBinding(meta.NamespacedName(namespaces["cicd"], roleBindingName), sa, "ClusterRole", roles.ClusterRoleName))

	o.Add(filepath.Join(pipelinePath, tasksPath), tasks.CreateDeployFromSourceTask(namespaces["cicd"], getPipelinesDir(prefix)))

	o.Add(filepath.Join(pipelinePath, ciPipelinesPath), createCIPipeline(meta.NamespacedName(namespaces["cicd"], "ci-dryrun-from-pr-pipeline"), namespaces["cicd"]))

	o.Add(filepath.Join(pipelinePath, cdPipelinesPath), createCDPipeline(meta.NamespacedName(namespaces["cicd"], "cd-deploy-from-push-pipeline"), namespaces["cicd"]))

	o.Add(filepath.Join(pipelinePath, prBindingPath), triggers.CreatePRBinding(namespaces["cicd"]))

	o.Add(filepath.Join(pipelinePath, pushBindingPath), triggers.CreatePushBinding(namespaces["cicd"]))

	o.Add(filepath.Join(pipelinePath, prTemplatePath), triggers.CreateCIDryRunTemplate(namespaces["cicd"], saName))

	o.Add(filepath.Join(pipelinePath, pushTemplatePath), triggers.CreateCDPushTemplate(namespaces["cicd"], saName))

	o.Add(filepath.Join(pipelinePath, eventListenerPath), eventlisteners.Generate(gitopsRepo, namespaces["cicd"], saName))

	o.Add(filepath.Join(pipelinePath, routePath), routes.Generate(namespaces["cicd"]))
}

func getPipelinesDir(prefix string) string {
	return filepath.Join(envsDir, addPrefix(prefix, cicdDir), pipelineDir)
}

func addKustomize(o out.Output, name string, items []string, path string) {
	content := make([]interface{}, 0)
	content = append(content, map[string]interface{}{name: items})
	o.Add(path, content)
}

func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}

func addPrefix(prefix, name string) string {
	if prefix != "" {
		return prefix + name
	}
	return name
}
