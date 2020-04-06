package pipelines

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEnv(t *testing.T) {

	gitopsPath, cleanUp := fakeGitopsDir(t)
	defer cleanUp()

	envParameters := EnvParameters{
		EnvName:    "dev",
		GitOpsRepo: "org/gitops",
		Output:     filepath.Dir(gitopsPath),
	}
	if err := Env(&envParameters); err != nil {
		t.Fatalf("Env() failed :%s", err.Error())
	}

	wantedPaths := []string{
		"envs/dev/base/kustomization.yaml",
		"envs/dev/base/namespace.yaml",
		"envs/dev/base/rolebinding.yaml",
		"envs/dev/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s", path), func(t *testing.T) {
			exists, err := isExisting(filepath.Join(gitopsPath, path))
			assertNoError(t, err)
			if !exists {
				t.Fatalf("path %s doesn't exist", path)
			}
		})
	}
}

func fakeGitopsDir(t *testing.T) (string, func()) {
	tmpDir, cleanUp := makeTempDir(t)
	gitopsDir := filepath.Join(tmpDir, "gitops")
	err := os.Mkdir(gitopsDir, 0755)
	if err != nil {
		t.Fatalf("failed to create gitops directory")
	}
	return gitopsDir, cleanUp
}
