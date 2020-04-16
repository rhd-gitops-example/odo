package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/manifest"
	"github.com/openshift/odo/pkg/manifest/eventlisteners"
	"github.com/openshift/odo/pkg/manifest/ioutils"
	"github.com/openshift/odo/pkg/manifest/pipelines"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/openshift/odo/pkg/manifest/meta"
	"github.com/openshift/odo/pkg/manifest/secrets"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
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
	ServiceGitRepo       string
	SkipChecks           bool
}

const (
	appDir           = "apps"
	appWebhookSecret = "app-webhook-secret"
	configDir        = "config"
	configSApath     = "base/config/serviceaccount.yaml"
	overlaysDir      = "overlays"
	// PatchPath path to eventlistener patch yaml
	PatchPath        = "overlays/eventlistener_patch.yaml"
	servicesDir      = "services"
	secretPath       = "base/config/secret.yaml"
	webhookPath      = "base/config/app-webhook-secret.yaml"
	kustomizeModPath = "base/config/kustomization.yaml"
	secretName       = "secret"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.

type patchStringValue struct {
	Op    string                          `json:"op"`
	Path  string                          `json:"path"`
	Value triggersv1.EventListenerTrigger `json:"value"`
}

// CreateApplication creates an application
func CreateApplication(o *AddParameters) error {

	if !o.SkipChecks {
		installed, err := pipelines.CheckTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	ServiceRepo := getGitopsRepoName(o.ServiceGitRepo)

	secretName := fmt.Sprintf("svc-%s-secret", ServiceRepo)

	// we simpily output to the output dir, no gitops repo in the output path
	gitopsPath := o.Output

	outputs := map[string]interface{}{}

	exists, _ := ioutils.IsExisting(gitopsPath)

	if !exists {
		return fmt.Errorf("Output does not exist at %s", gitopsPath)
	}

	configPath := filepath.Join(gitopsPath, servicesDir, ServiceRepo)

	createPatchFiles(outputs, o.ServiceGitRepo)

	CreatePatchKustomiseFile(outputs, filepath.Join(overlaysDir, manifest.Kustomize))

	environmentName := manifest.NamespaceNames(o.Prefix)

	createKustomizeMod(outputs, kustomizeModPath, environmentName["cicd"])

	files := createResourcesConfig(outputs, o.ServiceWebhookSecret, environmentName["cicd"], secretName)

	_, err := yaml.WriteResources(configPath, files)

	if err != nil {
		return err
	}
	if err := yaml.AddKustomize("bases", []string{"overlays"}, filepath.Join(gitopsPath, appDir, o.AppName, manifest.Kustomize)); err != nil {
		return err
	}

	if err := yaml.AddKustomize("bases", []string{"../base"}, filepath.Join(gitopsPath, appDir, o.AppName, overlaysDir, manifest.Kustomize)); err != nil {
		return err
	}

	if err := yaml.AddKustomize("bases", []string{fmt.Sprintf("../../../services/%s/overlays", ServiceRepo)}, filepath.Join(gitopsPath, appDir, o.AppName, manifest.BaseDir, manifest.Kustomize)); err != nil {
		return err
	}
	if err := yaml.AddKustomize("bases", []string{"../config"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, manifest.BaseDir, manifest.Kustomize)); err != nil {
		return err
	}
	if err := yaml.AddKustomize("bases", []string{"./config"}, filepath.Join(gitopsPath, servicesDir, ServiceRepo, manifest.BaseDir, manifest.Kustomize)); err != nil {
		return err
	}

	return nil
}

func addModKustomize(values map[string][]string, path string) error {
	content := make([]interface{}, 0)
	for name, items := range values {
		content = append(content, map[string]interface{}{name: items})
	}
	return yaml.MarshalItemToFile(path, content)
}

func createResourcesConfig(outputs map[string]interface{}, serviceWebhookSecret, environmentName, secretName string) map[string]interface{} {

	githubSecret, _ := secrets.CreateSealedSecret(meta.NamespacedName(environmentName, secretName),
		serviceWebhookSecret, eventlisteners.WebhookSecretKey)
	outputs[webhookPath] = githubSecret

	return outputs
}

func createPatchFiles(outputs map[string]interface{}, serviceRepo string) {
	t := []patchStringValue{
		{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-ci-build-from-pr", eventlisteners.StageCIDryRunFilters, serviceRepo, "github-pr-binding", "app-ci-template"),
		},
		{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-cd-deploy-from-master", eventlisteners.StageCDDeployFilters, serviceRepo, "github-push-binding", "app-cd-template"),
		},
	}
	outputs[PatchPath] = t

}

// CreatePatchKustomiseFile creates patch kustomization file
func CreatePatchKustomiseFile(outputs map[string]interface{}, path string) {

	bases := []string{"../base"}

	GVK := gvk.Gvk{
		Group:   "tekton.dev",
		Version: "v1alpha1",
		Kind:    "EventListener",
	}
	target := &types.PatchTarget{
		Gvk:  GVK,
		Name: "cicd-event-listener",
	}
	Patches := []types.PatchJson6902{
		{
			Target: target,
			Path:   "eventlistener_patch.yaml",
		},
	}
	file := types.Kustomization{
		Bases:           bases,
		PatchesJson6902: Patches,
	}
	outputs[path] = file

}

func createKustomizeMod(outputs map[string]interface{}, path, environmentName string) {

	bases := []string{fmt.Sprintf("../../../../environments/%s/overlays", environmentName)}
	resources := []string{"app-webhook-secret.yaml"}

	file := types.Kustomization{
		Bases:     bases,
		Resources: resources,
	}

	outputs[path] = file

}

func getGitopsRepoName(repo string) string {
	return strings.Split(repo, "/")[1]
}
