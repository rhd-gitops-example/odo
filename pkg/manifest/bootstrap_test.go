package manifest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
)

const (
	testSvcRepo    = "https://github.com/my-org/http-api.git"
	testGitOpsRepo = "https://github.com/my-org/gitops.git"
)

func TestBootstrapManifest(t *testing.T) {
	params := &BootstrapOptions{
		Prefix:              "tst-",
		GitOpsRepoURL:       testGitOpsRepo,
		GitOpsWebhookSecret: "123",
		AppRepoURL:          testSvcRepo,
		ImageRepo:           "image/repo",
	}

	r, err := bootstrapResources(params)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: pipeline output with names
	want := res.Resources{
		"manifest.yaml": &config.Manifest{
			Environments: []*config.Environment{
				{
					Pipelines: defaultPipelines,
					Name:      "tst-dev",
					Apps: []*config.Application{
						{
							Name: "http-api",
							Services: []*config.Service{
								{
									Name:      "http-api-svc",
									SourceURL: testSvcRepo,
								},
							},
						},
					},
				},
				{Name: "tst-stage"},
				{Name: "tst-cicd", IsCICD: true},
				{Name: "tst-argocd", IsArgoCD: true},
			},
		},
	}

	if diff := cmp.Diff(want, r, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		return k != "manifest.yaml"
	})); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}
}

func TestOrgRepoFromURL(t *testing.T) {
	want := "my-org/gitops"
	got, err := orgRepoFromURL(testGitOpsRepo)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("orgRepFromURL(%s) got %s, want %s", testGitOpsRepo, got, want)
	}
}

func TestApplicationFromRepo(t *testing.T) {
	want := &config.Application{
		Name: "http-api",
		Services: []*config.Service{
			{
				Name:      "http-api-svc",
				SourceURL: testSvcRepo,
			},
		},
	}

	got, err := applicationFromRepo(testSvcRepo)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}

}
