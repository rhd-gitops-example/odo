package pipelines

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidRepository(t *testing.T) {
	tmpDir := os.TempDir()
	repoPath := filepath.Join(tmpDir, "gitops")
	defer os.RemoveAll(repoPath)
	err := createRepositoryLayout(repoPath)
	if err != nil {
		t.Fatalf("create() failed:\n%v", err)
	}
	err = validateRepositoryLayout(repoPath)
	if err != nil {
		t.Fatalf("validate() failed:\n%v", err)
	}
}

func TestInvalidRepository(t *testing.T) {
	tmpDir := os.TempDir()
	repoPath := filepath.Join(tmpDir, "gitops")
	defer os.RemoveAll(repoPath)
	err := createRepositoryLayout(repoPath)
	if err != nil {
		t.Fatalf("create() failed:\n%v", err)
	}
	os.RemoveAll(filepath.Join(repoPath, "env"))
	err = validateRepositoryLayout(repoPath)
	if diff := cmp.Diff(err.Error(), "stat /tmp/gitops/env: no such file or directory"); diff != "" {
		t.Fatalf("validate() failed:\n%v", err)
	}
}
