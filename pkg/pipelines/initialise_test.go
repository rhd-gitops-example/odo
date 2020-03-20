package pipelines

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteResourcesToFile(t *testing.T) {
	t.Helper()
	tmpDir, cleanUp := makeTempDir(t)
	defer cleanUp()
	gitopsPath := filepath.Join(tmpDir, "gitops")
	resources := []string{"test1", "test2"}
	outputs := map[string]interface{}{
		"test1": fakeYamlDoc(1),
		"test2": fakeYamlDoc(2),
	}
	err := writeResourcesToFile(resources, gitopsPath, "", outputs)
	if err != nil {
		t.Fatalf("writeResourcesToFile() failed:\n%v", err)
	}
	AssertRepositoryLayout(t, []string{filepath.Join(gitopsPath, "01-test1.yaml"),
		filepath.Join(gitopsPath, "02-test2.yaml")})
}

func TestCreatePipelineResource(t *testing.T) {
	t.Helper()
	namespaces := namespaceNames("")
	outputs := createPipelineResources(namespaces, "gitops", "gitops", "")
	wantedResources := getOrderedResources()
	validResult := true
	for _, resource := range wantedResources {
		_, exists := outputs[resource]
		if diff := cmp.Diff(exists, validResult); diff != "" {
			t.Fatalf("resource %v not found", resource)
		}
	}
}

func AssertRepositoryLayout(t *testing.T, paths []string) {
	t.Helper()
	for _, path := range paths {
		if exists, err := isExisting(path); !exists {
			assertNoError(t, err)
		}
	}
}

func makeTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "test")
	assertNoError(t, err)
	return dir, func() {
		err := os.RemoveAll(dir)
		assertNoError(t, err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileNameCreation(t *testing.T) {
	validFileName := "01-test-foo.yaml"
	name := fileName(1, "test-", "foo")
	if diff := cmp.Diff(validFileName, name); diff != "" {
		t.Fatalf("fileName() failed:\n%v", diff)
	}
}

func fakeYamlDoc(n int) map[string]string {
	return map[string]string{
		fmt.Sprintf("key%d", n): fmt.Sprintf("value%d", n),
	}
}
