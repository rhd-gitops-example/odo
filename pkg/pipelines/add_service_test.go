package pipelines

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
)

func TestServiceWithoutCICD(t *testing.T) {

	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	manifestFile := filepath.Join(gitopsPath, pipelinesFile)

	afero.WriteFile(fakeFs, manifestFile, []byte("environments:\n - name: test\n "), 0644)
	if err := AddService(testSvcRepo, "123", "test", "app1", manifestFile, fakeFs); err != nil {
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

	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	manifestFile := filepath.Join(gitopsPath, pipelinesFile)
	kustomization := filepath.Join(gitopsPath, "environments/cicd/base/pipelines/kustomization.yaml")
	afero.WriteFile(fakeFs, kustomization, []byte(""), 0644)
	afero.WriteFile(fakeFs, manifestFile, []byte("environments:\n - name: test\n - cicd: true\n   name: cicd\n"), 0644)
	if err := AddService(testSvcRepo, "123", "test", "app1", manifestFile, fakeFs); err != nil {
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
