package manifest

import (
	"fmt"

	"github.com/openshift/odo/pkg/manifest/argocd"
	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"github.com/spf13/afero"
)

// BuildParameters is a struct that provides flags for the BuildResources
// command.
type BuildParameters struct {
	ManifestFilename string
	OutputDir        string
	RepositoryURL    string
}

// BuildResources builds all resources from a manifest.
func BuildResources(o *BuildParameters) error {
	m, err := config.ParseFile(o.ManifestFilename)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	appFs := afero.NewOsFs()
	resources, err := buildResources(appFs, o, m)
	if err != nil {
		return err
	}
	_, err = yaml.WriteResources(appFs, o.OutputDir, resources)
	return err
}

func buildResources(fs afero.Fs, o *BuildParameters, m *config.Manifest) (res.Resources, error) {
	resources := res.Resources{}
	envs, err := buildEnvironments(fs, m)
	if err != nil {
		return nil, err
	}
	resources = res.Merge(envs, resources)
	elFiles, err := buildEventlistener(m)
	if err != nil {
		return nil, err
	}
	resources = res.Merge(elFiles, resources)

	argoApps, err := argocd.Build(argocd.ArgoCDNamespace, o.RepositoryURL, m)
	if err != nil {
		return nil, err
	}
	resources = res.Merge(argoApps, resources)
	return resources, nil
}
