package pipelines

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
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
		v1rbac.PolicyRule{
			APIGroups: []string{"bitnami.com"},
			Resources: []string{"sealedsecrets"},
			Verbs:     []string{"get", "patch"},
		},
	}
)

const (
	pipelineDir       = "pipelines"
	cicdDir           = "cicd"
	envsDir           = "environments"
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
		installed, err := CheckTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	// check if the gitops dir already exists
	exists, err := isExisting(o.Output)
	if exists {
		return err
	}

	files, err := CreateResources(o.Prefix, o.GitOpsRepo, o.GitOpsWebhookSecret)
	if err != nil {
		return err
	}

	pipelinesPath := getPipelinesDir(o.Output, o.Prefix)

	fileNames, err := writeResources(pipelinesPath, files)
	if err != nil {
		return err
	}

	sort.Strings(fileNames)
	// kustomize file should refer all the pipeline resources
	if err := addKustomize("resources", fileNames, filepath.Join(pipelinesPath, kustomize)); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{"./pipelines"}, filepath.Join(getCICDDir(o.Output, o.Prefix), baseDir, kustomize)); err != nil {
		return err
	}

	// Add overlays
	if err := addKustomize("bases", []string{"../base"}, filepath.Join(getCICDDir(o.Output, o.Prefix), "overlays", kustomize)); err != nil {
		return err
	}

	return nil
}

func getCICDDir(path, prefix string) string {
	return filepath.Join(path, envsDir, addPrefix(prefix, cicdDir))
}

func CreateResources(prefix, gitOpsRepo, gitOpsWebhook string) (map[string]interface{}, error) {

	// key: path of the resource
	// value: YAML content of the resource
	outputs := map[string]interface{}{}
	cicdNamespace := addPrefix(prefix, "cicd")

	githubSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(cicdNamespace, eventlisteners.GitOpsWebhookSecret),
		gitOpsWebhook, eventlisteners.WebhookSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate GitHub Webhook Secret: %w", err)
	}

	outputs[secretsPath] = githubSecret

	outputs[namespacesPath] = CreateNamespace(cicdNamespace)

	outputs[rolesPath] = roles.CreateClusterRole(meta.NamespacedName("", roles.ClusterRoleName), rules)

	sa := roles.CreateServiceAccount(meta.NamespacedName(cicdNamespace, saName))

	outputs[rolebindingsPath] = roles.CreateRoleBinding(meta.NamespacedName(cicdNamespace, roleBindingName), sa, "ClusterRole", roles.ClusterRoleName)

	outputs[tasksPath] = tasks.CreateDeployFromSourceTask(cicdNamespace, getPipelinesDir("", prefix))

	outputs[ciPipelinesPath] = createCIPipeline(meta.NamespacedName(cicdNamespace, "ci-dryrun-from-pr-pipeline"), cicdNamespace)

	outputs[cdPipelinesPath] = createCDPipeline(meta.NamespacedName(cicdNamespace, "cd-deploy-from-push-pipeline"), cicdNamespace)

	outputs[prBindingPath] = triggers.CreatePRBinding(cicdNamespace)

	outputs[pushBindingPath] = triggers.CreatePushBinding(cicdNamespace)

	outputs[prTemplatePath] = triggers.CreateCIDryRunTemplate(cicdNamespace, saName)

	outputs[pushTemplatePath] = triggers.CreateCDPushTemplate(cicdNamespace, saName)

	outputs[eventListenerPath] = eventlisteners.Generate(gitOpsRepo, cicdNamespace, saName)

	outputs[routePath] = routes.Generate(cicdNamespace)
	return outputs, nil
}

func writeResources(path string, files map[string]interface{}) ([]string, error) {
	filenames := make([]string, 0)
	for filename, item := range files {
		err := marshalItemsToFile(filepath.Join(path, filename), item)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func marshalItemsToFile(filename string, item interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return marshalOutput(f, item)
}

func getPipelinesDir(rootPath, prefix string) string {
	return filepath.Join(rootPath, envsDir, addPrefix(prefix, cicdDir), baseDir, pipelineDir)
}

func addKustomize(name string, items []string, path string) error {
	content := []interface{}{}
	content = append(content, map[string]interface{}{name: items})
	return marshalItemsToFile(path, content)
}

func CheckTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}

func GetGitopsRepoName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func addPrefix(prefix, name string) string {
	if prefix != "" {
		return prefix + name
	}
	return name
}

func isExisting(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return true, fmt.Errorf("%q: Dir already exists at %s", filepath.Base(path), path)
	}
	return true, fmt.Errorf("%q: File already exists at %s", filepath.Base(path), path)
}
