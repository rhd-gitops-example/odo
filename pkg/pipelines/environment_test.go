package pipelines

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/tests/helper"
)

func TestAddEnv(t *testing.T) {
	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")
	pipelinesFile := filepath.Join(gitopsPath, pipelinesFile)
	envParameters := EnvParameters{
		PipelinesFilePath: pipelinesFile,
		EnvName:           "dev",
	}
	_ = afero.WriteFile(fakeFs, pipelinesFile, []byte("environments:"), 0644)

	if err := AddEnv(&envParameters, fakeFs); err != nil {
		t.Fatalf("AddEnv() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/dev/env/base/kustomization.yaml",
		"environments/dev/env/base/dev-environment.yaml",
		"environments/dev/env/overlays/kustomization.yaml",
	}
	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}

	got := mustReadFileAsMap(t, fakeFs, pipelinesFile)
	want := map[string]interface{}{
		"environments": []interface{}{
			map[string]interface{}{
				"name": "dev",
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("written environments failed:\n%s", diff)
	}

}

func TestAddEnvWithClusterProvided(t *testing.T) {
	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")
	pipelinesFile := filepath.Join(gitopsPath, pipelinesFile)
	envParameters := EnvParameters{
		PipelinesFilePath: pipelinesFile,
		EnvName:           "dev",
		Cluster:           "testing.cluster",
	}
	_ = afero.WriteFile(fakeFs, pipelinesFile, []byte("environments:"), 0644)

	if err := AddEnv(&envParameters, fakeFs); err != nil {
		t.Fatalf("AddEnv() failed :%s", err)
	}

	wantedPaths := []string{
		"environments/dev/env/base/kustomization.yaml",
		"environments/dev/env/base/dev-environment.yaml",
		"environments/dev/env/overlays/kustomization.yaml",
	}
	for _, path := range wantedPaths {
		t.Run(fmt.Sprintf("checking path %s already exists", path), func(rt *testing.T) {
			assertFileExists(rt, fakeFs, filepath.Join(gitopsPath, path))
		})
	}

	got := mustReadFileAsMap(t, fakeFs, pipelinesFile)
	want := map[string]interface{}{
		"environments": []interface{}{
			map[string]interface{}{
				"cluster": "testing.cluster",
				"name":    "dev",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("written environments failed:\n%s", diff)
	}
}

func TestAddEnvWithExistingName(t *testing.T) {
	fakeFs := ioutils.NewMapFilesystem()
	gitopsPath := afero.GetTempDir(fakeFs, "test")

	pipelinesFile := filepath.Join(gitopsPath, pipelinesFile)
	envParameters := EnvParameters{
		PipelinesFilePath: pipelinesFile,
		EnvName:           "dev",
	}
	_ = afero.WriteFile(fakeFs, pipelinesFile, []byte("environments:\n - name: dev\n"), 0644)

	if err := AddEnv(&envParameters, fakeFs); err == nil {
		t.Fatal("AddEnv() did not fail with duplicate environment")
	}
}

func TestNewEnvironment(t *testing.T) {
	tests := []struct {
		m      *config.Manifest
		name   string
		errMsg string
		want   *config.Environment
	}{
		{
			m: &config.Manifest{
				GitOpsURL: "https://github.com/foo/bar",
				Config: &config.Config{
					Pipelines: &config.PipelinesConfig{
						Name: "my-cicd",
					},
				},
				Environments: []*config.Environment{
					{
						Name: "myenv1",
					},
				},
			},
			name: "test-env",
			want: &config.Environment{
				Name: "test-env",
				Pipelines: &config.Pipelines{
					Integration: &config.TemplateBinding{
						Template: appCITemplateName,
						Bindings: []string{"github-push-binding"},
					},
				},
			},
		},
		{
			m: &config.Manifest{
				GitOpsURL: "https://gitlab.com/foo/bar",
				Config: &config.Config{
					Pipelines: &config.PipelinesConfig{
						Name: "my-cicd",
					},
				},
				Environments: []*config.Environment{
					{
						Name: "my-cicd",
					},
				},
			},
			name: "test-env",
			want: &config.Environment{
				Name: "test-env",
				Pipelines: &config.Pipelines{
					Integration: &config.TemplateBinding{
						Template: appCITemplateName,
						Bindings: []string{"gitlab-push-binding"},
					},
				},
			},
		},
		{
			m: &config.Manifest{
				// no GitOpsURL -> no Pipelines
				Config: &config.Config{
					Pipelines: &config.PipelinesConfig{
						Name: "my-cicd",
					},
				},
				Environments: []*config.Environment{
					{
						Name: "my-env2",
					},
				},
			},
			name: "test-env",
			want: &config.Environment{
				Name: "test-env",
			},
		},
		{
			m: &config.Manifest{
				GitOpsURL: "https://gitlab.com/foo/bar",
				Environments: []*config.Environment{
					{
						// no CICD -> no Pipelines
						Name: "my-env4",
					},
				},
			},
			name: "test-env",
			want: &config.Environment{
				Name: "test-env",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test_%d", i), func(rt *testing.T) {
			got, err := newEnvironment(tt.m, tt.name)

			if !helper.ErrorMatch(rt, tt.errMsg, err) {
				rt.Errorf("err mismatch want: %s got: %s: \n", tt.errMsg, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				rt.Errorf("env mismatch: \n%s", diff)
			}
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertFileExists(t *testing.T, testFs afero.Fs, path string) {
	t.Helper()
	exists, err := afero.Exists(testFs, path)
	assertNoError(t, err)
	if !exists {
		t.Fatalf("unable to find file %q", path)
	}
	isDir, err := afero.DirExists(testFs, path)
	assertNoError(t, err)
	if isDir {
		t.Fatalf("%q is a directory", path)
	}
}

func mustReadFileAsMap(t *testing.T, fs afero.Fs, filename string) map[string]interface{} {
	t.Helper()
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{}
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		t.Fatal(err)
	}
	return m

}
