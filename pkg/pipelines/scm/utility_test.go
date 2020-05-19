package scm

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func TestCreateListenerBinding(t *testing.T) {
	validListenerBinding := triggersv1.EventListenerBinding{
		Name: "sample",
	}
	listenerBinding := createListenerBinding("sample")
	if diff := cmp.Diff(validListenerBinding, *listenerBinding); diff != "" {
		t.Fatalf("createListenerBinding() failed:\n%s", diff)
	}
}

func TestCreateListenerTemplate(t *testing.T) {
	validListenerTemplate := triggersv1.EventListenerTemplate{
		Name: "sample",
	}
	listenerTemplate := createListenerTemplate("sample")
	if diff := cmp.Diff(validListenerTemplate, listenerTemplate); diff != "" {
		t.Fatalf("createListenerTemplate() failed:\n%s", diff)
	}
}

func TestCreateListenerTrigger(t *testing.T) {
	validListenerTrigger := triggersv1.EventListenerTrigger{
		Name: "sampleName",
		Interceptors: []*triggersv1.EventInterceptor{
			{
				CEL: &triggersv1.CELInterceptor{
					Filter: "sampleFilter sample",
				},
			},
			{
				GitHub: &triggersv1.GitHubInterceptor{
					SecretRef: &triggersv1.SecretRef{
						SecretName: "test",
						SecretKey:  WebhookSecretKey,
					},
				},
			},
		},
		Bindings: []*triggersv1.EventListenerBinding{
			{
				Name: "sampleBindingName",
			},
		},
		Template: triggersv1.EventListenerTemplate{
			Name: "sampleTemplateName",
		},
	}
	repo, err := NewRepository("http://github.com/org/test")
	if err != nil {
		t.Fatal(err)
	}
	listenerTrigger := createListenerTrigger(repo, "sampleName", "sampleFilter %s", "sample", "test", "", "sampleTemplateName", []string{"sampleBindingName"})
	if diff := cmp.Diff(validListenerTrigger, listenerTrigger); diff != "" {
		t.Fatalf("createListenerTrigger() failed:\n%s", diff)
	}
}

func TestCreateEventInterceptor(t *testing.T) {
	validEventInterceptor := triggersv1.EventInterceptor{
		CEL: &triggersv1.CELInterceptor{
			Filter: "sampleFilter sample",
		},
	}
	eventInterceptor := createEventInterceptor("sampleFilter %s", "sample")
	if diff := cmp.Diff(validEventInterceptor, *eventInterceptor); diff != "" {
		t.Fatalf("createEventInterceptor() failed:\n%s", diff)
	}
}
