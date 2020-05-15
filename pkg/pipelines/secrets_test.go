package pipelines

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateOpaqueSecret(t *testing.T) {
	data := []byte(`abcdefghijklmnop`)
	secret, err := CreateOpaqueSecret("github-auth", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	want := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "github-auth",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"token": data,
		},
	}

	if diff := cmp.Diff(want, secret); diff != "" {
		t.Fatalf("CreateOpaqueSecret() failed got\n%s", diff)
	}
}

func TestCreateOpaqueSecretWithErrorReading(t *testing.T) {
	testErr := errors.New("test failure")
	_, err := CreateOpaqueSecret("github-auth", errorReader{testErr})
	if !matchError(t, "failed to read .* test failure", err) {
		t.Fatalf("got an unexpected error: %#v", err)
	}
}

func TestCreateDockerConfigSecret(t *testing.T) {
	data := []byte(`abcdefghijklmnop`)
	secret, err := CreateDockerConfigSecret("regcred", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	want := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "regcred",
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": data,
		},
	}

	if diff := cmp.Diff(want, secret); diff != "" {
		t.Fatalf("CreateDockerConfigSecret() failed got\n%s", diff)
	}
}

type errorReader struct {
	err error
}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, e.err
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	if s == "" && e == nil {
		return true
	}
	if s != "" && e == nil {
		return false
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}
