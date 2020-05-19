package scm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewRepository(t *testing.T) {
	githubURL := "http://github.com/org/test"
	got, err := NewRepository(githubURL)
	assertNoError(t, err)
	want, err := NewGithubRepository(githubURL)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("NewRepository() failed:\n%s", diff)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
