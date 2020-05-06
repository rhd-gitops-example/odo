package pipelines

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
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
	orgRepo, err := orgRepoFromURL(gitRepoURL)
	if err != nil {
		return fmt.Errorf("Git repository URL is invalid: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to get CICD environment: %w", err)
	}
	cicdEnv, err := m.GetCICDEnvironment()

	secretName := secrets.MakeServiceWebhookSecretName(repoName)
	secretFileName := filepath.Join("03-secrets", secretName+".yaml")
	secretsPath := filepath.Join("environments", cicdEnv.Name, "base", "pipelines", secretFileName)
	pipelineFileName := filepath.Join("08-eventlisteners", repoName+"-cicd-event-listener.yaml")
	kustomizePath := filepath.Join("environments", "cicd", "base", "pipelines", "kustomization.yaml")

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
		app.Services = append(app.Services, GetService(repoName, gitRepoURL))
		env.Apps = []*config.Application{app}
	}

	cicdExists := m.GetEnvironment(envName)
	if cicdExists == nil {
		service := GetService(repoName, gitRepoURL)

		app.Services = append(app.Services, service)

		env.Apps = []*config.Application{app}

		err := generateService(manifest, fs, files, m)
		if err != nil {
			return err
		}
	} else {

		updatedService := GetServiceSecret(repoName, gitRepoURL, secretName, "cicd")
		if updatedService == nil {
			return fmt.Errorf("This failed bad")
		}
		app.Services = append(app.Services, updatedService)
		if app.Services == nil {
			return fmt.Errorf("This failed bad")
		}
		env.Apps = []*config.Application{app}

		err := appendKustomiseResources(pipelineFileName, secretFileName, kustomizePath, manifest, fs)
		if err != nil {
			return err
		}

		hookSecret, err := secrets.CreateSealedSecret(
			meta.NamespacedName(cicdEnv.Name, secretName),
			webhookSecret,
			eventlisteners.WebhookSecretKey)

		eventListenerService := eventlisteners.Generate(orgRepo, "cicd", saName, eventlisteners.GitOpsWebhookSecret)
		files[eventListenerPath] = eventListenerService
		files[secretsPath] = hookSecret

		err = generateService(manifest, fs, files, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateService(manifest string, fs afero.Fs, files resources.Resources, m *config.Manifest) error {
	outputPath := filepath.Dir(manifest)
	buildParams := &BuildParameters{
		ManifestFilename: manifest,
		OutputPath:       outputPath,
	}
	files[manifest] = m

	built, err := buildResources(fs, buildParams, m)
	if err != nil {
		return fmt.Errorf("failed to build resources: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to build resources: %w", err)
	}
	files = resources.Merge(built, files)

	_, err = yaml.WriteResources(fs, outputPath, files)
	if err != nil {
		return err
	}
	return nil
}

func checkServiceExists(services []*config.Service, serviceName string) error {
	for _, service := range services {
		if service.Name == serviceName {
			return fmt.Errorf("A Service with this name %s already exists", serviceName)
		}
	}
	return nil
}

func appendKustomiseResources(pipelineFileName, secretFileName, kustomizePath, manifest string, fs afero.Fs) error {
	outputPath := filepath.Dir(manifest)
	f, err := fs.OpenFile(filepath.Join(outputPath, kustomizePath), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Could not locate the file %s present", err)
	}
	defer f.Close()
	if _, err = f.WriteString("\n- " + secretFileName); err != nil {
		return fmt.Errorf("The secret kustomization file could not be appended to the kustomization.yaml file")
	}
	if _, err = f.WriteString("\n- " + pipelineFileName); err != nil {
		return fmt.Errorf("The eventlistener kustomization file could not be appended to the kustomization.yaml file")
	}
	return nil
}
