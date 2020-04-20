package manifest

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/openshift/odo/pkg/manifest/deployment"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/openshift/odo/pkg/manifest/secrets"
)

const (
	testSvcRepo    = "https://github.com/my-org/http-api.git"
	testGitOpsRepo = "https://github.com/my-org/gitops.git"
)

func TestBootstrapManifest(t *testing.T) {
	defer func(f secrets.PublicKeyFunc) {
		secrets.DefaultPublicKeyFunc = f
	}(secrets.DefaultPublicKeyFunc)

	secrets.DefaultPublicKeyFunc = func() (*rsa.PublicKey, error) {
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Fatalf("failed to generate a private RSA key: %s", err)
		}
		return &key.PublicKey, nil
	}

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

	want := res.Resources{
		"environments/tst-dev/services/http-api-svc/base/config/100-deployment.yaml": deployment.Create("tst-dev", "http-api-svc", bootstrapImage, deployment.ContainerPort(80)),

		"environments/tst-dev/services/http-api-svc/base/config/200-service.yaml":   createBootstrapService("tst-dev", "http-api-svc"),
		"environments/tst-dev/services/http-api-svc/base/config/kustomization.yaml": &res.Kustomization{Resources: []string{"100-deployment.yaml", "200-service.yaml"}},
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
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name: "gitops-webhook-secret",
										},
									},
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
		_, ok := want[k]
		t.Logf("got %s\n", k)
		return !ok
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
				Webhook: &config.Webhook{
					Secret: &config.Secret{
						Name: "test-svc-webhook-secret"},
				},
			},
		},
	}

	got, err := applicationFromRepo(testSvcRepo, "test-svc-webhook-secret")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}

}
