package manifest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/config"
)

func TestBootstrapManifest(t *testing.T) {
	params := &BootstrapParameters{Prefix: "tst-"}
	r, err := bootstrapResources(params)
	if err != nil {
		t.Fatal(err)
	}

	want := resources{
		"manifest.yaml": &config.Manifest{
			Environments: []*config.Environment{
				{Name: "tst-dev"},
				{Name: "tst-stage"},
				{Name: "tst-cicd", IsCICD: true},
			},
		},
	}

	if diff := cmp.Diff(want, r); diff != "" {
		t.Fatalf("bootstrapped resources:\n%s", diff)
	}
}
