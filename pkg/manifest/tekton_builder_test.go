package manifest

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/openshift/odo/pkg/manifest/eventlisteners"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func TestBuildEventListener(t *testing.T) {
	m := &config.Manifest{
		[]*config.Environment{
			&config.Environment{
				Name:   "test-cicd",
				IsCICD: true,
			},
			testEnv(),
		},
	}
	cicdPath := filepath.Join("environments", "test-cicd")
	gitOpsRepo := "http://github.com/org/gitops.git"
	got, err := buildEventlistenerResources(gitOpsRepo, m)
	assertNoError(t, err)
	triggers := fakeTriggers("org/gitops", "test-cicd", testService())

	want := res.Resources{
		getEventListenerPath(cicdPath): eventlisteners.CreateELFromTriggers("test-cicd", triggers),
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("resources didn't match:%s\n", diff)
	}
}

func fakeTriggers(gitopsRepo string, cicdNs string, svc *config.Service) []v1alpha1.EventListenerTrigger {
	return []v1alpha1.EventListenerTrigger{
		eventlisteners.CreateListenerTrigger(triggerName(svc.Name), eventlisteners.StageCIDryRunFilters, "org/test", "test-ci-binding", "test-ci-template", svc.Webhook.Secret.Name, svc.Webhook.Secret.Namespace),
		eventlisteners.CreateListenerTrigger("ci-dryrun-from-pr", eventlisteners.StageCIDryRunFilters, gitopsRepo, "github-pr-binding", "app-ci-template", eventlisteners.GitOpsWebhookSecret, cicdNs),
	}
}

func testService() *config.Service {
	return &config.Service{
		Name:      "test-svc",
		SourceURL: "http://github.com/org/test.git",
		Webhook: &config.Webhook{
			Secret: &config.Secret{
				Name:      "webhook-secret",
				Namespace: "webhook-ns",
			},
		},
	}
}

func testEnv() *config.Environment {
	return &config.Environment{
		Name:      "test-dev",
		Pipelines: testPipelines("test"),
		Apps: []*config.Application{
			&config.Application{
				Name: "test-app",
				Services: []*config.Service{
					testService(),
				},
			},
		},
	}
}

func testPipelines(name string) *config.Pipelines {
	return &config.Pipelines{
		Integration: &config.TemplateBinding{
			Template: fmt.Sprintf("%s-ci-template", name),
			Binding:  fmt.Sprintf("%s-ci-binding", name),
		},
	}
}

func TestGetPipelines(t *testing.T) {
	tests := []struct {
		desc string
		env  *config.Environment
		svc  *config.Service
		want *config.Pipelines
	}{
		{
			"Pipelines are provided by environment",
			&config.Environment{
				Name:      "test-env",
				Pipelines: testPipelines("env"),
			},
			&config.Service{
				Name: "test-svc",
			},
			testPipelines("env"),
		},
		{
			"Pipelines are provided by service",
			&config.Environment{
				Name: "test-env",
			},
			&config.Service{
				Name:      "test-service",
				Pipelines: testPipelines("svc"),
			},
			testPipelines("svc"),
		},
		{
			"Default pipelines are used",
			&config.Environment{
				Name: "test-env",
			},
			&config.Service{
				Name: "test-service",
			},
			defaultPipeline(),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(rt *testing.T) {
			got := getPipelines(test.env, test.svc)
			if diff := cmp.Diff(got, test.want); diff != "" {
				rt.Fatalf("getPipelines() failed:\n%v", diff)
			}
		})
	}
}
