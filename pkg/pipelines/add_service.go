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

type AddOptions struct {
	AppName       string
	EnvName       string
	GitRepoURL    string
	Manifest      string
	WebhookSecret string
}

func AddService(o *AddOptions, fs afero.Fs) error {
	m, err := config.ParseFile(fs, o.Manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	svc, err := createService(o.GitRepoURL)
	if err != nil {
		return err
	}
	files := res.Resources{}
	// add the secret only if CI/CD env is present
	cicdEnv, err := m.GetCICDEnvironment()
	if cicdEnv != nil {
		secretName := secrets.MakeServiceWebhookSecretName(svc.Name)
		hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName(cicdEnv.Name, secretName), o.WebhookSecret, eventlisteners.WebhookSecretKey)
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
	err = m.AddService(o.EnvName, o.AppName, svc)
	if err != nil {
		return err
	}
	err = m.Validate()
	if err != nil {
		return err
	}
	files[pipelinesFile] = m
	outputPath := filepath.Dir(o.Manifest)
	buildParams := &BuildParameters{
		ManifestFilename: o.Manifest,
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

func createService(url string) (*config.Service, error) {
	svcName, err := repoFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("Git repository URL is invalid: %w", err)
	}
	return &config.Service{
		Name:      svcName,
		SourceURL: url,
	}, nil
}
