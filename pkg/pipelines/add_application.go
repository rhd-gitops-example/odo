package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/manifest/eventlisteners"

	"github.com/openshift/odo/pkg/manifest/meta"
	"github.com/openshift/odo/pkg/manifest/roles"
	"github.com/openshift/odo/pkg/manifest/secrets"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/types"
)

// AddParameters is a struct that provides flags for add application command
type AddParameters struct {
	AppName              string
	EnvName              string
	GitopsRepo           string
	GitopsWebhookSecret  string
	Output               string
	Prefix               string
	ServiceWebhookSecret string
	ServiceImageRepo     string
	ServicesGitRepo      string
	SkipChecks           bool
}

const (
	appDir           = "apps"
	appWebhookSecret = "app-webhook-secret"
	configDir        = "config"
	configSApath     = "base/config/serviceaccount.yaml"
	overlaysDir      = "overlays"
	PatchPath        = "overlays/eventlistener_patch.yaml"
	servicesDir      = "services"
	secretName       = "quay"
	secretPath       = "base/config/secret.yaml"
	webhookPath      = "base/config/app-webhook-secret.yaml"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.

type patchStringValue struct {
	Op    string                          `json:"op"`
	Path  string                          `json:"path"`
	Value triggersv1.EventListenerTrigger `json:"value"`
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

	ServiceRepo := getGitopsRepoName(o.ServicesGitRepo)

	gitopsPath := filepath.Join(o.Output, gitopsName)

	outputs := map[string]interface{}{}

	exists, _ := isExisting(gitopsPath)

	if !exists {
		return fmt.Errorf("%s does not exist at %s", gitopsName, gitopsPath)
	}

	configPath := filepath.Join(gitopsPath, servicesDir, ServiceRepo)

	createPatchFiles(outputs, o.ServicesGitRepo)

	CreatePatchKustomiseFile(outputs, filepath.Join(overlaysDir, kustomize))

	environmentName := namespaceNames(o.Prefix)

	files := createResourcesConfig(outputs, o.ServiceWebhookSecret, environmentName["cicd"])

	_, err := writeResources(configPath, files)

	if err != nil {
		return err
	}
	if err := addKustomize("bases", []string{"overlays"}, filepath.Join(gitopsPath, appDir, o.AppName, kustomize)); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{"../base"}, filepath.Join(gitopsPath, appDir, o.AppName, overlaysDir, kustomize)); err != nil {
		return err
	}

	if err := addKustomize("bases", []string{fmt.Sprintf("../../../services/%s/overlays", ServiceRepo)}, filepath.Join(gitopsPath, appDir, o.AppName, baseDir, kustomize)); err != nil {
		return err
	}
	if err := addKustomize("bases", []string{"../config"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, baseDir, kustomize)); err != nil {
		return err
	}
	if err := addKustomize("bases", []string{"./config"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, baseDir, kustomize)); err != nil {
		return err
	}
	kustomize1 := map[string][]string{
		"bases":     []string{fmt.Sprintf("../../../../envs/%s/", environmentName["cicd"])},
		"resources": []string{"serviceaccount.yaml", "app-webhook-secret.yaml"},
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

func createResourcesConfig(outputs map[string]interface{}, serviceWebhookSecret, environmentName string) map[string]interface{} {
	sa := roles.CreateServiceAccount(meta.NamespacedName(environmentName, saName))
	ServiceAcc := roles.AddSecretToSA(sa, secretName)
	outputs[configSApath] = ServiceAcc
	githubSecret, _ := secrets.CreateSealedSecret(meta.NamespacedName(environmentName, eventlisteners.GitOpsWebhookSecret),
		serviceWebhookSecret, eventlisteners.WebhookSecretKey)
	outputs[webhookPath] = githubSecret

	return outputs
}

func createPatchFiles(outputs map[string]interface{}, servicesRepo string) {
	t := []patchStringValue{
		patchStringValue{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-ci-build-from-pr", eventlisteners.StageCIDryRunFilters, servicesRepo, "github-pr-binding", "app-ci-template"),
		},
		patchStringValue{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-cd-deploy-from-master", eventlisteners.StageCDDeployFilters, servicesRepo, "github-push-binding", "app-cd-template"),
		},
	}
	outputs[PatchPath] = &t

}

func CreatePatchKustomiseFile(outputs map[string]interface{}, path string) {

	bases := []string{"../base"}

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
		PatchesJson6902: Patches,
	}
	outputs[path] = &file

}
