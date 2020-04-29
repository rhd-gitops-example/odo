package webhook

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetRouteHost(t *testing.T) {
	resources, err := NewResources()
	if err != nil {
		t.Fatal(err)
	}

	route, err := resources.getListenerAddress("tst-cicd", "gitops-webhook-event-listener-route")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(route, "gitops-webhook-event-listener-route-tst-cicd.apps.gitops1.devcluster.openshift.com"); diff != "" {
		t.Errorf("driver errMsg mismatch got\n%s", diff)
	}
}

func TestGetSecret(t *testing.T) {
	resources, err := NewResources()
	if err != nil {
		t.Fatal(err)
	}

	secret, err := resources.getWebhookSecret("tst-cicd", "gitops-webhook-secret", "webhook-secret-key")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(secret, "testing"); diff != "" {
		t.Errorf("driver errMsg mismatch got\n%s", diff)
	}

}