package pipelines

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func TestBuildEventListener(t *testing.T) {
	m := &config.Manifest{
		Environments: []*config.Environment{
			{
				Name:   "test-cicd",
				IsCICD: true,
			},
			testEnv(testService()),
		},
	}
	cicdPath := filepath.Join("environments", "test-cicd")
	gitOpsRepo := "http://github.com/org/gitops.git"
	got, err := buildEventListenerResources(gitOpsRepo, m)
	assertNoError(t, err)
	want := res.Resources{
		getEventListenerPath(cicdPath): eventlisteners.CreateELFromTriggers("test-cicd", saName, fakeTiggers(t, m, gitOpsRepo)),
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("resources didn't match:%s\n", diff)
	}
}

func TestBuildEventListenerWithServiceWithNoURL(t *testing.T) {
	m := &config.Manifest{
		Environments: []*config.Environment{
			{
				Name:   "test-cicd",
				IsCICD: true,
			},
			testEnv(testService()),
		},
	}
	cicdPath := filepath.Join("environments", "test-cicd")
	gitOpsRepo := "http://github.com/org/gitops.git"
	got, err := buildEventListenerResources(gitOpsRepo, m)
	assertNoError(t, err)
	want := res.Resources{
		getEventListenerPath(cicdPath): eventlisteners.CreateELFromTriggers("test-cicd", saName, fakeTiggers(t, m, gitOpsRepo)),
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("resources didn't match:%s\n", diff)
	}
}

func TestBuildEventListenerWithNoGitOpsURL(t *testing.T) {
	m := &config.Manifest{
		Environments: []*config.Environment{
			{
				Name:   "test-cicd",
				IsCICD: true,
			},
			testEnv(testService()),
		},
	}
	got, err := buildEventListenerResources("", m)
	assertNoError(t, err)

	want := res.Resources{}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("resources didn't match:%s\n", diff)
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
			defaultPipelines,
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

func fakeTiggers(t *testing.T, m *config.Manifest, gitOpsRepo string) []v1alpha1.EventListenerTrigger {
	triggers := []v1alpha1.EventListenerTrigger{}
	devEnv := m.GetEnvironment("test-dev")
	cicdEnv, err := m.GetCICDEnvironment()
	assertNoError(t, err)
	devCITrigger, err := createCITrigger(gitOpsRepo, devEnv, testService())
	assertNoError(t, err)
	ciTrigger, err := createCITrigger(gitOpsRepo, cicdEnv, nil)
	assertNoError(t, err)
	cdTrigger, err := createCDTrigger(gitOpsRepo, cicdEnv, nil)
	assertNoError(t, err)
	triggers = append(triggers, devCITrigger, ciTrigger, cdTrigger)
	return triggers
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

func testEnv(svc *config.Service) *config.Environment {
	return &config.Environment{
		Name:      "test-dev",
		Pipelines: testPipelines("test"),
		Services: []*config.Service{
			svc,
		},
		Apps: []*config.Application{
			{
				Name: "test-app",
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
