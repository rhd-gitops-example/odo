package pipelines

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/environments"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/imagerepo"
	"github.com/openshift/odo/pkg/pipelines/meta"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/roles"
	"github.com/openshift/odo/pkg/pipelines/triggers"

	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
)

// AddServoceParameters are parameters passed to AddSerice function
type AddServoceParameters struct {
	AppName                  string
	EnvName                  string
	GitRepoURL               string
	ImageRepo                string
	InternalRegistryHostname string
	Manifest                 string
	ServiceName              string
	WebhookSecret            string
}

func AddService(p *AddServoceParameters, fs afero.Fs) error {

	m, err := config.ParseFile(fs, p.Manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	cicdEnv, err := m.GetCICDEnvironment()
	if err != nil {
		return err
	}
	outputPath := filepath.Dir(p.Manifest)

	files, err := serviceResources(m, fs, p)
	if err != nil {
		return err
	}

	_, err = yaml.WriteResources(fs, outputPath, files)
	if err != nil {
		return err
	}
	if cicdEnv != nil {
		base := filepath.Join(outputPath, config.PathForEnvironment(cicdEnv), "base", "pipelines")
		err = updateKustomization(fs, base)
		if err != nil {
			return err
		}
	}
	return nil
}

func serviceResources(m *config.Manifest, fs afero.Fs, p *AddServoceParameters) (res.Resources, error) {
	files := res.Resources{}

	svc, err := createService(p.ServiceName, p.GitRepoURL)
	if err != nil {
		return nil, err
	}

	cicdEnv, err := m.GetCICDEnvironment()
	if err != nil {
		return nil, err
	}
	if cicdEnv != nil && p.WebhookSecret == "" && p.GitRepoURL != "" {
		return nil, fmt.Errorf("The webhook secret is required")
	}
	// add the secret only if CI/CD env is present
	if cicdEnv != nil {
		secretName := secrets.MakeServiceWebhookSecretName(svc.Name)
		hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(cicdEnv.Name, secretName), p.WebhookSecret, eventlisteners.WebhookSecretKey)
		if err != nil {
			return nil, err
		}

		svc.Webhook = &config.Webhook{
			Secret: &config.Secret{
				Name:      secretName,
				Namespace: cicdEnv.Name,
			},
		}
		secretPath := filepath.Join(config.PathForEnvironment(cicdEnv), "base", "pipelines")
		files[filepath.Join(secretPath, "03-secrets", secretName+".yaml")] = hookSecret

		if p.ImageRepo != "" {
			resources, bindingName, err := createImageRepoResources(cicdEnv, p)
			if err != nil {
				return nil, err
			}

			env := m.GetEnvironment(p.EnvName)
			if env == nil {
				return nil, fmt.Errorf("environment %s does not exist.", p.EnvName)
			}

			files = res.Merge(resources, files)
			svc.Pipelines = &config.Pipelines{
				Integration: &config.TemplateBinding{
					Bindings: append([]string{bindingName}, env.Pipelines.Integration.Bindings[:]...),
				},
			}
		}

	}

	err = m.AddService(p.EnvName, p.AppName, svc)
	if err != nil {
		return nil, err
	}
	err = m.Validate()
	if err != nil {
		return nil, err
	}

	files[filepath.Base(p.Manifest)] = m
	outputPath := filepath.Dir(p.Manifest)
	buildParams := &BuildParameters{
		ManifestFilename: p.Manifest,
		OutputPath:       outputPath,
		RepositoryURL:    m.GitOpsURL,
	}
	built, err := buildResources(fs, buildParams, m)
	if err != nil {
		return nil, err
	}
	files = res.Merge(built, files)

	return files, nil

}

func createImageRepoResources(cicdEnv *config.Environment, p *AddServoceParameters) (res.Resources, string, error) {
	isInternalRegistry, imageRepo, err := imagerepo.ValidateImageRepo(p.ImageRepo, p.InternalRegistryHostname)
	if err != nil {
		return nil, "", err
	}

	resources := res.Resources{}

	bindingName, _, svcImageBinding := createSvcImageBinding(cicdEnv, p.EnvName, p.ServiceName, imageRepo, !isInternalRegistry)
	resources = res.Merge(svcImageBinding, resources)

	if isInternalRegistry {
		regRes, err := imagerepo.CreateInternalRegistryResources(cicdEnv, roles.CreateServiceAccount(meta.NamespacedName(cicdEnv.Name, saName)), imageRepo)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get resources for internal image repository: %w", err)
		}
		resources = res.Merge(regRes, resources)
	}

	return resources, bindingName, nil
}

func createService(serviceName, url string) (*config.Service, error) {
	if url == "" {
		return &config.Service{
			Name: serviceName,
		}, nil
	}
	return &config.Service{
		Name:      serviceName,
		SourceURL: url,
	}, nil
}

func updateKustomization(fs afero.Fs, base string) error {
	files := res.Resources{}
	list, err := environments.ListFiles(fs, base)
	if err != nil {
		return err
	}
	files[Kustomize] = &res.Kustomization{Resources: environments.ExtractFilenames(list)}
	_, err = yaml.WriteResources(fs, base, files)
	return err
}

func makeSvcImageBindingName(envName, svcName string) string {
	return fmt.Sprintf("%s-%s-binding", envName, svcName)
}

func makeSvcImageBindingFilename(bindingName string) string {
	return filepath.Join("06-bindings", bindingName+".yaml")
}

func makeImageBindingPath(cicdEnv *config.Environment, imageRepoBindingFilename string) string {
	return filepath.Join(config.PathForEnvironment(cicdEnv), "base", "pipelines", imageRepoBindingFilename)
}

func createSvcImageBinding(cicdEnv *config.Environment, envName, svcName, imageRepo string, isTLSVerify bool) (string, string, res.Resources) {
	name := makeSvcImageBindingName(envName, svcName)
	filename := makeSvcImageBindingFilename(name)
	resourceFilePath := makeImageBindingPath(cicdEnv, filename)
	return name, filename, res.Resources{resourceFilePath: triggers.CreateImageRepoBinding(cicdEnv.Name, name, imageRepo, strconv.FormatBool(isTLSVerify))}
}
