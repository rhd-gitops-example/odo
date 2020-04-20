package manifest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
)

var testCICDEnv = &config.Environment{Name: "tst-cicd", IsCICD: true}

func TestCreateManifest(t *testing.T) {
	want := &config.Manifest{
		Environments: []*config.Environment{
			testCICDEnv,
		},
	}
	got := createManifest(testCICDEnv)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("manifest didn't match: %s\n", diff)
	}
}

func TestInitialFiles(t *testing.T) {
	prefix := "tst-"
	gitOpsRepo := "test-repo"
	gitOpsWebhook := "123"
	imageRepo := "image/repo"

	got, err := createInitialFiles(prefix, gitOpsRepo, gitOpsWebhook, "", imageRepo)

	if err != nil {
		t.Fatal(err)
	}
	want := res.Resources{
		"manifest.yaml": createManifest(testCICDEnv),
	}

	cicdResources, err := CreateResources(prefix, gitOpsRepo, gitOpsWebhook, "", imageRepo)
	if err != nil {
		t.Fatalf("CreatePipelineResources() failed due to :%s\n", err)
	}
	files := getResourceFiles(cicdResources)

	want = res.Merge(addPrefixToResources("environments/tst-cicd/base/pipelines", cicdResources), want)

	want = res.Merge(addPrefixToResources("environments/tst-cicd", getCICDKustomization(files)), want)

	if diff := cmp.Diff(want, got, cmpopts.IgnoreMapEntries(ignoreSecrets)); diff != "" {
		t.Fatalf("outputs didn't match: %s\n", diff)
	}
}

func ignoreSecrets(k string, v interface{}) bool {
	if k == "environments/tst-cicd/base/pipelines/03-secrets/gitops-webhook-secret.yaml" {
		return true
	}
	return false
}

func TestGetCICDKustomization(t *testing.T) {
	want := res.Resources{
		"base/kustomization.yaml": map[string]interface{}{
			"bases": []string{"./pipelines"},
		},
		"overlays/kustomization.yaml": map[string]interface{}{
			"bases": []string{"../base"},
		},
		"base/pipelines/kustomization.yaml": map[string]interface{}{
			"resources": []string{"resource1", "resource2"},
		},
	}
	got := getCICDKustomization([]string{"resource1", "resource2"})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("getCICDKustomization was not correct: %s\n", diff)
	}

}

func TestAddPrefixToResources(t *testing.T) {
	files := map[string]interface{}{
		"base/kustomization.yaml": map[string]interface{}{
			"resources": []string{},
		},
		"overlays/kustomization.yaml": map[string]interface{}{
			"bases": []string{"../base"},
		},
	}

	want := map[string]interface{}{
		"test-prefix/base/kustomization.yaml": map[string]interface{}{
			"resources": []string{},
		},
		"test-prefix/overlays/kustomization.yaml": map[string]interface{}{
			"bases": []string{"../base"},
		},
	}
	if diff := cmp.Diff(want, addPrefixToResources("test-prefix", files)); diff != "" {
		t.Fatalf("addPrefixToResources failed, diff %s\n", diff)
	}
}

func TestMerge(t *testing.T) {
	map1 := map[string]interface{}{
		"test-1": "value-1",
	}
	map2 := map[string]interface{}{
		"test-1": "value-a",
		"test-2": "value-2",
	}
	map3 := map[string]interface{}{
		"test-1": "value-a",
		"test-2": "value-2",
	}

	want := res.Resources{
		"test-1": "value-1",
		"test-2": "value-2",
	}
	if diff := cmp.Diff(want, res.Merge(map1, map2)); diff != "" {
		t.Fatalf("merge failed: %s\n", diff)
	}
	if diff := cmp.Diff(map2, map3); diff != "" {
		t.Fatalf("original map changed %s\n", diff)
	}

}
