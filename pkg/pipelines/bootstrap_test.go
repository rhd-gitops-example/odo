package pipelines

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/deployment"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/meta"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	"k8s.io/apimachinery/pkg/types"
)

const (
	testSvcRepo    = "https://github.com/my-org/http-api.git"
	testGitOpsRepo = "https://github.com/my-org/gitops.git"
)

func TestBootstrapManifest(t *testing.T) {
	defer func(f secrets.PublicKeyFunc) {
		secrets.DefaultPublicKeyFunc = f
	}(secrets.DefaultPublicKeyFunc)

	secrets.DefaultPublicKeyFunc = func(service types.NamespacedName) (*rsa.PublicKey, error) {
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Fatalf("failed to generate a private RSA key: %s", err)
		}
		return &key.PublicKey, nil
	}

	params := &BootstrapOptions{
		InitOptions: &InitOptions{
			Prefix:              "tst-",
			GitOpsRepoURL:       testGitOpsRepo,
			ImageRepo:           "image/repo",
			GitOpsWebhookSecret: "123",
		},
		ServiceRepoURL:       testSvcRepo,
		ServiceWebhookSecret: "456",
	}

	r, err := bootstrapResources(params, ioutils.NewMapFilesystem())
	if err != nil {
		t.Fatal(err)
	}
	hookSecret, err := secrets.CreateSealedSecret(meta.NamespacedName("tst-cicd", "webhook-secret-tst-dev-http-api"), meta.NamespacedName("test-ns", "service"), "456", eventlisteners.WebhookSecretKey)
	if err != nil {
		t.Fatal(err)
	}
	want := res.Resources{
		"config/tst-cicd/base/03-secrets/webhook-secret-tst-dev-http-api.yaml":   hookSecret,
		"environments/tst-dev/services/http-api/base/config/100-deployment.yaml": deployment.Create("app-taxi", "tst-dev", "http-api", bootstrapImage, deployment.ContainerPort(8080)),
		"environments/tst-dev/services/http-api/base/config/200-service.yaml":    createBootstrapService("app-taxi", "tst-dev", "http-api"),
		"environments/tst-dev/services/http-api/base/config/kustomization.yaml":  &res.Kustomization{Resources: []string{"100-deployment.yaml", "200-service.yaml"}},
		pipelinesFile: &config.Manifest{
			GitOpsURL: "https://github.com/my-org/gitops.git",
			Environments: []*config.Environment{
				{
					Pipelines: &config.Pipelines{
						Integration: &config.TemplateBinding{
							Template: "app-ci-template",
							Bindings: []string{"github-push-binding"},
						},
					},
					Name: "tst-dev",
					Services: []*config.Service{
						{
							Name:      "http-api",
							SourceURL: testSvcRepo,
							Webhook: &config.Webhook{
								Secret: &config.Secret{
									Name:      "webhook-secret-tst-dev-http-api",
									Namespace: "tst-cicd",
								},
							},
							Pipelines: &config.Pipelines{
								Integration: &config.TemplateBinding{Bindings: []string{"tst-dev-http-api-binding", "github-push-binding"}},
							},
						},
					},

					Apps: []*config.Application{
						{
							Name:        "app-http-api",
							ServiceRefs: []string{"http-api"},
						},
					},
				},
				{Name: "tst-stage"},
			},
			Config: &config.Config{
				Pipelines: &config.PipelinesConfig{Name: "tst-cicd"},
				ArgoCD:    &config.ArgoCDConfig{Namespace: "argocd"},
			},
		},
	}

	if diff := cmp.Diff(want, r, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		_, ok := want[k]
		return !ok
	})); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}

	wantResources := []string{"01-namespaces/cicd-environment.yaml",
		"01-namespaces/image.yaml",
		"02-rolebindings/internal-registry-image-binding.yaml",
		"02-rolebindings/pipeline-service-role.yaml",
		"02-rolebindings/pipeline-service-rolebinding.yaml",
		"03-secrets/gitops-webhook-secret.yaml",
		"03-secrets/webhook-secret-tst-dev-http-api.yaml",
		"04-tasks/deploy-from-source-task.yaml",
		"04-tasks/deploy-using-kubectl-task.yaml",
		"05-pipelines/app-ci-pipeline.yaml",
		"05-pipelines/ci-dryrun-from-push-pipeline.yaml",
		"06-bindings/github-push-binding.yaml",
		"06-bindings/tst-dev-http-api-binding.yaml",
		"07-templates/app-ci-build-from-push-template.yaml",
		"07-templates/ci-dryrun-from-push-template.yaml",
		"08-eventlisteners/cicd-event-listener.yaml",
		"09-routes/gitops-webhook-event-listener.yaml",
	}
	k := r["config/tst-cicd/base/kustomization.yaml"].(res.Kustomization)
	if diff := cmp.Diff(wantResources, k.Resources); diff != "" {
		t.Fatalf("did not add the secret to the base kustomization: %s\n", diff)
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
		Name:        "app-http-api",
		ServiceRefs: []string{"http-api"},
	}

	got, err := applicationFromRepo(testSvcRepo, "http-api")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}
}
