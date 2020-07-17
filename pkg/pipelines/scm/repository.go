package scm

import (
	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/triggers"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

const (
	webhookSecretKey = "webhook-secret-key"
)

var (
	gits = make(map[string]func(string) (Repository, error))
)

type repository struct {
	url  string
	path string // Repository path eg: (org/.../repo)
	spec triggerSpec
}

type triggerSpec interface {
	pushBindingParams() []triggersv1.Param
	pushEventFilters() string
	eventInterceptor(secretNamespace, secretName string) *triggersv1.EventInterceptor
	pushBindingName() string
}

// NewRepository returns a suitable Repository instance
// based on the driver name (github,gitlab,etc)
func NewRepository(url string) (Repository, error) {
	name, err := GetDriverName(url)
	if err != nil {
		return nil, err
	}

	git := gits[name]
	if git == nil {
		return nil, unsupportedGitTypeError(name)
	}

	return git(url)
}

func (r *repository) CreatePushBinding(ns string) (triggersv1.TriggerBinding, string) {
	return triggersv1.TriggerBinding{
		TypeMeta:   triggers.TriggerBindingTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName(ns, r.spec.pushBindingName())),
		Spec: triggersv1.TriggerBindingSpec{
			Params: r.spec.pushBindingParams(),
		},
	}, r.spec.pushBindingName()
}

func (r *repository) CreatePushTrigger(name, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return r.createTrigger(name, r.spec.pushEventFilters(),
		template, bindings,
		r.spec.eventInterceptor(secretNS, secretName))
}

func (r *repository) createTrigger(name, filters, template string, bindings []string, interceptor *triggersv1.EventInterceptor) triggersv1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			interceptor,
			createEventInterceptor(filters, r.path),
		},
		Bindings: createBindings(bindings),
		Template: createListenerTemplate(template),
	}
}

func (r *repository) PushBindingName() string {
	return r.spec.pushBindingName()
}

// URL returns the URL of the GitHub repository
func (r *repository) URL() string {
	return r.url
}
