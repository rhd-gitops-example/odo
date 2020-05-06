package pipelines

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/namespaces"
	res "github.com/openshift/odo/pkg/pipelines/resources"
	"github.com/openshift/odo/pkg/pipelines/secrets"
	yamlWriter "github.com/openshift/odo/pkg/pipelines/yaml"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

func TestServiceWithoutCICD(t *testing.T) {
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
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	manifestFile := filepath.Join(gitopsPath, pipelinesFile)

	afero.WriteFile(fakeFs, manifestFile, []byte("gitops_url: http://github.com/org/test \nenvironments:\n - name: test\n "), 0644)
	o := &AddOptions{
		AppName:       "app1",
		EnvName:       "test",
		Manifest:      manifestFile,
		WebhookSecret: "123",
		GitRepoURL:    "http://github.com/org/http-api",
	}
	if err := AddService(o, fakeFs); err != nil {
		t.Fatalf("AddService() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/test/apps/app1/base/kustomization.yaml",
		"environments/test/apps/app1/kustomization.yaml",
		"environments/test/apps/app1/overlays/kustomization.yaml",
		"environments/test/services/http-api/base/kustomization.yaml",
		"environments/test/services/http-api/kustomization.yaml",
		"environments/test/services/http-api/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}
}

func TestServiceWithCICD(t *testing.T) {
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
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	manifestFile := filepath.Join(gitopsPath, pipelinesFile)
	o := &AddOptions{
		AppName:       "app1",
		EnvName:       "test",
		Manifest:      manifestFile,
		WebhookSecret: "123",
		GitRepoURL:    "http://github.com/org/http-api",
	}
	kustomization := filepath.Join(gitopsPath, "environments/cicd/base/pipelines/kustomization.yaml")
	afero.WriteFile(fakeFs, kustomization, []byte("resources:\n"), 0644)
	afero.WriteFile(fakeFs, manifestFile, []byte("gitops_url: http://github.com/org/test \nenvironments:\n - name: test\n - cicd: true\n   name: cicd\n"), 0644)
	if err := AddService(o, fakeFs); err != nil {
		t.Fatalf("AddService() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/test/apps/app1/base/kustomization.yaml",
		"environments/test/apps/app1/kustomization.yaml",
		"environments/test/apps/app1/overlays/kustomization.yaml",
		"environments/test/services/http-api/base/kustomization.yaml",
		"environments/test/services/http-api/kustomization.yaml",
		"environments/test/services/http-api/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}
}

func TestUpdateKustomization(t *testing.T) {
	files := res.Resources{}
	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")
	files["namespace-1.yaml"] = namespaces.Create("namespace-1")
	files["namespace-2.yaml"] = namespaces.Create("namespace-2")
	files["namespace-3.yaml"] = namespaces.Create("namespace-3")
	files[Kustomize] = res.Kustomization{Resources: []string{"namespace-1.yaml", "namespace-2.yaml"}}
	_, err := yamlWriter.WriteResources(fakeFs, gitopsPath, files)
	assertNoError(t, err)
	want := res.Kustomization{Resources: []string{"namespace-1.yaml", "namespace-2.yaml", "namespace-3.yaml"}}
	err = updateKustomization(fakeFs, gitopsPath)
	assertNoError(t, err)
	b, err := afero.ReadFile(fakeFs, filepath.Join(gitopsPath, Kustomize))
	assertNoError(t, err)
	got := res.Kustomization{}
	err = yaml.Unmarshal(b, got)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("updateKustomization() failed: %v", diff)
	}
}
