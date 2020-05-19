package github

import (
	"net/url"

	"github.com/openshift/odo/pkg/odo/cli/pipelines/scm/utility"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/meta"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

var (
	triggerBindingTypeMeta = meta.TypeMeta("TriggerBinding", "tekton.dev/v1alpha1")
)

const (
	githubPRBindingName   = "github-pr-binding"
	githubPushBindingName = "github-push-binding"
)

// GithubRepository represents a service on a Github repo
type GithubRepository struct {
	URL *url.URL
}

func NewRepository(rawURL string) (*GithubRepository, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &GithubRepository{URL: parsedURL}, nil
}

// CreatePRBinding returns a TriggerBinding for Github PullRequest hooks.
func (repo *GithubRepository) CreatePRBinding(ns string) (triggersv1.TriggerBinding, string) {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, githubPRBindingName)),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.pull_request.head.ref)"),
				createBindingParam("gitsha", "$(body.pull_request.head.sha)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
				createBindingParam("fullname", "$(body.repository.full_name)"),
			},
		},
	}, githubPRBindingName
}

// CreatePushBinding returns a TriggerBinding for Github Push hooks.
func (repo *GithubRepository) CreatePushBinding(ns string) (triggersv1.TriggerBinding, string) {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, githubPushBindingName)),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.ref)"),
				createBindingParam("gitsha", "$(body.head_commit.id)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
			},
		},
	}, githubPushBindingName
}

// GetURL returns the URL of the Github repository
func (repo *GithubRepository) GetURL() string {
	return repo.URL.String()
}

// CreateCITrigger creates a CI eventlistener trigger for Github
func (repo *GithubRepository) CreateCITrigger(name, secretName, secretNs, template string, bindings []string) (v1alpha1.EventListenerTrigger, error) {
	repoName, err := utility.GetRepoName(repo.GetURL())
	if err != nil {
		return v1alpha1.EventListenerTrigger{}, err
	}
	return eventlisteners.CreateListenerTrigger(repo, name, repo.GetCIFilters(), repoName, secretName, secretNs, template, bindings), nil
}

// CreateCDTrigger creates a CD eventlistener trigger for Github
func (repo *GithubRepository) CreateCDTrigger(name, secretName, secretNs, template string, bindings []string) (v1alpha1.EventListenerTrigger, error) {
	repoName, err := utility.GetRepoName(repo.GetURL())
	if err != nil {
		return v1alpha1.EventListenerTrigger{}, err
	}
	return eventlisteners.CreateListenerTrigger(repo, name, repo.GetCDFilters(), repoName, secretName, secretNs, template, bindings), nil
}

// CreateInterceptor returns a Github event interceptor
func (repo *GithubRepository) CreateInterceptor(secretName, secretNs string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitHub: &triggersv1.GitHubInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  eventlisteners.WebhookSecretKey,
				Namespace:  secretNs,
			},
		},
	}
}

// GetCIFilters returns the CEL interceptor filters
// for a Github CI event
func (repo *GithubRepository) GetCIFilters() string {
	return eventlisteners.StageCIDryRunFilters
}

// GetCDFilters returns the CEL event interceptor filters
// for a Github CD event
func (repo *GithubRepository) GetCDFilters() string {
	return eventlisteners.StageCDDeployFilters
}

func createBindingParam(name string, value string) pipelinev1.Param {
	return pipelinev1.Param{
		Name: name,
		Value: pipelinev1.ArrayOrString{
			StringVal: value,
			Type:      pipelinev1.ParamTypeString,
		},
	}
}
