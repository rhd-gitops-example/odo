package scm

import (
	"fmt"
	"net/url"
	"strings"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func getDriverName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if s := strings.TrimSuffix(u.Host, ".com"); s != u.Host {
		return strings.ToLower(s), nil
	}
	return "", invalidRepoTypeError(rawURL)
}

func getRepoName(rawURL string) (string, error) {
	var components []string
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	for _, s := range strings.Split(u.Path, "/") {
		if s != "" {
			components = append(components, s)
		}
	}
	if len(components) < 2 {
		return "", invalidRepoNameError(rawURL)
	}
	return components[0] + "/" + strings.TrimSuffix(components[1], ".git"), nil
}

func invalidRepoTypeError(repoURL string) error {
	return fmt.Errorf("unable to determine type of Git host from: %s", repoURL)
}

func invalidRepoNameError(repoURL string) error {
	return fmt.Errorf("unable to determine repo name from: %s", repoURL)
}

func createListenerTrigger(repo Repository, name, filter, repoName, secretName, secretNS, template string, bindings []string) triggersv1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			createEventInterceptor(filter, repoName),
			repo.CreateInterceptor(secretName, secretNS),
		},
		Bindings: createBindings(bindings),
		Template: createListenerTemplate(template),
	}
}

func createEventInterceptor(filter string, repoName string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		CEL: &triggersv1.CELInterceptor{
			Filter: fmt.Sprintf(filter, repoName),
		},
	}
}

func createListenerTemplate(name string) triggersv1.EventListenerTemplate {
	return triggersv1.EventListenerTemplate{
		Name: name,
	}
}

func createListenerBinding(name string) *triggersv1.EventListenerBinding {
	return &triggersv1.EventListenerBinding{
		Name: name,
	}
}

func createBindings(names []string) []*triggersv1.EventListenerBinding {
	bindings := make([]*triggersv1.EventListenerBinding, len(names))
	for i, name := range names {
		bindings[i] = createListenerBinding(name)
	}
	return bindings
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
