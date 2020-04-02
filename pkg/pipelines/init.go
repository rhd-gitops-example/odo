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
	DeploymentPath           string
	GitOpsRepo               string
	GitOpsWebhookSecret      string
	Output                   string
	Prefix                   string
	SkipChecks               bool
	ImageRepo                string
	InternalRegistryHostname string
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
	pipelineDir              = "pipelines"
	cicdDir                  = "cicd-environment"
	envsDir                  = "envs"
	baseDir                  = "base"
	kustomize                = "kustomization.yaml"
	namespacesPath           = "01-namespaces/cicd-environment.yaml"
	rolesPath                = "02-rolebindings/pipeline-service-role.yaml"
	rolebindingsPath         = "02-rolebindings/pipeline-service-rolebinding.yaml"
	secretsPath              = "03-secrets/gitops-webhook-secret.yaml"
	tasksPath                = "04-tasks/deploy-from-source-task.yaml"
	ciPipelinesPath          = "05-pipelines/ci-dryrun-from-pr-pipeline.yaml"
	cdPipelinesPath          = "05-pipelines/cd-deploy-from-push-pipeline.yaml"
	appCIPipelinesPath       = "05-pipelines/app-ci-pipeline.yaml"
	appCDPipelinesPath       = "05-pipelines/app-cd-pipeline.yaml"
	prBindingPath            = "06-bindings/github-pr-binding.yaml"
	pushBindingPath          = "06-bindings/github-push-binding.yaml"
	prTemplatePath           = "07-templates/ci-dryrun-from-pr-template.yaml"
	pushTemplatePath         = "07-templates/cd-deploy-from-push-template.yaml"
	appCIBuildPRTemplatePath = "07-templates/app-ci-build-pr-template.yaml"
	appCDBuildPRTemplatePath = "07-templates/app-cd-build-pr-template.yaml"
	eventListenerPath        = "08-eventlisteners/cicd-event-listener.yaml"
	routePath                = "09-routes/gitops-webhook-event-listener.yaml"
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

	isInternalRegistry, imageRepo, err := validatingImageRepo(o)
	if err != nil {
		return err
	}

	namespaces := namespaceNames(o.Prefix)
	gitopsName := getGitopsRepoName(o.GitOpsRepo)
	gitopsPath := filepath.Join(o.Output, gitopsName)

	// check if the gitops dir already exists
	exists, _ := isExisting(gitopsPath)
	if exists {
		return fmt.Errorf("%s already exists at %s", gitopsName, gitopsPath)
	}

	// key: path of the resource
	// value: YAML content of the resource
	outputs := map[string]interface{}{}

	if o.GitOpsWebhookSecret != "" {
		githubSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(namespaces["cicd"], eventlisteners.GitOpsWebhookSecret),
			o.GitOpsWebhookSecret, eventlisteners.WebhookSecretKey)
		if err != nil {
			return fmt.Errorf("failed to generate GitHub Webhook Secret: %w", err)
		}

		outputs[secretsPath] = githubSecret
	}

	// create gitops pipeline
	files := createPipelineResources(outputs, namespaces, o.GitOpsRepo, o.Prefix, o.DeploymentPath, imageRepo, isInternalRegistry)

	pipelinesPath := getPipelinesDir(gitopsPath, o.Prefix)

	fileNames, err := writeResources(pipelinesPath, files)
	if err != nil {
		return err
	}

	sort.Strings(fileNames)
	// kustomize file should refer all the pipeline resources
	if err := addKustomize("resources", fileNames, filepath.Join(pipelinesPath, kustomize)); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{"./pipelines"}, filepath.Join(getCICDDir(gitopsPath, o.Prefix), kustomize)); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{}, filepath.Join(gitopsPath, envsDir, baseDir, kustomize)); err != nil {
		return err
	}

	return nil
}

func getCICDDir(path, prefix string) string {
	return filepath.Join(path, envsDir, addPrefix(prefix, cicdDir))
}

func createPipelineResources(outputs map[string]interface{}, namespaces map[string]string, gitopsRepo, prefix, deploymentPath, imageRepo string, isInternalRegistry bool) map[string]interface{} {

	outputs[namespacesPath] = createNamespace(namespaces["cicd"])

	outputs[rolesPath] = roles.CreateClusterRole(meta.NamespacedName("", roles.ClusterRoleName), rules)

	sa := roles.CreateServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))

	outputs[rolebindingsPath] = roles.CreateRoleBinding(meta.NamespacedName(namespaces["cicd"], roleBindingName), sa, "ClusterRole", roles.ClusterRoleName)

	outputs[tasksPath] = tasks.CreateDeployFromSourceTask(namespaces["cicd"], getPipelinesDir("", prefix))

	outputs[ciPipelinesPath] = createCIPipeline(meta.NamespacedName(namespaces["cicd"], "ci-dryrun-from-pr-pipeline"), namespaces["cicd"])

	outputs[cdPipelinesPath] = createCDPipeline(meta.NamespacedName(namespaces["cicd"], "cd-deploy-from-push-pipeline"), namespaces["cicd"])

	outputs[appCIPipelinesPath] = createDevCIPipeline(meta.NamespacedName(namespaces["cicd"], "app-ci-pipeline"), isInternalRegistry)

	outputs[appCDPipelinesPath] = createDevCDPipeline(meta.NamespacedName(namespaces["cicd"], "app-cd-pipeline"), deploymentPath, namespaces["dev"], isInternalRegistry)

	outputs[prBindingPath] = triggers.CreatePRBinding(namespaces["cicd"])

	outputs[pushBindingPath] = triggers.CreatePushBinding(namespaces["cicd"])

	outputs[prTemplatePath] = triggers.CreateCIDryRunTemplate(namespaces["cicd"], saName)

	outputs[pushTemplatePath] = triggers.CreateCDPushTemplate(namespaces["cicd"], saName)

	outputs[appCIBuildPRTemplatePath] = triggers.CreateDevCIBuildPRTemplate(namespaces["cicd"], saName, imageRepo)

	outputs[appCDBuildPRTemplatePath] = triggers.CreateDevCDDeployTemplate(namespaces["cicd"], saName, imageRepo)

	outputs[eventListenerPath] = eventlisteners.Generate(gitopsRepo, namespaces["cicd"], saName)

	outputs[routePath] = routes.Generate(namespaces["cicd"])

	return outputs
}

func writeResources(path string, files map[string]interface{}) ([]string, error) {
	filenames := make([]string, 0)
	for filename, item := range files {
		err := marshalItemsToFile(filepath.Join(path, filename), list(item))
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func marshalItemsToFile(filename string, items []interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return marshalOutputs(f, items)
}

func list(i interface{}) []interface{} {
	return []interface{}{i}
}

func getPipelinesDir(rootPath, prefix string) string {
	return filepath.Join(rootPath, envsDir, addPrefix(prefix, cicdDir), pipelineDir)
}

func addKustomize(name string, items []string, path string) error {
	content := make([]interface{}, 0)
	content = append(content, map[string]interface{}{name: items})
	return marshalItemsToFile(path, content)
}

func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}

func getGitopsRepoName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func addPrefix(prefix, name string) string {
	if prefix != "" {
		return prefix + name
	}
	return name
}

func isExisting(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

// validatingImageRepo validates the input image repo.  It determines if it is
// for internal registry and prepend internal registry hostname if neccessary.
func validatingImageRepo(o *InitParameters) (bool, string, error) {
	components := strings.Split(o.ImageRepo, "/")

	// repo url has minimum of 2 components
	if len(components) < 2 {
		return false, "", imageRepoValidationErrors(o.ImageRepo)
	}

	for _, v := range components {
		// check for empty components
		if strings.TrimSpace(v) == "" {
			return false, "", imageRepoValidationErrors(o.ImageRepo)
		}
		// check for white spaces
		if len(v) > len(strings.TrimSpace(v)) {
			return false, "", imageRepoValidationErrors(o.ImageRepo)
		}
	}

	if len(components) == 2 {
		if components[0] == "docker.io" || components[0] == "quay.io" {
			// we recognize docker.io and quay.io.  It is missing one component
			return false, "", imageRepoValidationErrors(o.ImageRepo)
		}
		// We have format like <project>/<app> which is an internal registry.
		// We prepend the internal registry hostname.
		return true, o.InternalRegistryHostname + "/" + o.ImageRepo, nil
	}

	// Check the first component to see if it is an internal registry
	if len(components) == 3 {
		return components[0] == o.InternalRegistryHostname, o.ImageRepo, nil
	}

	// > 3 components.  invalid repo
	return false, "", imageRepoValidationErrors(o.ImageRepo)
}

func imageRepoValidationErrors(imageRepo string) error {
	return fmt.Errorf("failed to parse image repo:%s, expected image repository in the form <registry>/<username>/<repository> or <project>/<app> for internal registry", imageRepo)
}
