package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/eventlisteners"

	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/roles"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/types"
)

// AddParameters is a struct that provides flags for add application command
type AddParameters struct {
	GitopsRepo           string
	GitopsWebhookSecret  string
	Output               string
	Prefix               string
	AppName              string
	ServiceWebhookSecret string
	ServiceImageRepo     string
	EnvName              string
	DockerCfgJson        string
	ServicesGitRepo      string
	SkipChecks           bool
}

const (
	overlays     = "overlays"
	appDir       = "apps"
	configDir    = "config"
	servicesDir  = "services"
	secretName   = "quay"
	configSApath = "base/config/serviceaccount.yaml"
	secretPath   = "base/config/secret.yaml"
	PatchPath    = "overlays/eventlistener_patch.yaml"
	overlaysDir  = "overlays"
	webhookPath  = "base/config/app-webhook-secret.yaml"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.

type patchStringValue struct {
	Op    string                          `yaml:"op"`
	Path  string                          `yaml:"path"`
	Value triggersv1.EventListenerTrigger `yaml:"value"`
}

// Add_Application function will initialise the gitops directory
func Add_Application(o *AddParameters) error {

	if !o.SkipChecks {
		installed, err := checkTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	gitopsName := getGitopsRepoName(o.GitopsRepo)

	// AppRepo := getGitopsRepoName(o.AppGitRepo)

	ServiceRepo := getGitopsRepoName(o.ServicesGitRepo)

	gitopsPath := filepath.Join(o.Output, gitopsName)

	outputs := map[string]interface{}{}

	exists, _ := isExisting(gitopsPath)

	createPatchFiles(outputs, o.EnvName, o.AppName)

	if !exists {
		return fmt.Errorf("%s does not exist at %s", gitopsName, gitopsPath)
	}
	configPath := filepath.Join(gitopsPath, servicesDir, ServiceRepo)

	CreateNewKustomiseFile(outputs, filepath.Join(overlaysDir, kustomize))

	files := createResourcesConfig(outputs, o.EnvName, o.DockerCfgJson, o.EnvName, o.ServiceWebhookSecret)

	_, err := writeResources(configPath, files)

	if err != nil {
		return err
	}

	if err := addKustomize("bases", []string{"../base"}, filepath.Join(gitopsPath, appDir, o.AppName, overlays, kustomize)); err != nil {
		return err
	}
	if err := addKustomize("bases", []string{fmt.Sprintf("../../../services/%s/overlays", o.ServicesGitRepo)}, filepath.Join(gitopsPath, appDir, o.AppName, baseDir, kustomize)); err != nil {
		return err
	}
	if err := addKustomize("bases", []string{"../config"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, baseDir, kustomize)); err != nil {
		return err
	}
	if err := addKustomize("bases", []string{"../base"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, baseDir, kustomize)); err != nil {
		return err
	}

	kustomize1 := map[string][]string{
		"bases":     []string{fmt.Sprintf("../../../../envs/%s/overlays/", o.EnvName)},
		"resources": []string{"secret.yaml", "serviceaccount.yaml", "app-webhook-secret.yaml"},
	}

	if err := addModKustomize(kustomize1, filepath.Join(gitopsPath, servicesDir, ServiceRepo, baseDir, configDir, kustomize)); err != nil {
		return err
	}

	return nil
}

func addModKustomize(values map[string][]string, path string) error {
	content := make([]interface{}, 0)
	for name, items := range values {
		content = append(content, map[string]interface{}{name: items})
	}
	return marshalItemsToFile(path, content)
}

func createResourcesConfig(outputs map[string]interface{}, namespace, dockerFilePath, envName, serviceWebhookSecret string) map[string]interface{} {
	sa := roles.CreateServiceAccount(meta.NamespacedName(namespace, saName))
	ServiceAcc := roles.AddSecretToSA(sa, secretName)
	outputs[configSApath] = ServiceAcc
	secretsDocker, _ := CreateDockerSecret(dockerFilePath, namespace)
	outputs[secretPath] = secretsDocker
	githubSecret, _ := secrets.CreateSealedSecret(meta.NamespacedName(envName, eventlisteners.GitOpsWebhookSecret),
		serviceWebhookSecret, eventlisteners.WebhookSecretKey)
	outputs[webhookPath] = githubSecret

	return outputs
}

func createPatchFiles(outputs map[string]interface{}, name, repo string) {
	t := []patchStringValue{
		patchStringValue{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger(name, fmt.Sprintf(eventlisteners.StageCIDryRunFilters, repo), "sample", "github-pr-binding", "app-ci-template"),
		},
		patchStringValue{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger(name, fmt.Sprintf(eventlisteners.StageCDDeployFilters, repo), "sample", "github-push-binding", "app-cd-template"),
		},
	}
	outputs[PatchPath] = &t

}

func CreateNewKustomiseFile(outputs map[string]interface{}, path string) {

	bases := []string{"../resources", "../resources2"}

	GVK := resid.Gvk{
		Group:   "tekton.dev",
		Version: "v1alpha1",
		Kind:    "EventListener",
	}
	target := &types.PatchTarget{
		Gvk:  GVK,
		Name: "cicd-event-listener",
	}
	Patches := []types.PatchJson6902{
		types.PatchJson6902{
			Target: target,
			Path:   "eventlistener_patch.yaml",
		},
	}
	file := types.Kustomization{
		Bases:           bases,
		Resources:       bases,
		PatchesJson6902: Patches,
	}
	outputs[path] = &file

}
