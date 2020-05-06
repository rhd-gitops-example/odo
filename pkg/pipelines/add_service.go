package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/environments"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
)

func AddService(GitRepoURL, AppName, ServiceName, WebhookSecret, EnvName, Manifest string, fs afero.Fs) error {
	m, err := config.ParseFile(fs, Manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	svc, err := createService(ServiceName, GitRepoURL)
	if err != nil {
		return err
	}
	files := res.Resources{}
	// add the secret only if CI/CD env is present
	cicdEnv, err := m.GetCICDEnvironment()
	if cicdEnv != nil && WebhookSecret != "" {
		secretName := secrets.MakeServiceWebhookSecretName(svc.Name)
		hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(cicdEnv.Name, secretName), WebhookSecret, eventlisteners.WebhookSecretKey)
		if err != nil {
			return err
		}
		svc.Webhook = &config.Webhook{
			Secret: &config.Secret{
				Name:      secretName,
				Namespace: cicdEnv.Name,
			},
		}
		secretPath := filepath.Join(config.PathForEnvironment(cicdEnv), "base", "pipelines")
		files[filepath.Join(secretPath, "03-secrets", secretName+".yaml")] = hookSecret
	}
	err = m.AddService(EnvName, AppName, svc)
	if err != nil {
		return err
	}
	err = m.Validate()
	if err != nil {
		return err
	}
	files[pipelinesFile] = m
	outputPath := filepath.Dir(Manifest)
	buildParams := &BuildParameters{
		ManifestFilename: Manifest,
		OutputPath:       outputPath,
		RepositoryURL:    m.GitOpsURL,
	}
	built, err := buildResources(fs, buildParams, m)
	if err != nil {
		return nil
	}
	files = res.Merge(built, files)
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
