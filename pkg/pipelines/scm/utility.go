package scm

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func GetDriverName(rawURL string) (string, error) {

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if s := strings.TrimSuffix(u.Host, ".com"); s != u.Host {
		return strings.ToLower(s), nil
	}

	if s := strings.TrimSuffix(u.Host, ".org"); s != u.Host {
		return strings.ToLower(s), nil
	}

	return "", errors.New("unknown Git server: " + u.Host)
}

func GetRepoName(rawURL string) (string, error) {

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

	if len(components) != 2 {
		return "", errors.New("failed to get Git repo: " + u.Path)
	}

	components[1] = strings.TrimSuffix(components[1], ".git")

	for _, s := range components {
		if strings.Index(s, ".") != -1 {
			return "", errors.New("failed to get Git repo: " + u.Path)
		}
	}

	return components[0] + "/" + components[1], nil
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
