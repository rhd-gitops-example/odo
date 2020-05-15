package triggers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreatePRBinding(t *testing.T) {
	validPRBinding := triggersv1.TriggerBinding{
		TypeMeta: triggerBindingTypeMeta,
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
	binding := CreatePRBinding("testns")
	if diff := cmp.Diff(validPRBinding, binding); diff != "" {
		t.Fatalf("createPRBinding() failed:\n%s", diff)
	}
}

func TestCreatePushBinding(t *testing.T) {
	validPushBinding := triggersv1.TriggerBinding{
		TypeMeta: triggerBindingTypeMeta,
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
	binding := CreatePushBinding("testns")
	if diff := cmp.Diff(validPushBinding, binding); diff != "" {
		t.Fatalf("CreatePushBinding() failed:\n%s", diff)
	}
}

func TestCreateBindingParam(t *testing.T) {
	validParam := triggersv1.Param{
		Name:  "gitref",
		Value: "$(body.ref)",
	}
	bindingParam := createBindingParam("gitref", "$(body.ref)")
	if diff := cmp.Diff(validParam, bindingParam); diff != "" {
		t.Fatalf("createBindingParam() failed\n%s", diff)
	}
}

func TestCreatePiplineBindingParam(t *testing.T) {
	validParam := pipelinev1.Param{
		Name: "gitref",
		Value: pipelinev1.ArrayOrString{
			StringVal: "$(body.ref)",
			Type:      pipelinev1.ParamTypeString,
		},
	}
	bindingParam := createPipelineBindingParam("gitref", "$(body.ref)")
	if diff := cmp.Diff(validParam, bindingParam); diff != "" {
		t.Fatalf("createPipelineBindingParam() failed\n%s", diff)
	}
}

func TestCreateImageRepoBinding(t *testing.T) {
	imageRepoBinding := triggersv1.TriggerBinding{
		TypeMeta: triggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-binding",
			Namespace: "testns",
		},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "imageRepo",
					Value: "quay.io/user/testing",
				},
			},
		},
	}
	binding := CreateImageRepoBinding("testns", "test-binding", "quay.io/user/testing")
	if diff := cmp.Diff(imageRepoBinding, binding); diff != "" {
		t.Fatalf("CreateImageRepoBinding() failed:\n%s", diff)
	}
}
