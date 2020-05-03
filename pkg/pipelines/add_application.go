package pipelines

import (
	"fmt"
	"log"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/yaml"

	"github.com/openshift/odo/pkg/pipelines/resources"

	"github.com/spf13/afero"
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

	m, err := config.ParseFile(fs, o.Manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	env, err := m.GetEnvironment(o.EnvName)
	if env == nil {
		return fmt.Errorf("environment does not exist already", o.EnvName)
	}

	repoName, err := repoFromURL(o.ServiceGitRepo)

	secretName := "github-webhook-secret-" + repoName + "-svc"

	app, err := m.GetApplication(o.EnvName, o.AppName)

	if app == nil {
		Newapp, _ := ApplicationFromRepo(o.AppName, o.ServiceGitRepo, secretName, "cicd")
		env.Apps = append(env.Apps, Newapp)
		log.Println("was here in the if")
	} else {
		log.Println("Was here in the else")
		service := GetService(repoName, o.ServiceGitRepo, secretName, "cicd")
		app.Services = append(app.Services, service)
		env.Apps = []*config.Application{app}
	}

	files := resources.Resources{}

	files[o.Manifest] = m

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

	_, err = yaml.WriteResources(fs, o.Output, files)
	return err
}
