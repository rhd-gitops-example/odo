package scm

import (
	"net/url"
	"strings"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

const (
	gitlabCIDryRunFilters = "header.match('X-Gitlab-Event','Merge Request Hook') && body.object_attributes.state == 'opened' && body.project.path_with_namespace == '%s'  && body.project.default_branch == body.object_attributes.target_branch"
	gitlabCDDeployFilters = "header.match('X-Gitlab-Event','Push Hook') && body.project.path_with_namespace == '%s' && body.ref.endsWith(body.project.default_branch)"
	gitlabPRBindingName   = "gitlab-pr-binding"
	gitlabPushBindingName = "gitlab-push-binding"
	gitlabType            = "gitlab"
)

type gitlab struct {
	repository
}

type gitlabSpec struct {
}

func init() {
	gits[gitlabType] = newGitLab
}

func newGitLab(rawURL string) (Repository, error) {
	path, err := processRawURL(rawURL, proccessGitLabPath)
	if err != nil {
		return nil, err
	}
	return &gitlab{repository{url: rawURL, path: path, spec: &gitlabSpec{}}}, nil
}

func proccessGitLabPath(parsedURL *url.URL) (string, error) {
	components, err := splitRepositoryPath(parsedURL)
	if err != nil {
		return "", err
	}
	if len(components) < 2 {
		return "", invalidRepoPathError(gitlabType, parsedURL.Path)
	}
	path := strings.Join(components, "/")
	return path, nil
}

func (r *gitlab) PRBindingName() string {
	return r.spec.prBindingName()
}

func (r *gitlab) PushBindingName() string {
	return r.spec.pushBindingName()
}

func (r *gitlabSpec) prBindingName() string {
	return gitlabPRBindingName
}

func (r *gitlabSpec) pushBindingName() string {
	return gitlabPushBindingName
}

func (r *gitlabSpec) prBindingParams() []triggersv1.Param {
	return []triggersv1.Param{
		createBindingParam("gitref", "$(body.object_attributes.source_branch)"),
		createBindingParam("gitsha", "$(body.object_attributes.last_commit.id)"),
		createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
		createBindingParam("fullname", "$(body.project.path_with_namespace)"),
	}
}

func (r *gitlabSpec) pushBindingParams() []triggersv1.Param {
	return []triggersv1.Param{
		createBindingParam("gitref", "$(body.ref)"),
		createBindingParam("gitsha", "$(body.after)"),
		createBindingParam("gitrepositoryurl", "$(body.project.git_http_url)"),
	}
}

func (r *gitlabSpec) ciDryRunFilters() string {
	return gitlabCIDryRunFilters
}

func (r *gitlabSpec) cdDeployFilters() string {
	return gitlabCDDeployFilters
}

func (r *gitlabSpec) eventInterceptor(secretNamespace, secretName string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitLab: &triggersv1.GitLabInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  webhookSecretKey,
				Namespace:  secretNamespace,
			},
		},
	}
}
