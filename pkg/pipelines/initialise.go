package pipelines

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/tasks"
	"github.com/openshift/odo/pkg/pipelines/triggers"
	v1rbac "k8s.io/api/rbac/v1"
	"sigs.k8s.io/yaml"
)

// InitialiseParameters is a struct that provides flags for initialise command
type InitialiseParameters struct {
	GitOpsRepo               string
	Output                   string
	AppGithubSecret          string
	AppGitRepo               string
	AppImageRepo             string
	Prefix                   string
	InternalRegistryHostname string
	DockerConfigJSONFileName string
	SkipChecks               bool
}

// PolicyRules to be bound to service account
var (
	rules = []v1rbac.PolicyRule{
		v1rbac.PolicyRule{
			APIGroups: []string{""},
			Resources: []string{"namespace"},
			Verbs:     []string{"patch"},
		},
	}
)

var (
	pipelineDir = "pipelines"
	cicdDir     = "cicd-environment"
	envsDir     = "envs"
	baseDir     = "base"
	appsDir     = "apps"
	kustomize   = "kustomization.yaml"
)

const (
	cicdNamespace     = "cicd-environment"
	cicdRole          = "pipeline-service-role"
	cicdRolebinding   = "pipeline-service-role-binding"
	cicdTask          = "deploy-from-source-task"
	cicdPrPipeline    = "ci-dryrun-from-pr-pipeline"
	cicdPushPipeline  = "cd-deploy-from-push-pipeline"
	githubPrBinding   = "github-pr-binding"
	githubPushBinding = "github-push-binding"
	prTemplate        = "ci-dryrun-from-pr-template"
	pushTemplate      = "cd-deploy-from-push-template"
	cicdEventlistener = "cicd-event-listener"
	cicdRoute         = "github-webhook-event-listener"
)

func getOrderedResources() []string {
	return []string{
		cicdNamespace,
		cicdRole,
		cicdRolebinding,
		cicdTask,
		cicdPrPipeline,
		cicdPushPipeline,
		githubPrBinding,
		githubPushBinding,
		prTemplate,
		pushTemplate,
		cicdEventlistener,
		cicdRoute,
	}
}

// Initialise function will initialise the gitops directory
func Initialise(o *InitialiseParameters) error {

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

	// create gitops pipeline resources
	outputs := createPipelineResources(namespaces, o.GitOpsRepo, gitopsPath, o.Prefix)

	pipelinesPath := getPipelinesDir(gitopsPath, o.Prefix)

	if err := writeResourcesToFile(pipelinesPath, o.Prefix, outputs); err != nil {
		return err
	}

	if err := addKustomization(gitopsPath, o.Prefix); err != nil {
		return err
	}

	fmt.Println("\n\nGitops initialised")
	return nil
}

func createPipelineResources(namespaces map[string]string, gitopsRepo, rootPath, prefix string) map[string]interface{} {

	// key: name of the resource
	// value: YAML content of the resource
	outputs := make(map[string]interface{}, 0)

	// create namespace
	namespace := createNamespace(namespaces["cicd"])
	outputs[namespace.Name] = namespace

	// create roles and rolebindings for pipeline service account
	role := createRole(meta.NamespacedName(namespaces["cicd"], cicdRole), rules)
	outputs[role.Name] = role

	sa := createServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))
	rolebinding := createRoleBinding(meta.NamespacedName(namespaces["cicd"], cicdRolebinding), sa, role.Kind, role.Name)
	outputs[rolebinding.Name] = rolebinding

	// create deploy from source task
	task := tasks.GenerateDeployFromSourceTask(namespaces["cicd"], getPipelinesDir("", prefix))
	outputs[task.Name] = task

	// create ci-cd pipelines for gitops repo
	ciPipeline := createStageCIPipeline(meta.NamespacedName(namespaces["cicd"], cicdPrPipeline), namespaces["cicd"])
	outputs[ciPipeline.Name] = ciPipeline

	cdPipeline := createStageCDPipeline(meta.NamespacedName(namespaces["cicd"], cicdPushPipeline), namespaces["stage"])
	outputs[cdPipeline.Name] = cdPipeline

	// create trigger bindings
	prBinding := triggers.CreatePRBinding(namespaces["cicd"])
	outputs[prBinding.Name] = prBinding

	pushBinding := triggers.CreatePushBinding(namespaces["cicd"])
	outputs[pushBinding.Name] = pushBinding

	// create trigger templates
	ciTemplate := triggers.CreateStageCIDryRunTemplate(namespaces["cicd"], saName)
	outputs[ciTemplate.Name] = ciTemplate

	cdTemplate := triggers.CreateStageCDPushTemplate(namespaces["cicd"], saName)
	outputs[cdTemplate.Name] = cdTemplate

	// create eventlisteners and route
	eventListener := eventlisteners.Generate(gitopsRepo, namespaces["cicd"], saName)
	outputs[eventListener.Name] = eventListener

	route := routes.Generate(namespaces["cicd"])
	outputs[route.Name] = route

	return outputs
}

func getKustomizePaths(rootPath, prefix string) []string {
	return []string{
		filepath.Join(rootPath, envsDir, baseDir),
		getPipelinesDir(rootPath, prefix),
		filepath.Dir(getPipelinesDir(rootPath, prefix)),
	}
}

func writeResourcesToFile(path string, prefix string, outputs map[string]interface{}) error {

	// MkdirAll() creates necessary parent directories
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	for index, resource := range getOrderedResources() {
		if resource == cicdNamespace {
			resource = prefix + resource
		}
		filePath := filepath.Join(path, fileName(index, prefix, resource))

		// marshal data from interface{} to []byte
		data, err := yaml.Marshal(outputs[resource])
		if err != nil {
			return err
		}

		// write data to file
		if err := ioutil.WriteFile(filePath, data, 0755); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		fmt.Printf("created %v\n", resource)
	}

	return nil
}

func createFile(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return nil
}

func fileName(index int, prefix, name string) string {
	return fmt.Sprintf("%02d-%v%v.yaml", index, prefix, name)
}

func getPipelinesDir(rootPath, prefix string) string {
	return filepath.Join(rootPath, envsDir, addPrefix(prefix, cicdDir), pipelineDir)
}

func addKustomization(rootPath, prefix string) error {
	paths := getKustomizePaths(rootPath, prefix)
	for _, path := range paths {

		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create dirs:%w", err)
		}
		data := []byte("sample")
		if err := ioutil.WriteFile(filepath.Join(path, kustomize), data, 0755); err != nil {
			return fmt.Errorf("failed to write to file:%w", err)
		}

	}
	return nil
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
