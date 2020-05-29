package scm

import (
	"net/url"

	"github.com/openshift/odo/pkg/pipelines/meta"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

var (
	supportedTriggerBindings []string
)

type repository struct {
	url  *url.URL
	path string // Repository path eg: (org/.../repo)
}

// NewRepository returns a suitable Repository instance
// based on the driver name (github,gitlab,etc)
func NewRepository(rawURL string) (Repository, error) {
	repoType, err := getDriverName(rawURL)
	if err != nil {
		return nil, err
	}
	switch repoType {
	case "github":
		return NewGitHubRepository(rawURL)
	case "gitlab":
		return NewGitLabRepository(rawURL)
	}
	return nil, invalidRepoTypeError(rawURL)
}

// IsRepositoryBinding checks if a binding
// is a supported binding type
func IsRepositoryBinding(b string) bool {
	for _, binding := range supportedTriggerBindings {
		if binding == b {
			return true
		}
	}
	return false
}

func (r *repository) createBinding(ns, name string, params []triggersv1.Param) triggersv1.TriggerBinding {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, name)),
		Spec: triggersv1.TriggerBindingSpec{
			Params: params,
		},
	}
}

func (r *repository) createTrigger(name, filters, template string, bindings []string, interceptor *triggersv1.EventInterceptor) triggersv1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			createEventInterceptor(filters, r.path),
			interceptor,
		},
		Bindings: createBindings(bindings),
		Template: createListenerTemplate(template),
	}
}

func (r *repository) String() string {
	return r.url.String()
}
