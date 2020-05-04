package pipelines

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
)

func TestService(t *testing.T) {

	fakeFs := ioutils.NewFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")
	manifestFile := "config/testdata/pipelines.yaml"
	svcParameters := ServiceParameters{
		ServiceGitRepo:       testSvcRepo,
		ServiceWebhookSecret: "123",
		AppName:              "app1",
		EnvName:              "test",
		Output:               gitopsPath,
		Manifest:             manifestFile,
	}

	if err := AddService(&svcParameters, fakeFs); err != nil {
		t.Fatalf("AddService() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/test/apps/app1/base/kustomization.yaml",
		"environments/test/apps/app1/kustomization.yaml",
		"environments/test/apps/app1/overlays/kustomization.yaml",
		"environments/test/services/myservice/base/kustomization.yaml",
		"environments/test/services/myservice/kustomization.yaml",
		"environments/test/services/myservice/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}
}
