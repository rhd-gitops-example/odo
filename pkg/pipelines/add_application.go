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

// EnvParameters encapsulates parameters for add env command
type ServiceParameters struct {
	ServiceGitRepo       string
	ServiceWebhookSecret string
	AppName              string
	EnvName              string
	Output               string
	Manifest             string
}

const filename = "pipelines.yaml"

// Env will bootstrap a new environment directory
func AddService(o *ServiceParameters, fs afero.Fs) error {

	repoName, err := repoFromURL(o.ServiceGitRepo)
	ns := NamespaceNames("")
	secretName := "github-webhook-secret-" + repoName + "-svc"

	secretFileName := filepath.Join("03-secrets", secretName+".yaml")
	secretsPath := filepath.Join("environments", ns["cicd"], "base", "pipelines", secretFileName)
	secretsKustomisePath := filepath.Join("environments", o.EnvName, "services", repoName, "base", "kustomization.yaml")

	files := resources.Resources{}

	m, err := config.ParseFile(fs, o.Manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	env, err := m.GetEnvironment(o.EnvName)
	if env == nil {
		return fmt.Errorf("environment does not exist already at %s", o.EnvName)
	}

	app, err := m.GetApplication(o.EnvName, o.AppName)
	if err != nil {
		return err
	}
	if app == nil {
		Newapp, _ := ApplicationFromRepo(o.AppName, o.ServiceGitRepo, secretName, ns["cicd"])
		env.Apps = append(env.Apps, Newapp)
	} else {
		err := checkServiceExists(app.Services, repoName)
		if err != nil {
			return err
		}
		service := GetService(repoName, o.ServiceGitRepo, secretName, ns["cicd"])
		app.Services = append(app.Services, service)
		env.Apps = []*config.Application{app}
	}

	hookSecret, err := secrets.CreateSealedSecret(
		meta.NamespacedName(ns["cicd"], secretName),
		o.ServiceWebhookSecret,
		eventlisteners.WebhookSecretKey)

	files[o.Manifest] = m
	files[secretsPath] = hookSecret

	// m1, err := config.ParseFile(fs, o.Manifest)

	// env1, err := m1.GetEnvironment("cicd")
	// log.Println(env1)

	buildParams := &BuildParameters{
		ManifestFilename: o.Manifest,
		OutputPath:       o.Output,
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

	_, err = yaml.WriteResources(fs, o.Output, files)
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
