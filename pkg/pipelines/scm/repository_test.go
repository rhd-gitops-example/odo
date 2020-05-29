package scm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewRepository(t *testing.T) {
	githubURL := "http://github.com/org/test"
	got, err := NewRepository(githubURL)
	assertNoError(t, err)
	want, err := NewGitHubRepository(githubURL)
	assertNoError(t, err)
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(GitHubRepository{})); diff != "" {
		t.Fatalf("NewRepository() failed:\n%s", diff)
	}
}

func TestNewRepositoryForInvalidRepoType(t *testing.T) {
	githubURL := "http://test.com/org/test"
	repoType := "test"
	_, gotErr := NewRepository(githubURL)
	if gotErr == nil {
		t.Fatalf("NewRepository() returned an invalid repository of type: %s", repoType)
	}
	wantErr := invalidRepoTypeError(githubURL)
	if diff := cmp.Diff(wantErr.Error(), gotErr.Error()); diff != "" {
		t.Fatalf("Errors don't match: got %v want %v", gotErr, wantErr)
	}
}

func TestIsRepositoryBinding(t *testing.T) {
	tests := []struct {
		desc string
		name string
		want bool
	}{
		{
			"valid repository binding",
			"github-pr-binding",
			true,
		},
		{
			"valid repository binding",
			"gitlab-pr-binding",
			true,
		},
		{
			"invalid repository binding",
			"test-binding",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(rt *testing.T) {
			got := IsRepositoryBinding(test.name)
			if test.want != got {
				rt.Fatalf("IsRepositoryBinding failed: got %v want %v", got, test.want)
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
