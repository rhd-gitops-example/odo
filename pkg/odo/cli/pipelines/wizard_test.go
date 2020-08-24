package pipelines

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

type mockSpinner struct {
	writer io.Writer
	start  bool
	end    bool
}

func (m *mockSpinner) Start(status string, debug bool) {
	m.start = true
	fmt.Fprintf(m.writer, "\n%s", status)
}

func (m *mockSpinner) End(status bool) {
	m.end = status
}

func (m *mockSpinner) WarningStatus(status string) {
	fmt.Fprintf(m.writer, "[%s]", status)
}

func TestCheckSpinner(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		endStatus bool
		wantMsg   string
	}{
		{
			"No error",
			nil,
			true,
			"\nChecking if abcd is installed",
		},
		{
			"Resource not found error",
			errors.NewNotFound(schema.GroupResource{}, "abcd"),
			false,
			"\nChecking if abcd is installed[Please install abcd]",
		},
		{
			"Random cluster error",
			fmt.Errorf("Sample cluster error"),
			false,
			"\nChecking if abcd is installed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buff := &bytes.Buffer{}

			fakeSpinner := &mockSpinner{writer: buff}
			fakeSpinner.Start("Checking if abcd is installed", false)
			setSpinnerStatus(fakeSpinner, "Please install abcd", test.err)

			if fakeSpinner.end != test.endStatus {
				t.Errorf("Spinner status mismatch: got %v, want %v", fakeSpinner.end, test.endStatus)
			}
			assertMessage(t, buff.String(), test.wantMsg)
		})
	}
}

func TestDependenciesWithNothingInstalled(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()

	wantMsg := `
Checking if Sealed Secrets is installed with the default configuration[Please install Sealed Secrets from https://github.com/bitnami-labs/sealed-secrets/releases]
Checking if ArgoCD Operator is installed with the default configuration[Please install ArgoCD operator from OperatorHub]
Checking if OpenShift Pipelines Operator is installed with the default configuration[Please install OpenShift Pipelines operator from OperatorHub]`

	buff := &bytes.Buffer{}
	fakeSpinner := &mockSpinner{writer: buff}
	err := checkBootstrapDependencies(&WizardParameters{&pipelines.BootstrapOptions{}, &genericclioptions.Context{}}, fakeClient, fakeSpinner)
	wantErr := "Failed to satisfy the required dependencies"

	assertError(t, err, wantErr)
	assertMessage(t, buff.String(), wantMsg)
}

func TestDependenciesWithAllInstalled(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(sealedSecretService(), argoCDOperator(), pipelinesOperator())

	wantMsg := `
Checking if Sealed Secrets is installed with the default configuration
Checking if ArgoCD Operator is installed with the default configuration
Checking if OpenShift Pipelines Operator is installed with the default configuration`

	buff := &bytes.Buffer{}
	fakeSpinner := &mockSpinner{writer: buff}
	wizardParams := &WizardParameters{&pipelines.BootstrapOptions{}, &genericclioptions.Context{}}
	err := checkBootstrapDependencies(wizardParams, fakeClient, fakeSpinner)

	assertError(t, err, "")
	if wizardParams.SealedSecretsService.Name != "sealed-secrets-controller" && wizardParams.SealedSecretsService.Namespace != "kube-system" {
		t.Fatalf("Expected sealed secrets to be set")
	}
	assertMessage(t, buff.String(), wantMsg)
}

func TestDependenciesWithNoArgoCD(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(sealedSecretService(), pipelinesOperator())

	wantMsg := `
Checking if Sealed Secrets is installed with the default configuration
Checking if ArgoCD Operator is installed with the default configuration[Please install ArgoCD operator from OperatorHub]
Checking if OpenShift Pipelines Operator is installed with the default configuration`

	buff := &bytes.Buffer{}
	fakeSpinner := &mockSpinner{writer: buff}
	wizardParams := &WizardParameters{&pipelines.BootstrapOptions{}, &genericclioptions.Context{}}
	err := checkBootstrapDependencies(wizardParams, fakeClient, fakeSpinner)
	wantErr := "Failed to satisfy the required dependencies"

	assertError(t, err, wantErr)
	assertMessage(t, buff.String(), wantMsg)
}

func TestDependenciesWithNoPipelines(t *testing.T) {
	fakeClient := fake.NewSimpleClientset(sealedSecretService(), argoCDOperator())

	wantMsg := `
Checking if Sealed Secrets is installed with the default configuration
Checking if ArgoCD Operator is installed with the default configuration
Checking if OpenShift Pipelines Operator is installed with the default configuration[Please install OpenShift Pipelines operator from OperatorHub]`

	buff := &bytes.Buffer{}
	fakeSpinner := &mockSpinner{writer: buff}
	wizardParams := &WizardParameters{&pipelines.BootstrapOptions{}, &genericclioptions.Context{}}
	err := checkBootstrapDependencies(wizardParams, fakeClient, fakeSpinner)
	wantErr := "Failed to satisfy the required dependencies"

	assertError(t, err, wantErr)
	assertMessage(t, buff.String(), wantMsg)
}

func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		if msg != "" {
			t.Fatalf("Error mismatch: got %v, want %v", err, msg)
		}
		return
	}
	if err.Error() != msg {
		t.Fatalf("Error mismatch: got %s, want %s", err.Error(), msg)
	}
}

func assertMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Message mismatch: got %s, want %s", got, want)
	}
}

func sealedSecretService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sealed-secrets-controller",
			Namespace: "kube-system",
		},
	}
}

func argoCDOperator() *appv1.DeploymentList {
	return &appv1.DeploymentList{
		Items: []appv1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "argocd-operator",
					Namespace: "argocd",
				},
			}, {
				ObjectMeta: metav1.ObjectMeta{
					Name:      "argocd-server",
					Namespace: "argocd",
				},
			},
		},
	}
}

func pipelinesOperator() *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openshift-pipelines-operator",
			Namespace: "openshift-operators",
		},
	}
}
