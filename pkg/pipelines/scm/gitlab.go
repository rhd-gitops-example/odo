package scm

import (
	"net/url"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
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
	url  *url.URL
	path string // GitLab repo path eg: (group/subgroup/../repo)
}

// NewGitLabRepository returns an instance of GitLabRepository
// NewGitLabRepository returns an instance of GitLabRepository
func NewGitLabRepository(rawURL string) (*GitLabRepository, error) {
	parsedURL, path, err := processRawURL(rawURL, func(parsedURL *url.URL) (string, err) {
		var components []string
		for _, s := range strings.Split(parsedURL.Path, "/") {
			if s != "" {
				components = append(components, s)
			}
		}
		if len(components) < 2 {
			return nil, invalidRepoPathError(rawURL)
		}
		components[len(components)-1] = strings.TrimSuffix(components[len(components)-1], ".git")
		path := strings.Join(components, "/")
	})
	if err != nil {
		return nil, err
	}
	return &GitLabRepository{url: parsedURL, path: path}, nil
}

// CreatePRBinding returns a TriggerBinding for GitLab merge request hooks
func (repo *GitLabRepository) CreatePRBinding(ns string) (triggersv1.TriggerBinding, string) {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, gitlabPRBindingName)),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				createBindingParam("gitref", "$(body.object_attributes.source_branch)"),
				createBindingParam("gitsha", "$(body.object_attributes.last_commit.id)"),
				createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
				createBindingParam("fullname", "$(body.project.path_with_namespace)"),
			},
		},
	}, gitlabPRBindingName

}

// CreatePushBinding returns a TriggerBinding for GitLab push hooks
func (repo *GitLabRepository) CreatePushBinding(ns string) (triggersv1.TriggerBinding, string) {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, gitlabPushBindingName)),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				createBindingParam("gitref", "$(body.ref)"),
				createBindingParam("gitsha", "$(body.after)"),
				createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
			},
		},
	}, gitlabPushBindingName
}

// URL returns the URL of the GitLab repository
func (repo *GitLabRepository) URL() string {
	return repo.url.String()
}

// CreateCITrigger creates a CI eventlistener trigger for GitLab
func (repo *GitLabRepository) CreateCITrigger(name, secretName, secretNS, template string, bindings []string) v1alpha1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			createEventInterceptor(gitlabCIDryRunFilters, repo.path),
			repo.CreateInterceptor(secretName, secretNS),
		},
		Bindings: createBindings(bindings),
		Template: createListenerTemplate(template),
	}
}

// CreateCDTrigger creates a CD eventlistener trigger for GitLab
func (repo *GitLabRepository) CreateCDTrigger(name, secretName, secretNS, template string, bindings []string) v1alpha1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			createEventInterceptor(gitlabCDDeployFilters, repo.path),
			repo.CreateInterceptor(secretName, secretNS),
		},
		Bindings: createBindings(bindings),
		Template: createListenerTemplate(template),
	}
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
