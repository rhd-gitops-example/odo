package pipelines

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {

	gitopsPath, cleanUp := fakeGitopsDir(t)
	outputPath := filepath.Dir(gitopsPath)
	defer cleanUp()

	envParameters := EnvParameters{
		EnvName:    "dev",
		GitOpsRepo: "org/gitops",
		Output:     outputPath,
	}
	if err := Env(&envParameters); err != nil {
		t.Fatalf("Env() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/dev/base/kustomization.yaml",
		"environments/dev/base/namespace.yaml",
		"environments/dev/base/rolebinding.yaml",
		"environments/dev/overlays/kustomization.yaml",
	}

	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(t *testing.T) {
			assert.FileExists(t, filepath.Join(gitopsPath, path))
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
