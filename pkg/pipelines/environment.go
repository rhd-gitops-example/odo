package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/config"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
)

const (
	envNamespace   = "namespace.yaml"
	envRoleBinding = "rolebinding.yaml"
)

// EnvParameters encapsulates parameters for add env command
type EnvParameters struct {
	ManifestFilename string
	OutputPath       string
	EnvName          string
}

// AddEnv adds a new environment to the manifest.
func AddEnv(o *EnvParameters, appFs afero.Fs) error {
	m, err := config.ParseFile(appFs, o.ManifestFilename)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	env := m.GetEnvironment(o.EnvName)
	if env != nil {
		return fmt.Errorf("environment %s already exists", o.EnvName)
	}
	files := res.Resources{}
	m.Environments = append(m.Environments, &config.Environment{Name: o.EnvName})
	manifestPath, err := filepath.Rel(o.OutputPath, o.ManifestFilename)
	if err != nil {
		return err
	}
	files[manifestPath] = m
	buildParams := &BuildParameters{
		ManifestFilename: o.ManifestFilename,
		OutputPath:       o.OutputPath,
	}
	built, err := buildResources(appFs, buildParams, m)
	if err != nil {
		return fmt.Errorf("failed to build resources: %w", err)
	}
	files = res.Merge(built, files)
	_, err = yaml.WriteResources(appFs, o.OutputPath, files)
	return err
}
