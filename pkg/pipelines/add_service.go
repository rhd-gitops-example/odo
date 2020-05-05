package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
	"sigs.k8s.io/kustomize/pkg/types"
)

// AddService adds a new service to an environment
func AddService(gitRepoURL, webhookSecret, envName, appName, manifest string, fs afero.Fs) error {

	m, err := config.ParseFile(fs, manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	if err := m.Validate(); err != nil {
		return err
	}

	env := m.GetEnvironment(envName)
	if env == nil {
		return fmt.Errorf("environment %s does not exist", envName)
	}

	repoName, err := repoFromURL(gitRepoURL)
	if err != nil {
		return fmt.Errorf("Git repository URL is invalid: %w", err)
	}

	cicdEnv, err := m.GetCICDEnvironment()
	if err != nil {
		return fmt.Errorf("failed to get CICD environment: %w", err)
	}

	secretName := secrets.MakeServiceWebhookSecretName(repoName)
	secretFileName := filepath.Join("03-secrets", secretName+".yaml")
	secretsPath := filepath.Join("environments", cicdEnv.Name, "base", "pipelines", secretFileName)
	secretsKustomisePath := filepath.Join("environments", envName, "services", repoName, "base", "kustomization.yaml")

	files := resources.Resources{}

	app, _ := m.GetApplication(envName, appName)

	if app == nil {
		app, err := ApplicationFromName(appName, gitRepoURL, secretName, cicdEnv.Name)
		if err != nil {
			return err
		}
		env.Apps = append(env.Apps, app)
	} else {
		err := checkServiceExists(app.Services, repoName)
		if err != nil {
			return err
		}
		app.Services = append(app.Services, GetService(repoName, gitRepoURL, secretName, cicdEnv.Name))
		env.Apps = []*config.Application{app}
	}

	hookSecret, err := secrets.CreateSealedSecret(
		meta.NamespacedName(cicdEnv.Name, secretName),
		webhookSecret,
		eventlisteners.WebhookSecretKey)

	files[filepath.Base(manifest)] = m
	files[secretsPath] = hookSecret

	outputPath := filepath.Dir(manifest)
	buildParams := &BuildParameters{
		ManifestFilename: manifest,
		OutputPath:       outputPath,
	}

	built, err := buildResources(fs, buildParams, m)
	if err != nil {
		return fmt.Errorf("failed to build resources: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to build resources: %w", err)
	}

	files = resources.Merge(built, files)
	kustFiles := appKustomise([]string{secretsPath}, []string{"./config"})
	files[secretsKustomisePath] = kustFiles

	_, err = yaml.WriteResources(fs, outputPath, files)
	if err != nil {
		return err
	}
	return nil
}

func appKustomise(bases, resources []string) *types.Kustomization {
	return &types.Kustomization{
		Resources: bases,
		Bases:     resources,
	}
}

func checkServiceExists(services []*config.Service, serviceName string) error {
	for _, service := range services {
		if service.Name == serviceName {
			return fmt.Errorf("A Service with this name %s already exists", serviceName)
		}
	}
	return nil
}
