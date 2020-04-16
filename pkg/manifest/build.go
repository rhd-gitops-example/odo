package manifest

import (
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"

	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/openshift/odo/pkg/manifest/meta"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/openshift/odo/pkg/manifest/yaml"
)

// BuildParameters is a struct that provides flags for the BuildResources
// command.
type BuildParameters struct {
	ManifestFilename string
	OutputDir        string
}

// BuildResources builds all resources from a manifest.
func BuildResources(o *BuildParameters) error {
	m, err := config.ParseFile(o.ManifestFilename)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	resources := buildResources(m)
	_, err = yaml.WriteResources(o.OutputDir, resources)
	return err
}

func buildResources(m *config.Manifest) map[string]interface{} {
	resources := map[string]interface{}{}
	resources = res.Merge(buildEnvironments(m), resources)
	return resources
}

func buildEnvironments(m *config.Manifest) map[string]interface{} {
	files := map[string]interface{}{}
	for _, env := range m.Environments {
		filename := configPathForEnvironment(env, "01_namespace.yaml")
		files[filename] = namespace(env)
	}

	return files
}

func configPathForEnvironment(env *config.Environment, elements ...string) string {
	envPath := pathForEnvironment(env)
	return filepath.Join(append([]string{envPath, "base", "config"}, elements...)...)
}

func namespace(env *config.Environment) corev1.Namespace {
	return corev1.Namespace{
		TypeMeta:   meta.TypeMeta("Namespace", "v1"),
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("", env.Name)),
	}
}
