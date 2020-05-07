package pipelines

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/openshift/odo/pkg/pipelines/ioutils"

	"github.com/openshift/odo/pkg/pipelines/secrets"

	"github.com/spf13/afero"
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

	if err := afero.WriteFile(fakeFs, manifestFile, []byte("environments:\n - name: test\ngitops_url: http://github.com/org/test"), 0644); err != nil {
		t.Fatalf("failed to write manifest to fake filesystem: %s", err)
	}

	if err := AddService("", "test", "app", "svc-1", "123", manifestFile, fakeFs); err != nil {
		t.Fatalf("AddService() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/test/apps/app/base/kustomization.yaml",
		"environments/test/apps/app/kustomization.yaml",
		"environments/test/apps/app/overlays/kustomization.yaml",
		"environments/test/services/svc-1/base/kustomization.yaml",
		"environments/test/services/svc-1/kustomization.yaml",
		"environments/test/services/svc-1/overlays/kustomization.yaml",
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

	kustomization := filepath.Join(gitopsPath, "environments/cicd/base/pipelines/kustomization.yaml")
	afero.WriteFile(fakeFs, kustomization, []byte("resources:\n"), 0644)
	afero.WriteFile(fakeFs, manifestFile, []byte("gitops_url: http://github.com/org/test \nenvironments:\n - name: test\n - cicd: true\n   name: cicd\n"), 0644)
	if err := AddService("", "test", "app1", "svc-1", "123", manifestFile, fakeFs); err != nil {
		t.Fatalf("AddService() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/test/apps/app1/base/kustomization.yaml",
		"environments/test/apps/app1/kustomization.yaml",
		"environments/test/apps/app1/overlays/kustomization.yaml",
		"environments/test/services/svc-1/base/kustomization.yaml",
		"environments/test/services/svc-1/kustomization.yaml",
		"environments/test/services/svc-1/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}
}

func TestEvironmentExists(t *testing.T) {

	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	manifestFile := filepath.Join(gitopsPath, pipelinesFile)

	afero.WriteFile(fakeFs, manifestFile, []byte("gitops_url: http://github.com/org/test \nenvironments:\n - cicd: true\n   name: cicd\n"), 0644)

	err := AddService("", "test", "app1", "svc-1", "123", manifestFile, fakeFs)

	want := `environment test does not exist`
	if err.Error() != want {
		t.Errorf("got %s, want %s", err, want)
	}

}
