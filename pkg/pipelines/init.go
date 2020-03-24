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
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/tasks"
	"github.com/openshift/odo/pkg/pipelines/triggers"
	v1rbac "k8s.io/api/rbac/v1"
)

// InitParameters is a struct that provides flags for initialise command
type InitParameters struct {
	GitOpsWebhookSecret string
	GitOpsRepo          string
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
	appsDir           = "apps"
	kustomize         = "kustomization.yaml"
	namespacesPath    = "01-namespaces/cicd-environment.yaml"
	rolesPath         = "02-rolebindings/pipeline-service-role.yaml"
	rolebindingsPath  = "02-rolebindings/pipeline-service-rolebinding.yaml"
	secretsPath       = "03-secrets/github-webhook-secret.yaml"
	tasksPath         = "04-tasks/deploy-from-source-task.yaml"
	ciPipelinesPath   = "05-pipelines/ci-dryrun-from-pr-pipeline.yaml"
	cdPipelinesPath   = "05-pipelines/cd-deploy-from-push-pipeline.yaml"
	prBindingPath     = "06-bindings/github-pr-binding.yaml"
	pushBindingPath   = "06-bindings/github-push-binding.yaml"
	prTemplatePath    = "07-templates/ci-dryrun-from-pr-template.yaml"
	pushTemplatePath  = "07-templates/cd-deploy-from-push-template.yaml"
	eventListenerPath = "08-eventlisteners/cicd-event-listener.yaml"
	routePath         = "09-routes/github-webhook-event-listener.yaml"
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

	gitopsName := getGitopsRepoName(o.GitOpsRepo)
	gitopsPath := filepath.Join(o.Output, gitopsName)

	// check if the gitops dir already exists
	exists, _ := isExisting(gitopsPath)
	if exists {
		return fmt.Errorf("%s already exists at %s", gitopsName, gitopsPath)
	}

	// key: path of the resource
	// value: YAML content of the resource
	outputs := map[string][]interface{}{}

	if o.GitOpsWebhookSecret != "" {
		githubSecret, err := createOpaqueSecret(meta.NamespacedName(namespaces["cicd"], eventlisteners.GithubWebhookSecret), o.GitOpsWebhookSecret, eventlisteners.WebhookSecretKey)
		if err != nil {
			return fmt.Errorf("failed to generate GitHub Webhook Secret: %w", err)
		}

		outputs[secretsPath] = append(outputs[secretsPath], githubSecret)
	}

	// create gitops pipeline
	files := createPipelineResources(outputs, namespaces, o.GitOpsRepo, o.Prefix)

	pipelinesPath := getPipelinesDir(gitopsPath, o.Prefix)

	fileNames, err := writeResources(pipelinesPath, files)
	if err != nil {
		return err
	}

	sort.Strings(fileNames)
	// kustomize file should refer all the pipeline resources
	if err := addKustomize("resources", fileNames, filepath.Join(pipelinesPath, kustomize), o.Prefix); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{"./pipelines"}, filepath.Join(getCICDDir(gitopsPath, o.Prefix), kustomize), o.Prefix); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{}, filepath.Join(gitopsPath, envsDir, baseDir, kustomize), o.Prefix); err != nil {
		return err
	}

	return nil
}

func getCICDDir(path, prefix string) string {
	return filepath.Join(path, envsDir, addPrefix(prefix, cicdDir))
}

func createPipelineResources(outputs map[string][]interface{}, namespaces map[string]string, gitopsRepo, prefix string) map[string][]interface{} {

	// create namespace
	outputs[namespacesPath] = append(outputs[namespacesPath], createNamespace(namespaces["cicd"]))

	// create roles and rolebindings for pipeline service account
	role := createClusterRole(meta.NamespacedName("", clusterRoleName), rules)
	outputs[rolesPath] = append(outputs[rolesPath], role)

	sa := createServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))
	rolebinding := createRoleBinding(meta.NamespacedName(namespaces["cicd"], roleBindingName), sa, role.Kind, role.Name)
	outputs[rolebindingsPath] = append(outputs[rolebindingsPath], rolebinding)

	// create deploy from source task
	task := tasks.GenerateDeployFromSourceTask(namespaces["cicd"], getPipelinesDir("", prefix))
	outputs[tasksPath] = append(outputs[tasksPath], task)

	// create ci-cd pipelines for gitops repo
	ciPipeline := createStageCIPipeline(meta.NamespacedName(namespaces["cicd"], "ci-dryrun-from-pr-pipeline"), namespaces["cicd"])
	outputs[ciPipelinesPath] = append(outputs[ciPipelinesPath], ciPipeline)

	cdPipeline := createStageCDPipeline(meta.NamespacedName(namespaces["cicd"], "cd-deploy-from-push-pipeline"), namespaces["cicd"])
	outputs[cdPipelinesPath] = append(outputs[cdPipelinesPath], cdPipeline)

	// create trigger bindings
	prBinding := triggers.CreatePRBinding(namespaces["cicd"])
	outputs[prBindingPath] = append(outputs[prBindingPath], prBinding)

	pushBinding := triggers.CreatePushBinding(namespaces["cicd"])
	outputs[pushBindingPath] = append(outputs[pushBindingPath], pushBinding)

	// create trigger templates
	prTemplate := triggers.CreateStageCIDryRunTemplate(namespaces["cicd"], saName)
	outputs[prTemplatePath] = append(outputs[prTemplatePath], prTemplate)

	pushTemplate := triggers.CreateStageCDPushTemplate(namespaces["cicd"], saName)
	outputs[pushTemplatePath] = append(outputs[pushTemplatePath], pushTemplate)

	// create eventlisteners and route
	eventListener := eventlisteners.Generate(gitopsRepo, namespaces["cicd"], saName)
	outputs[eventListenerPath] = append(outputs[eventListenerPath], eventListener)

	route := routes.Generate(namespaces["cicd"])
	outputs[routePath] = append(outputs[routePath], route)

	return outputs
}

func writeResources(path string, files map[string][]interface{}) ([]string, error) {
	filenames := make([]string, 0)
	for filename, items := range files {
		err := marshalItemsToFile(filepath.Join(path, filename), items)
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

// generate file name in the format index-prefix-filename.yaml
func fileName(index int, prefix, name string) string {
	return fmt.Sprintf("%02d-%v%v.yaml", index, prefix, name)
}

func getPipelinesDir(rootPath, prefix string) string {
	return filepath.Join(rootPath, envsDir, addPrefix(prefix, cicdDir), pipelineDir)
}

func addKustomize(name string, items []string, path, prefix string) error {
	content := make([]interface{}, 0)
	content = append(content, map[string]interface{}{name: items})
	return marshalItemsToFile(path, content)
}

// create and invoke a Tekton Checker
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
