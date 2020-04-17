package manifest

import (
	"net/url"
	"sort"
	"strings"

	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"github.com/spf13/afero"
)

// BootstrapOptions is a struct that provides the optional flags
type BootstrapOptions struct {
	GitOpsRepoURL            string // This is where the manifest and configuration are.
	GitOpsWebhookSecret      string // This is the secret for authenticating hooks from your GitOps repo.
	AppRepoURL               string // This is the full URL to your GitHub repository for your app source.
	AppWebhookSecret         string // This is the secret for authenticating hooks from your app source.
	InternalRegistryHostname string // This is the internal registry hostname used for pushing images.
	ImageRepo                string // This is where built images are pushed to.
	Prefix                   string // Used to prefix generated environment names in a shared cluster.
	OutputPath               string // Where to write the bootstrapped files to?
	DockerConfigJSONFilename string
	SkipChecks               bool // Don't check that the local cluster has Tekton installed
}

// Bootstrap bootstraps a GitOps manifest and repository structure.
func Bootstrap(o *BootstrapOptions) error {
	bootstrapped, err := bootstrapResources(o)

	appFs := afero.NewOsFs()
	_, err = yaml.WriteResources(appFs, o.OutputPath, bootstrapped)
	return err
}

func bootstrapResources(p *BootstrapOptions) (res.Resources, error) {
	orgRepo, err := orgRepoFromURL(p.GitOpsRepoURL)
	if err != nil {
		return nil, err
	}
	bootstrapped, err := createInitialFiles(p.Prefix, orgRepo, p.GitOpsWebhookSecret, p.DockerConfigJSONFilename, "")
	if err != nil {
		return nil, err
	}
	ns := NamespaceNames(p.Prefix)
	envs, err := bootstrapEnvironments(p.Prefix, p.AppRepoURL, ns)
	if err != nil {
		return nil, err
	}
	bootstrapped["manifest.yaml"] = createManifest(envs...)
	return bootstrapped, nil
}

func bootstrapEnvironments(prefix, repoURL string, ns map[string]string) ([]*config.Environment, error) {
	envs := []*config.Environment{}
	for k, v := range ns {
		env := &config.Environment{Name: v}
		if k == "cicd" {
			env.IsCICD = true
		}
		if k == "dev" {
			app, err := applicationFromRepo(repoURL)
			if err != nil {
				return nil, err
			}
			env.Apps = []*config.Application{app}
		}
		envs = append(envs, env)
	}
	envs = append(envs, &config.Environment{Name: prefix + "argocd", IsArgoCD: true})
	sort.Sort(config.ByName(envs))
	return envs, nil
}

func applicationFromRepo(repoURL string) (*config.Application, error) {
	repo, err := repoFromURL(repoURL)
	if err != nil {
		return nil, err
	}
	return &config.Application{
		Name: repo,
		Services: []*config.Service{
			{
				Name:      repo + "-svc",
				SourceURL: repoURL,
			},
		},
	}, nil
}

func repoFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Path, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
}

func orgRepoFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Path, "/")
	orgRepo := strings.Join(parts[len(parts)-2:], "/")
	return strings.TrimSuffix(orgRepo, ".git"), nil
}
