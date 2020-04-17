package manifest

import (
	"sort"

	"github.com/openshift/odo/pkg/manifest/config"
)

// BootstrapParameters is a struct that provides the optional flags
type BootstrapParameters struct {
	GitOpsWebhookSecret      string
	GitOpsRepo               string
	Output                   string
	DeploymentPath           string
	GitHubHookSecret         string
	GitHubToken              string
	GitRepo                  string
	InternalRegistryHostname string
	ImageRepo                string
	Prefix                   string
	DockerConfigJSONFilename string
	SkipChecks               bool
}

// Bootstrap bootstraps a GitOps manifest and repository structure.
func Bootstrap(o *BootstrapParameters) error {
	return nil
}

func bootstrapResources(p *BootstrapParameters) (resources, error) {
	bootstrapped := resources{}
	ns := NamespaceNames(p.Prefix)
	bootstrapped["manifest.yaml"] = createManifest(bootstrapEnvironments(ns)...)

	return bootstrapped, nil
}

func bootstrapEnvironments(ns map[string]string) []*config.Environment {
	envs := []*config.Environment{}
	for k, v := range ns {
		env := &config.Environment{Name: v}
		if k == "cicd" {
			env.IsCICD = true
		}
		envs = append(envs, env)
	}
	envs = append(envs, &config.Environment{Name: "argocd", IsArgoCD: true})
	sort.Sort(config.ByName(envs))
	return envs
}
