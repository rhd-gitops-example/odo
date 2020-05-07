package pipelines

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift/odo/pkg/pipelines/argocd"
	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/meta"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/openshift/odo/pkg/pipelines/secrets"
)

func TestServiceResourcesWithCICD(t *testing.T) {
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

	fakeFs := ioutils.NewMapFilesystem()
	m := buildManifest(true, false)
	hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName("cicd", "github-webhook-secret-test-svc"), "123", eventlisteners.WebhookSecretKey)
	assertNoError(t, err)

	want := res.Resources{
		"environments/test-dev/services/test/base/kustomization.yaml":                     &res.Kustomization{Bases: []string{"./config"}},
		"environments/test-dev/services/test/kustomization.yaml":                          &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/services/test/overlays/kustomization.yaml":                 &res.Kustomization{Bases: []string{"../base"}},
		"environments/cicd/base/pipelines/03-secrets/github-webhook-secret-test-svc.yaml": hookSecret,
		"pipelines.yaml": &config.Manifest{
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							Services: []*config.Service{
								{
									Name:      "svc",
									SourceURL: "https://github.com/myproject/test-svc",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc-svc",
											Namespace: "cicd",
										},
									},
								},
								{
									Name:      "test",
									SourceURL: "http://github.com/org/test",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc",
											Namespace: "cicd",
										},
									},
								},
							},
						},
					},
				},
				{Name: "cicd", IsCICD: true},
			},
		},
	}

	got, err := serviceResources(m, fakeFs, "http://github.com/org/test", "test-dev", "test-app", "test", "123", pipelinesFile)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		_, ok := want[k]
		return !ok
	})); diff != "" {
		t.Fatalf("serviceResources() failed: %v", diff)
	}
}

func TestServiceResourcesWithoutCICD(t *testing.T) {
	fakeFs := ioutils.NewMapFilesystem()
	m := buildManifest(false, false)
	want := res.Resources{
		"environments/test-dev/services/test/base/kustomization.yaml":     &res.Kustomization{Bases: []string{"./config"}},
		"environments/test-dev/services/test/kustomization.yaml":          &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/services/test/overlays/kustomization.yaml": &res.Kustomization{Bases: []string{"../base"}},
		"pipelines.yaml": &config.Manifest{
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							Services: []*config.Service{
								{
									Name:      "svc",
									SourceURL: "https://github.com/myproject/test-svc",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc-svc",
											Namespace: "cicd",
										},
									},
								},
								{
									Name:      "test",
									SourceURL: "http://github.com/org/test",
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := serviceResources(m, fakeFs, "http://github.com/org/test", "test-dev", "test-app", "test", "123", pipelinesFile)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		_, ok := want[k]
		return !ok
	})); diff != "" {
		t.Fatalf("serviceResources() failed: %v", diff)
	}
}

func TestAddServiceWithoutApp(t *testing.T) {
	fakeFs := ioutils.NewMapFilesystem()
	m := buildManifest(false, false)
	want := res.Resources{
		"environments/test-dev/apps/new-app/base/kustomization.yaml":                      &res.Kustomization{Bases: []string{"../../../services/test"}},
		"environments/test-dev/apps/new-app/overlays/kustomization.yaml":                  &res.Kustomization{Bases: []string{"../base"}},
		"environments/test-dev/apps/new-app/kustomization.yaml":                           &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/services/test/base/kustomization.yaml":                     &res.Kustomization{Bases: []string{"./config"}},
		"environments/test-dev/services/test/kustomization.yaml":                          &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/services/test/overlays/kustomization.yaml":                 &res.Kustomization{Bases: []string{"../base"}},
		"environments/cicd/base/pipelines/03-secrets/github-webhook-secret-test-svc.yaml": nil,
		"pipelines.yaml": &config.Manifest{
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							Services: []*config.Service{
								{
									Name:      "svc",
									SourceURL: "https://github.com/myproject/test-svc",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc-svc",
											Namespace: "cicd",
										},
									},
								},
							},
						},
						{
							Name: "new-app",
							Services: []*config.Service{
								{Name: "test", SourceURL: "http://github.com/org/test"},
							},
						},
					},
				},
			},
		},
	}

	got, err := serviceResources(m, fakeFs, "http://github.com/org/test", "test-dev", "new-app", "test", "123", pipelinesFile)
	assertNoError(t, err)
	for w := range want {
		if diff := cmp.Diff(got[w], want[w]); diff != "" {
			t.Fatalf("serviceResources() failed: %v", diff)
		}
	}
}

func TestAddService(t *testing.T) {
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

	fakeFs := ioutils.NewMapFilesystem()
	outputPath := afero.GetTempDir(fakeFs, "test")
	manifestPath := filepath.Join(outputPath, pipelinesFile)
	m := buildManifest(true, true)
	b, err := yaml.Marshal(m)
	assertNoError(t, err)
	err = afero.WriteFile(fakeFs, manifestPath, b, 0644)
	assertNoError(t, err)
	wantedPaths := []string{
		"environments/test-dev/apps/new-app/base/kustomization.yaml",
		"environments/test-dev/apps/new-app/overlays/kustomization.yaml",
		"environments/test-dev/apps/new-app/kustomization.yaml",
		"environments/test-dev/services/test/base/kustomization.yaml",
		"environments/test-dev/services/test/overlays/kustomization.yaml",
		"environments/test-dev/services/test/kustomization.yaml",
		"environments/cicd/base/pipelines/03-secrets/github-webhook-secret-test-svc.yaml",
		"environments/cicd/base/pipelines/kustomization.yaml",
		"pipelines.yaml",
		"environments/argocd/config/test-dev-test-app-app.yaml",
		"environments/argocd/config/test-dev-new-app-app.yaml",
	}
	err = AddService("http://github.com/org/test", "test-dev", "new-app", "test", "123", manifestPath, fakeFs)
	assertNoError(t, err)
	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(outputPath, path))
		})
	}
}

func TestServiceWithArgoCD(t *testing.T) {
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

	fakeFs := ioutils.NewMapFilesystem()
	m := buildManifest(true, true)
	want := res.Resources{
		"pipelines.yaml": &config.Manifest{
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							Services: []*config.Service{
								{
									Name:      "svc",
									SourceURL: "https://github.com/myproject/test-svc",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc-svc",
											Namespace: "cicd",
										},
									},
								},
								{
									Name:      "test",
									SourceURL: "http://github.com/org/test",
									Webhook: &config.Webhook{
										Secret: &config.Secret{
											Name:      "github-webhook-secret-test-svc",
											Namespace: "cicd",
										},
									},
								},
							},
						},
					},
				},
				{Name: "cicd", IsCICD: true},
				{Name: "argocd", IsArgoCD: true},
			},
		},
	}
	argo, err := argocd.Build("argocd", "http://github.com/org/test", m)
	assertNoError(t, err)
	want = res.Merge(argo, want)
	got, err := serviceResources(m, fakeFs, "http://github.com/org/test", "test-dev", "test-app", "test", "123", pipelinesFile)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		_, ok := want[k]
		return !ok
	})); diff != "" {
		t.Fatalf("serviceResources() failed: %v", diff)
	}
}

func buildManifest(withCICD, withArgoCD bool) *config.Manifest {
	cfg := &config.Manifest{
		GitOpsURL: "http://github.com/org/test",
		Environments: []*config.Environment{
			{
				Name: "test-dev",
				Apps: []*config.Application{
					{
						Name: "test-app",
						Services: []*config.Service{
							{
								Name:      "svc",
								SourceURL: "https://github.com/myproject/test-svc",
								Webhook: &config.Webhook{
									Secret: &config.Secret{
										Name:      "github-webhook-secret-test-svc-svc",
										Namespace: "cicd",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if withCICD == true {
		cfg.Environments = append(cfg.Environments, &config.Environment{
			Name:   "cicd",
			IsCICD: true,
		})
	}
	if withArgoCD == true {
		cfg.Environments = append(cfg.Environments, &config.Environment{
			Name:     "argocd",
			IsArgoCD: true,
		})
	}
	return cfg
}
