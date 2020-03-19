package pipelines

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidRepository(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create a temp directory:\n%v", err)
	}
	gitopsPath := filepath.Join(tmpDir, "gitops")
	defer os.RemoveAll(tmpDir)
	folder := getGitopsFolder(gitopsPath, "test")
	if err := folder.create(); err != nil {
		t.Fatalf("create() failed:\n%v", err)
	}
	if err := folder.validate(); err != nil {
		t.Fatalf("validate() failed:\n%v", err)
	}
}

func TestInvalidRepository(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create a temp directory:\n%v", err)
	}
	gitopsPath := filepath.Join(tmpDir, "gitops")
	defer os.RemoveAll(tmpDir)
	folder := getGitopsFolder(gitopsPath, "test")
	if err := folder.create(); err != nil {
		t.Fatalf("create() failed:\n%v", err)
	}
	envs := filepath.Join(folder.path, "envs")
	// alter gitops dir structure
	os.RemoveAll(envs)

	validErrMsg := fmt.Sprintf("stat %s: no such file or directory", envs)
	err = folder.validate()
	if diff := cmp.Diff(err.Error(), validErrMsg); diff != "" {
		t.Fatalf("validate() failed:\n%v", err)
	}
}

func TestAddKusomize(t *testing.T) {
	validDir := folder{
		path:    "test/kustomization.yaml",
		isDir:   false,
		folders: []*folder{},
	}
	dir := addKustomize("test")
	if diff := cmp.Diff(validDir, *dir, cmp.AllowUnexported(validDir, *dir)); diff != "" {
		t.Fatalf("addKustomize() failed:\n%v", diff)
	}
}
