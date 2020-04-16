package manifest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/spf13/afero"
)

func TestBuildEnvironmentsCreatesEnvironments(t *testing.T) {
	var appFs = afero.NewMemMapFs()
	m := &config.Manifest{
		Environments: []*config.Environment{
			&config.Environment{Name: "test-dev"},
			&config.Environment{Name: "test-staging"},
		},
	}

	files, err := buildEnvironments(appFs, m)
	if err != nil {
		t.Fatal(err)
	}

	want := res.Resources{
		"environments/test-dev/base/test-dev-environment.yaml":         CreateNamespace("test-dev"),
		"environments/test-dev/base/kustomization.yaml":                &res.Kustomization{Resources: []string{"test-dev-environment.yaml"}},
		"environments/test-staging/base/test-staging-environment.yaml": CreateNamespace("test-staging"),
		"environments/test-staging/base/kustomization.yaml":            &res.Kustomization{Resources: []string{"test-staging-environment.yaml"}},
		"environments/test-staging/overlays/kustomization.yaml":        &res.Kustomization{Bases: []string{"../base"}},
		"environments/test-dev/overlays/kustomization.yaml":            &res.Kustomization{Bases: []string{"../base"}},
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Fatalf("files didn't match: %s\n", diff)
	}
}

func TestBuildEnvironmentsDoesNotOutputCIorArgo(t *testing.T) {
	var appFs = afero.NewMemMapFs()
	m := &config.Manifest{
		Environments: []*config.Environment{
			&config.Environment{Name: "test-ci", IsCICD: true},
			&config.Environment{Name: "test-argo", IsArgoCD: true},
		},
	}

	files, err := buildEnvironments(appFs, m)
	if err != nil {
		t.Fatal(err)
	}

	want := res.Resources{}
	if diff := cmp.Diff(want, files); diff != "" {
		t.Fatalf("files didn't match: %s\n", diff)
	}
}

func TestBuildEnvironmentsAddsKustomizedFiles(t *testing.T) {
	var appFs = afero.NewMemMapFs()
	appFs.MkdirAll("environments/test-dev/base", 0755)
	afero.WriteFile(appFs, "environments/test-dev/base/volume.yaml", []byte(`this is a file`), 0644)
	afero.WriteFile(appFs, "environments/test-dev/base/test-dev-environment.yaml", []byte(`this is a file`), 0644)
	afero.WriteFile(appFs, "environments/test-dev/base/routes/01-route.yaml", []byte(`this is a file`), 0644)

	m := &config.Manifest{
		Environments: []*config.Environment{
			&config.Environment{Name: "test-dev"},
		},
	}

	files, err := buildEnvironments(appFs, m)
	if err != nil {
		t.Fatal(err)
	}

	want := res.Resources{
		"environments/test-dev/base/test-dev-environment.yaml": CreateNamespace("test-dev"),
		"environments/test-dev/base/kustomization.yaml":        &res.Kustomization{Resources: []string{"routes/01-route.yaml", "test-dev-environment.yaml", "volume.yaml"}},
		"environments/test-dev/overlays/kustomization.yaml":    &res.Kustomization{Bases: []string{"../base"}},
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Fatalf("files didn't match: %s\n", diff)
	}
}
