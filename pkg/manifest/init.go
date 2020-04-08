package manifest

import (
	"path/filepath"
	"sort"

	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"github.com/openshift/odo/pkg/pipelines"
)

type resources map[string]interface{}

// InitParameters is a struct that provides flags for the Init command.
type InitParameters struct {
	GitOpsRepo          string
	GitOpsWebhookSecret string
	Output              string
	Prefix              string
}

// Init bootstraps a GitOps manifest and repository structure.
func Init(o *InitParameters) error {

	outputs, err := createInitialFiles(o.Prefix, o.GitOpsRepo, o.GitOpsWebhookSecret)
	if err != nil {
		return err
	}
	_, err = yaml.WriteResources(o.Output, outputs)
	return err
}

func createInitialFiles(prefix, gitOpsRepo, gitOpsWebhook string) (resources, error) {
	manifest := createManifest(prefix)
	initialFiles := resources{
		"manifest.yaml": manifest,
	}

	cicdResources, err := pipelines.CreateResources(prefix, gitOpsRepo, gitOpsWebhook)
	if err != nil {
		return nil, err
	}
	files := getResourceFiles(cicdResources)

	prefixedResources := addPrefixToResources(filepath.Join(cicdEnvironmentPath(manifest), "base/pipelines"), cicdResources)
	initialFiles = merge(prefixedResources, initialFiles)

	cicdKustomizations := addPrefixToResources(cicdEnvironmentPath(manifest), getCICDKustomization(files))
	initialFiles = merge(cicdKustomizations, initialFiles)

	return initialFiles, nil
}

func createManifest(prefix string) *config.Manifest {
	return &config.Manifest{
		Environments: []*config.Environment{
			{
				Name:   prefix + "cicd",
				IsCICD: true,
			},
		},
	}
}

func getCICDKustomization(files []string) resources {
	return resources{
		"base/kustomization.yaml": map[string]interface{}{
			"bases": []string{"./pipelines"},
		},
		"overlays/kustomization.yaml": map[string]interface{}{
			"bases": []string{"../base"},
		},
		"base/pipelines/kustomization.yaml": map[string]interface{}{
			"resources": files,
		},
	}
}

func pathForEnvironment(env *config.Environment) string {
	return filepath.Join("environments", env.Name)
}

func addPrefixToResources(prefix string, files resources) map[string]interface{} {
	updated := map[string]interface{}{}
	for k, v := range files {
		updated[filepath.Join(prefix, k)] = v
	}
	return updated
}

func merge(from, to resources) resources {
	merged := resources{}
	for k, v := range to {
		merged[k] = v
	}
	for k, v := range from {
		merged[k] = v
	}
	return merged
}

// TODO: this should probably use the .FindCICDEnvironment on the manifest.
func cicdEnvironmentPath(m *config.Manifest) string {
	return pathForEnvironment(m.Environments[0])
}

func getResourceFiles(res resources) []string {
	files := []string{}
	for k := range res {
		files = append(files, k)
	}
	sort.Strings(files)
	return files
}
