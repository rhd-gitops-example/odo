package pipelines

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func writeResources(prefix string, files map[string][]interface{}) ([]string, error) {
	filenames := make([]string, len(files))
	for filename, items := range files {
		err := marshalItemsToFile(filepath.Join(prefix, filename), items)
		if err != nil {
			return err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func marshalItemsToFile(filename string, items []interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return marshalOutputs(f, items)
}

func TestWriteResources(t *testing.T) {
	tmpDir, cleanUp := makeTempDir(t)
	defer cleanUp()
	resources := map[string][]interface{}{
		"01_roles/serviceaccount.yaml": []interface{}{fakeYamlDoc(1)},
		"02_tasks/buildah_task.yaml":   []interface{}{fakeYamlDoc(1), fakeYamlDoc(2)},
	}

	err := writeResources(tmpDir, resources)
	if err != nil {
		t.Fatalf("failed to writeResources: %v", err)
	}
	assertFileContents(t, filepath.Join(tmpDir, "01_roles/serviceaccount.yaml"), []byte("key1: value1\n---\n"))
	assertFileContents(t, filepath.Join(tmpDir, "02_tasks/buildah_task.yaml"), []byte("key1: value1\n---\nkey2: value2\n---\n"))
}

func assertFileContents(t *testing.T, filename string, want []byte) {
	t.Helper()
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %s", filename)
	}

	if diff := cmp.Diff(body, want); diff != "" {
		t.Fatalf("file %s diff = \n%s\n", filename, diff)
	}
}

func TestWriteResourcesToFile(t *testing.T) {
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
	assertRepositoryLayout(t, []string{filepath.Join(gitopsPath, "01-test1.yaml"),
		filepath.Join(gitopsPath, "02-test2.yaml")})
}

func TestCreatePipelineResource(t *testing.T) {
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

func assertRepositoryLayout(t *testing.T, paths []string) {
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
