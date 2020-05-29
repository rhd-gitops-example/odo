package scm

import (
	"net/url"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/meta"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

var (
	triggerBindingTypeMeta = meta.TypeMeta("TriggerBinding", "triggers.tekton.dev/v1alpha1")
)

// Filters for interceptors
const (
	githubCIDryRunFilters = "(header.match('X-GitHub-Event', 'pull_request') && body.action == 'opened' || body.action == 'synchronize') && body.pull_request.head.repo.full_name == '%s'"

	githubCDDeployFilters = "(header.match('X-GitHub-Event', 'push') && body.repository.full_name == '%s') && body.ref.startsWith('refs/heads/master')"

	webhookSecretKey = "webhook-secret-key"

	githubPRBindingName   = "github-pr-binding"
	githubPushBindingName = "github-push-binding"
)

func init() {
	supportedTriggerBindings = append(supportedTriggerBindings, githubPRBindingName, githubPushBindingName)
}

// GitHubRepository represents a service on a GitHub repo
type GitHubRepository struct {
	repository
}

// NewGitHubRepository returns an instance of GitHubRepository
func NewGitHubRepository(rawURL string) (*GitHubRepository, error) {
	parsedURL, path, err := processRawURL(rawURL, func(parsedURL *url.URL) (string, error) {
		components, err := splitRepositoryPath(parsedURL)
		if err != nil {
			return "", err
		}
		path := components[0] + "/" + strings.TrimSuffix(components[1], ".git")
		return path, nil
	})
	if err != nil {
		return nil, err
	}
	return &GitHubRepository{repository{url: parsedURL, path: path}}, nil
}

// CreatePRBinding returns a TriggerBinding for GitHub PullRequest hooks.
func (repo *GitHubRepository) CreatePRBinding(ns string) (triggersv1.TriggerBinding, string) {
	return repo.createBinding(ns, githubPRBindingName, []triggersv1.Param{
		createBindingParam("gitref", "$(body.pull_request.head.ref)"),
		createBindingParam("gitsha", "$(body.pull_request.head.sha)"),
		createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
		createBindingParam("fullname", "$(body.repository.full_name)"),
	}), githubPRBindingName
}

// CreatePushBinding returns a TriggerBinding for GitHub Push hooks.
func (repo *GitHubRepository) CreatePushBinding(ns string) (triggersv1.TriggerBinding, string) {
	return repo.createBinding(ns, githubPushBindingName, []triggersv1.Param{
		createBindingParam("gitref", "$(body.ref)"),
		createBindingParam("gitsha", "$(body.head_commit.id)"),
		createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
	}), githubPushBindingName
}

// URL returns the URL of the GitHub repository
func (repo *GitHubRepository) URL() string {
	return repo.String()
}

// CreateCITrigger creates a CI eventlistener trigger for GitHub
func (repo *GitHubRepository) CreateCITrigger(name, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return repo.createTrigger(name, githubCIDryRunFilters,
		template, bindings,
		repo.CreateInterceptor(secretName, secretNS))
}

// CreateCDTrigger creates a CD eventlistener trigger for GitHub
func (repo *GitHubRepository) CreateCDTrigger(name, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return repo.repository.createTrigger(name, githubCDDeployFilters,
		template, bindings,
		repo.CreateInterceptor(secretName, secretNS))
}

// CreateInterceptor returns a GitHub event interceptor
func (repo *GitHubRepository) CreateInterceptor(secretName, secretNs string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitHub: &triggersv1.GitHubInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  webhookSecretKey,
				Namespace:  secretNs,
			},
		},
	}
}
