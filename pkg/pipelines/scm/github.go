package scm

import (
	"net/url"
	"strings"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

const (
	githubCIDryRunFilters = "(header.match('X-GitHub-Event', 'pull_request') && body.action == 'opened' || body.action == 'synchronize') && body.pull_request.head.repo.full_name == '%s'"
	githubCDDeployFilters = "(header.match('X-GitHub-Event', 'push') && body.repository.full_name == '%s') && body.ref.startsWith('refs/heads/master')"
	githubPRBindingName   = "github-pr-binding"
	githubPushBindingName = "github-push-binding"
	githubType            = "github"
)

type github struct {
	repository
}

type githubSpec struct {
}

func init() {
	gits[githubType] = newGitHub
}

func newGitHub(rawURL string) (Repository, error) {
	path, err := processRawURL(rawURL, proccessGitHubPath)
	if err != nil {
		return nil, err
	}
	return &github{repository{url: rawURL, path: path, spec: &githubSpec{}}}, nil
}

func proccessGitHubPath(parsedURL *url.URL) (string, error) {
	components, err := splitRepositoryPath(parsedURL)
	if err != nil {
		return "", err
	}

	if len(components) != 2 {
		return "", invalidRepoPathError(githubType, parsedURL.Path)
	}
	path := strings.Join(components, "/")
	return path, nil
}

func (r *github) PRBindingName() string {
	return r.spec.prBindingName()
}

func (r *github) PushBindingName() string {
	return r.spec.pushBindingName()
}

func (r *githubSpec) prBindingName() string {
	return githubPRBindingName
}

func (r *githubSpec) pushBindingName() string {
	return githubPushBindingName
}

func (r *githubSpec) prBindingParams() []triggersv1.Param {
	return []triggersv1.Param{
		createBindingParam("gitref", "$(body.pull_request.head.ref)"),
		createBindingParam("gitsha", "$(body.pull_request.head.sha)"),
		createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
		createBindingParam("fullname", "$(body.repository.full_name)"),
	}
}

func (r *githubSpec) pushBindingParams() []triggersv1.Param {
	return []triggersv1.Param{
		createBindingParam("gitref", "$(body.ref)"),
		createBindingParam("gitsha", "$(body.head_commit.id)"),
		createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
	}
}

func (r *githubSpec) ciDryRunFilters() string {
	return githubCIDryRunFilters
}

func (r *githubSpec) cdDeployFilters() string {
	return githubCDDeployFilters
}

func (r *githubSpec) eventInterceptor(secretNamespace, secretName string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitHub: &triggersv1.GitHubInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  webhookSecretKey,
				Namespace:  secretNamespace,
			},
		},
	}
}
