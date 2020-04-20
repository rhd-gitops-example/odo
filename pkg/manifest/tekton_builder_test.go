package manifest

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
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
	elPatch, err := eventlistenerPatch(testEnv(), testService())
	assertNoError(t, err)
	got, err := buildEventlistenerResources(m)
	assertNoError(t, err)
	want := res.Resources{
		"environments/test-cicd/overlays/eventlistener_patches/test-svc_patch.yaml": elPatch,
		"environments/test-cicd/overlays/kustomization.yaml":                        elKustomiseTarget("test-cicd", "../base", []string{"test-svc"}),
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("resources didn't match:%s\n", diff)
	}
}

func testService() *config.Service {
	return &config.Service{
		Name:      "test-svc",
		SourceURL: "http://github.com/org/test.git",
		Webhook: &config.Webhook{
			Secret: &config.Secret{
				Name:      "webhook-secret",
				Namespace: "webhook-key",
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
			testPipelines("app"),
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
