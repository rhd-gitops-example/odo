package scm

import (
	"net/url"
	"strings"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

const (
	gitlabPRBindingName   = "gitlab-pr-binding"
	gitlabPushBindingName = "gitlab-push-binding"

	gitlabCIDryRunFilters = "header.match('X-Gitlab-Event','Merge Request Hook') && body.object_attributes.state == 'opened' && body.project.path_with_namespace == '%s'  && body.project.default_branch == body.object_attributes.target_branch"

	gitlabCDDeployFilters = "header.match('X-Gitlab-Event','Push Hook') && body.project.path_with_namespace == '%s' && body.ref.endsWith(body.project.default_branch)"
)

func init() {
	supportedTriggerBindings = append(supportedTriggerBindings, gitlabPRBindingName, gitlabPushBindingName)
}

// GitLabRepository represents a service on a GitLab repo
type GitLabRepository struct {
	repository
}

// NewGitLabRepository returns an instance of GitLabRepository
// NewGitLabRepository returns an instance of GitLabRepository
func NewGitLabRepository(rawURL string) (*GitLabRepository, error) {
	parsedURL, path, err := processRawURL(rawURL, func(parsedURL *url.URL) (string, error) {
		components, err := splitRepositoryPath(parsedURL)
		if err != nil {
			return "", err
		}
		components[len(components)-1] = strings.TrimSuffix(components[len(components)-1], ".git")
		path := strings.Join(components, "/")
		return path, nil
	})
	if err != nil {
		return nil, err
	}
	return &GitLabRepository{repository{url: parsedURL, path: path}}, nil
}

// CreatePRBinding returns a TriggerBinding for GitLab merge request hooks
func (repo *GitLabRepository) CreatePRBinding(ns string) (triggersv1.TriggerBinding, string) {
	return repo.createBinding(ns, gitlabPRBindingName, []triggersv1.Param{
		createBindingParam("gitref", "$(body.object_attributes.source_branch)"),
		createBindingParam("gitsha", "$(body.object_attributes.last_commit.id)"),
		createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
		createBindingParam("fullname", "$(body.project.path_with_namespace)"),
	}), gitlabPRBindingName
}

// CreatePushBinding returns a TriggerBinding for GitLab push hooks
func (repo *GitLabRepository) CreatePushBinding(ns string) (triggersv1.TriggerBinding, string) {
	return repo.createBinding(ns, gitlabPushBindingName, []triggersv1.Param{
		createBindingParam("gitref", "$(body.ref)"),
		createBindingParam("gitsha", "$(body.after)"),
		createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
	}), gitlabPushBindingName
}

// URL returns the URL of the GitLab repository
func (repo *GitLabRepository) URL() string {
	return repo.String()
}

// CreateCITrigger creates a CI eventlistener trigger for GitLab
func (repo *GitLabRepository) CreateCITrigger(name, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return repo.createTrigger(name, gitlabCIDryRunFilters,
		template, bindings,
		repo.CreateInterceptor(secretName, secretNS))
}

// CreateCDTrigger creates a CD eventlistener trigger for GitLab
func (repo *GitLabRepository) CreateCDTrigger(name, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return repo.createTrigger(name, gitlabCDDeployFilters,
		template, bindings,
		repo.CreateInterceptor(secretName, secretNS))
}

// CreateInterceptor returns a GitLab eventlistener
func (repo *GitLabRepository) CreateInterceptor(secretName, secretNs string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitLab: &triggersv1.GitLabInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  webhookSecretKey,
				Namespace:  secretNs,
			},
		},
	}
}
