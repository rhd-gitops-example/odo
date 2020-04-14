package manifest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/config"
)

func TestBuildCreatesNamespaces(t *testing.T) {
	m := &config.Manifest{
		Environments: []*config.Environment{
			{Name: "dev"},
			{Name: "staging"},
			{Name: "production"},
		},
	}

	files := buildEnvironments(m)
	want := map[string]interface{}{
		"environments/dev/base/config/01_namespace.yaml":        namespace(&config.Environment{Name: "dev"}),
		"environments/staging/base/config/01_namespace.yaml":    namespace(&config.Environment{Name: "staging"}),
		"environments/production/base/config/01_namespace.yaml": namespace(&config.Environment{Name: "production"}),
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Fatalf("didn't get correct files: %s\n", diff)
	}

}

func mustParse(t *testing.T, filename string) *config.Manifest {
	t.Helper()
	m, err := config.ParseFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
