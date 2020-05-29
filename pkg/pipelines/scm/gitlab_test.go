package scm

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPRbindingForGitlab(t *testing.T) {
	repo, err := NewGitLabRepository("http://gitlab.com/org/test")
	assertNoError(t, err)
	want := triggersv1.TriggerBinding{
		TypeMeta: triggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "gitlab-pr-binding",
			Namespace: "testns",
		},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "gitref",
					Value: "$(body.object_attributes.source_branch)",
				},
				{
					Name:  "gitsha",
					Value: "$(body.object_attributes.last_commit.id)",
				},
				{
					Name:  "gitrepositoryurl",
					Value: "$(body.project.git_http_url)",
				},
				{
					Name:  "fullname",
					Value: "$(body.project.path_with_namespace)",
				},
			},
		},
	}
	got, name := repo.CreatePRBinding("testns")
	if name != gitlabPRBindingName {
		t.Fatalf("CreatePushBinding() returned a wrong binding: want %v got %v", gitlabPRBindingName, name)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("createPRBinding() failed:\n%s", diff)
	}
}

func TestCreatePushBindingForGitlab(t *testing.T) {
	repo, err := NewGitLabRepository("https://gitlab.com/org/fullname/subgroup/repository/subrepo/test")
	assertNoError(t, err)
	want := triggersv1.TriggerBinding{
		TypeMeta: triggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "gitlab-push-binding",
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
					Value: "$(body.after)",
				},
				{
					Name:  "gitrepositoryurl",
					Value: "$(body.project.git_http_url)",
				},
			},
		},
	}
	got, name := repo.CreatePushBinding("testns")
	if name != gitlabPushBindingName {
		t.Fatalf("CreatePushBinding() returned a wrong binding: want %v got %v", gitlabPushBindingName, name)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreatePushBinding() failed:\n%s", diff)
	}
}

func TestCreateCITriggerForGitLab(t *testing.T) {
	repo, err := NewGitLabRepository("http://gitlab.com/org/test")
	assertNoError(t, err)
	want := triggersv1.EventListenerTrigger{
		Name: "test",
		Bindings: []*triggersv1.EventListenerBinding{
			&triggersv1.EventListenerBinding{Name: "test-binding"},
		},
		Template: triggersv1.EventListenerTemplate{Name: "test-template"},
		Interceptors: []*triggersv1.EventInterceptor{
			&triggersv1.EventInterceptor{
				CEL: &triggersv1.CELInterceptor{
					Filter: fmt.Sprintf(gitlabCIDryRunFilters, "org/test"),
				},
			},
			&triggersv1.EventInterceptor{
				GitLab: &triggersv1.GitLabInterceptor{
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

func TestCreateCDTriggersForGitLab(t *testing.T) {
	repo, err := NewGitLabRepository("http://gitlab.com/org/test")
	assertNoError(t, err)
	want := triggersv1.EventListenerTrigger{
		Name: "test",
		Bindings: []*triggersv1.EventListenerBinding{
			&triggersv1.EventListenerBinding{Name: "test-binding"},
		},
		Template: triggersv1.EventListenerTemplate{Name: "test-template"},
		Interceptors: []*triggersv1.EventInterceptor{
			&triggersv1.EventInterceptor{
				CEL: &triggersv1.CELInterceptor{
					Filter: fmt.Sprintf(gitlabCDDeployFilters, "org/test"),
				},
			},
			&triggersv1.EventInterceptor{
				GitLab: &triggersv1.GitLabInterceptor{
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

func TestNewGitlabRepository(t *testing.T) {
	tests := []struct {
		url      string
		repoPath string
		errMsg   string
	}{
		{
			"http://gitlab.com",
			"",
			"unable to determine repo path from: http://gitlab.com",
		},
		{
			"http://gitlab.com/",
			"",
			"unable to determine repo path from: http://gitlab.com/",
		},
		{
			"http://gitlab.com/foo/bar",
			"foo/bar",
			"",
		},
		{
			"https://gitlab.com/group/subgroup/subgroup/repo.git",
			"group/subgroup/subgroup/repo",
			"",
		},
		{
			"https://gitlaB.com/foo/bar/test.git",
			"foo/bar/test",
			"",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(rt *testing.T) {
			repo, err := NewGitLabRepository(tt.url)
			if err != nil {
				if diff := cmp.Diff(tt.errMsg, err.Error()); diff != "" {
					rt.Fatalf("repo path errMsg mismatch: \n%s", diff)
				}
			}
			if repo != nil {
				if diff := cmp.Diff(tt.repoPath, repo.path); diff != "" {
					rt.Fatalf("repo path mismatch: got\n%s", diff)
				}
			}
		})
	}
}
