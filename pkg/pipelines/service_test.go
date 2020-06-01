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
	"github.com/openshift/odo/pkg/pipelines/meta"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"github.com/spf13/afero"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/odo/pkg/pipelines/ioutils"
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
	hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName("cicd", "github-webhook-secret-test"), "123", eventlisteners.WebhookSecretKey)
	assertNoError(t, err)

	want := res.Resources{
		"config/cicd/base/pipelines/03-secrets/github-webhook-secret-test.yaml": hookSecret,
		"environments/test-dev/apps/test-app/base/kustomization.yaml":           &res.Kustomization{Bases: []string{"../../../services/test-svc", "../../../services/test"}},
		"environments/test-dev/apps/test-app/kustomization.yaml":                &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/apps/test-app/overlays/kustomization.yaml":       &res.Kustomization{Bases: []string{"../base"}},
		"pipelines.yaml": &config.Manifest{
			Config: &config.Config{
				CICD: &config.Cicd{
					Namespace: "cicd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
								"test",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "github-webhook-secret-test-svc",
									Namespace: "cicd",
								},
							},
						},
						{
							Name:      "test",
							SourceURL: "http://github.com/org/test",
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "github-webhook-secret-test",
									Namespace: "cicd",
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := serviceResources(m, fakeFs, &AddServiceParameters{
		AppName:       "test-app",
		EnvName:       "test-dev",
		GitRepoURL:    "http://github.com/org/test",
		Manifest:      pipelinesFile,
		WebhookSecret: "123",
		ServiceName:   "test",
	})
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
	m := buildManifest(false, true)
	want := res.Resources{
		"environments/test-dev/apps/test-app/base/kustomization.yaml":     &res.Kustomization{Bases: []string{"../../../services/test-svc", "../../../services/test"}},
		"environments/test-dev/apps/test-app/kustomization.yaml":          &res.Kustomization{Bases: []string{"overlays"}},
		"environments/test-dev/apps/test-app/overlays/kustomization.yaml": &res.Kustomization{Bases: []string{"../base"}},
		"environments/test-dev/env/base/kustomization.yaml":               &res.Kustomization{Resources: []string{"test-dev-environment.yaml"}},
		"pipelines.yaml": &config.Manifest{
			Config: &config.Config{
				Argo: &config.Argo{
					Namespace: "argocd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
								"test",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
							Webhook: &config.Webhook{
								Secret: &config.Secret{Name: "github-webhook-secret-test-svc", Namespace: "cicd"},
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
	}

	got, err := serviceResources(m, fakeFs, &AddServiceParameters{
		AppName:       "test-app",
		EnvName:       "test-dev",
		GitRepoURL:    "http://github.com/org/test",
		Manifest:      pipelinesFile,
		WebhookSecret: "123",
		ServiceName:   "test",
	})
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
							Name:        "test-app",
							ServiceRefs: []string{"test-svc"},
						},
						{
							Name:        "new-app",
							ServiceRefs: []string{"test"},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "github-webhook-secret-test-svc",
									Namespace: "cicd",
								},
							},
						},
						{Name: "test", SourceURL: "http://github.com/org/test"},
					},
				},
			},
		},
	}

	got, err := serviceResources(m, fakeFs, &AddServiceParameters{
		AppName:       "new-app",
		EnvName:       "test-dev",
		GitRepoURL:    "http://github.com/org/test",
		Manifest:      pipelinesFile,
		WebhookSecret: "123",
		ServiceName:   "test",
	})
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
		"environments/test-dev/apps/test-app/base/kustomization.yaml",
		"environments/test-dev/apps/test-app/overlays/kustomization.yaml",
		"environments/test-dev/apps/test-app/kustomization.yaml",
		"environments/test-dev/services/test/base/kustomization.yaml",
		"environments/test-dev/services/test/overlays/kustomization.yaml",
		"environments/test-dev/services/test/kustomization.yaml",
		"config/cicd/base/pipelines/03-secrets/github-webhook-secret-test.yaml",
		"config/cicd/base/pipelines/kustomization.yaml",
		"pipelines.yaml",
		// "config/argocd/config/test-dev-test-app-app.yaml",
		// "config/argocd/config/kustomization.yaml",
	}
	err = AddService(&AddServiceParameters{
		AppName:       "test-app",
		EnvName:       "test-dev",
		GitRepoURL:    "http://github.com/org/test",
		Manifest:      manifestPath,
		WebhookSecret: "123",
		ServiceName:   "test",
	}, fakeFs)
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
			Config: &config.Config{
				CICD: &config.Cicd{
					Namespace: "cicd",
				},
				Argo: &config.Argo{
					Namespace: "argocd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name:        "test-app",
							ServiceRefs: []string{"test-svc", "test"},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "github-webhook-secret-test-svc",
									Namespace: "cicd",
								},
							},
						},
						{
							Name:      "test",
							SourceURL: "http://github.com/org/test",
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "github-webhook-secret-test",
									Namespace: "cicd",
								},
							},
						},
					},
				},
			},
		},
	}
	argo, err := argocd.Build("argocd", "http://github.com/org/test", m)
	assertNoError(t, err)
	want = res.Merge(argo, want)
	got, err := serviceResources(m, fakeFs, &AddServiceParameters{
		AppName:       "test-app",
		EnvName:       "test-dev",
		GitRepoURL:    "http://github.com/org/test",
		Manifest:      pipelinesFile,
		WebhookSecret: "123",
		ServiceName:   "test",
	})
	assertNoError(t, err)
	if diff := cmp.Diff(got, want, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		_, ok := want[k]
		return !ok
	})); diff != "" {
		t.Fatalf("serviceResources() failed: %v", diff)
	}
}

func buildManifest(withCICD, withArgoCD bool) *config.Manifest {
	if !withCICD && withArgoCD {
		return &config.Manifest{
			Config: &config.Config{
				Argo: &config.Argo{
					Namespace: "argocd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
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
		}

	}
	if withCICD && withArgoCD {
		return &config.Manifest{
			Config: &config.Config{
				CICD: &config.Cicd{
					Namespace: "cicd",
				},
				Argo: &config.Argo{
					Namespace: "argocd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
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
		}
	}
	if withCICD && !withArgoCD {
		return &config.Manifest{
			Config: &config.Config{
				CICD: &config.Cicd{
					Namespace: "cicd",
				},
			},
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
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
		}
	}
	if !withCICD && !withArgoCD {
		return &config.Manifest{
			GitOpsURL: "http://github.com/org/test",
			Environments: []*config.Environment{
				{
					Name: "test-dev",
					Apps: []*config.Application{
						{
							Name: "test-app",
							ServiceRefs: []string{
								"test-svc",
							},
						},
					},
					Services: []*config.Service{
						{
							Name:      "test-svc",
							SourceURL: "https://github.com/myproject/test-svc",
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
		}
	}
	return nil
}

func TestCreateSvcImageBinding(t *testing.T) {
	CICD := &config.Cicd{
		Namespace: "cicd",
	}
	env := &config.Environment{
		Name: "new-env",
	}
	bindingName, bindingFilename, resources := createSvcImageBinding(CICD, env, "new-svc", "quay.io/user/app", false)

	if diff := cmp.Diff(bindingName, "new-env-new-svc-binding"); diff != "" {
		t.Errorf("bindingName failed: %v", diff)
	}

	if diff := cmp.Diff(bindingFilename, "06-bindings/new-env-new-svc-binding.yaml"); diff != "" {
		t.Errorf("bindingFilename failed: %v", diff)
	}

	triggerBinding := triggersv1.TriggerBinding{
		TypeMeta:   v1.TypeMeta{Kind: "TriggerBinding", APIVersion: "triggers.tekton.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Name: "new-env-new-svc-binding", Namespace: "cicd"},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "imageRepo",
					Value: "quay.io/user/app",
				},
				{
					Name:  "tlsVerify",
					Value: "false",
				},
			},
		},
	}

	wantResources := res.Resources{"config/cicd/base/pipelines/06-bindings/new-env-new-svc-binding.yaml": triggerBinding}
	if diff := cmp.Diff(resources, wantResources); diff != "" {
		t.Errorf("resources failed: %v", diff)
	}

}
