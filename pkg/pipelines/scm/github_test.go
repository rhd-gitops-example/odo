package scm

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/pipelines/triggers"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreatePRBindingForGithub(t *testing.T) {
	repo, err := NewRepository("http://github.com/org/test")
	assertNoError(t, err)
	want := triggersv1.TriggerBinding{
		TypeMeta: triggers.TriggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "github-pr-binding",
			Namespace: "testns",
		},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "gitref",
					Value: "$(body.pull_request.head.ref)",
				},
				{
					Name:  "gitsha",
					Value: "$(body.pull_request.head.sha)",
				},
				{
					Name:  "gitrepositoryurl",
					Value: "$(body.repository.clone_url)",
				},
				{
					Name:  "fullname",
					Value: "$(body.repository.full_name)",
				},
			},
		},
	}
	got, name := repo.CreatePRBinding("testns")
	if name != githubPRBindingName {
		t.Fatalf("CreatePushBinding() returned a wrong binding: want %v got %v", githubPRBindingName, name)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("createPRBinding() failed:\n%s", diff)
	}
}

func TestCreatePushBindingForGithub(t *testing.T) {
	repo, err := NewRepository("http://github.com/org/test")
	assertNoError(t, err)
	want := triggersv1.TriggerBinding{
		TypeMeta: triggers.TriggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "github-push-binding",
			Namespace: "testns",
		},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "gitref",
					Value: "$(body.ref)",
				},
				{
					Name:  "gitsha",
					Value: "$(body.head_commit.id)",
				},
				{
					Name:  "gitrepositoryurl",
					Value: "$(body.repository.clone_url)",
				},
			},
		},
	}
	got, name := repo.CreatePushBinding("testns")
	if name != githubPushBindingName {
		t.Fatalf("CreatePushBinding() returned a wrong binding: want %v got %v", githubPushBindingName, name)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreatePushBinding() failed:\n%s", diff)
	}
}

func TestCreateCITriggerForGithub(t *testing.T) {
	repo, err := NewRepository("http://github.com/org/test")
	assertNoError(t, err)
	want := triggersv1.EventListenerTrigger{
		Name: "test",
		Bindings: []*triggersv1.EventListenerBinding{
			{Name: "test-binding"},
		},
		Template: triggersv1.EventListenerTemplate{Name: "test-template"},
		Interceptors: []*triggersv1.EventInterceptor{
			{
				CEL: &triggersv1.CELInterceptor{
					Filter: fmt.Sprintf(githubCIDryRunFilters, "org/test"),
				},
			},
			{
				GitHub: &triggersv1.GitHubInterceptor{
					SecretRef: &triggersv1.SecretRef{SecretKey: "webhook-secret-key", SecretName: "secret", Namespace: "ns"},
				},
			},
		},
	}
	got := repo.CreateCITrigger("test", "secret", "ns", "test-template", []string{"test-binding"})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreateCITrigger() failed:\n%s", diff)
	}
}

func TestCreateCDTriggersForGithub(t *testing.T) {
	repo, err := NewRepository("http://github.com/org/test")
	assertNoError(t, err)
	want := triggersv1.EventListenerTrigger{
		Name: "test",
		Bindings: []*triggersv1.EventListenerBinding{
			{Name: "test-binding"},
		},
		Template: triggersv1.EventListenerTemplate{Name: "test-template"},
		Interceptors: []*triggersv1.EventInterceptor{
			{
				CEL: &triggersv1.CELInterceptor{
					Filter: fmt.Sprintf(githubCDDeployFilters, "org/test"),
				},
			},
			{
				GitHub: &triggersv1.GitHubInterceptor{
					SecretRef: &triggersv1.SecretRef{SecretKey: "webhook-secret-key", SecretName: "secret", Namespace: "ns"},
				},
			},
		},
	}
	got := repo.CreateCDTrigger("test", "secret", "ns", "test-template", []string{"test-binding"})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreateCDTrigger() failed:\n%s", diff)
	}
}

func TestNewGitHubRepository(t *testing.T) {
	tests := []struct {
		url      string
		repoPath string
		errMsg   string
	}{
		{
			"http://github.org",
			"",
			"unable to determine type of Git host from: http://github.org",
		},
		{
			"http://github.com/",
			"",
			"invalid repository URL http://github.com/: path is empty",
		},
		{
			"http://github.com/foo/bar",
			"foo/bar",
			"",
		},
		{
			"https://githuB.com/foo/bar.git",
			"foo/bar",
			"",
		},
		{
			"https://githuB.com/foo/bar/test.git",
			"",
			"invalid repository path for github: /foo/bar/test.git",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(rt *testing.T) {
			repo, err := NewRepository(tt.url)
			if err != nil {
				if diff := cmp.Diff(tt.errMsg, err.Error()); diff != "" {
					rt.Fatalf("repo path errMsg mismatch: \n%s", diff)
				}
			}
			if repo != nil {
				if diff := cmp.Diff(tt.repoPath, repo.(*github).path); diff != "" {
					rt.Fatalf("repo path mismatch: got\n%s", diff)
				}
			}
		})
	}
}
