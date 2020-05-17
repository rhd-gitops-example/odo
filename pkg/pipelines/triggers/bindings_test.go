package triggers

import (
	"fmt"
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
			Params: []pipelinev1.Param{
				{
					Name: "gitref",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.pull_request.head.ref)",
						Type:      pipelinev1.ParamTypeString,
					},
				},
				{
					Name: "gitsha",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.pull_request.head.sha)",
						Type:      pipelinev1.ParamTypeString,
					},
				},
				{
					Name: "gitrepositoryurl",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.repository.clone_url)",
						Type:      pipelinev1.ParamTypeString,
					},
				},
				{
					Name: "fullname",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.repository.full_name)",
						Type:      pipelinev1.ParamTypeString,
					},
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
			Params: []pipelinev1.Param{
				{
					Name: "gitref",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.ref)",
						Type:      pipelinev1.ParamTypeString,
					},
				},
				{
					Name: "gitsha",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.head_commit.id)",
						Type:      pipelinev1.ParamTypeString,
					},
				},
				{
					Name: "gitrepositoryurl",
					Value: pipelinev1.ArrayOrString{
						StringVal: "$(body.repository.clone_url)",
						Type:      pipelinev1.ParamTypeString,
					},
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
	validParam := pipelinev1.Param{
		Name: "gitref",
		Value: pipelinev1.ArrayOrString{
			StringVal: "$(body.ref)",
			Type:      pipelinev1.ParamTypeString,
		},
	}
	bindingParam := createBindingParam("gitref", "$(body.ref)")
	if diff := cmp.Diff(validParam, bindingParam); diff != "" {
		t.Fatalf("createBindingParam() failed\n%s", diff)
	}
}

func TestCreateImageRepoBinding(t *testing.T) {

	testcases := []struct {
		ns        string
		name      string
		imageRepo string
		tlsVerify string
		want      triggersv1.TriggerBinding
	}{
		{
			ns:        "testns",
			name:      "test-binding",
			imageRepo: "quay.io/user/testing",
			tlsVerify: "true",
			want: triggersv1.TriggerBinding{
				TypeMeta: triggerBindingTypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-binding",
					Namespace: "testns",
				},
				Spec: triggersv1.TriggerBindingSpec{
					Params: []pipelinev1.Param{
						{
							Name: "imageRepo",
							Value: pipelinev1.ArrayOrString{
								StringVal: "quay.io/user/testing",
								Type:      pipelinev1.ParamTypeString,
							},
						},
						{
							Name: "tlsVerify",
							Value: pipelinev1.ArrayOrString{
								StringVal: "true",
								Type:      pipelinev1.ParamTypeString,
							},
						},
					},
				},
			},
		},
		{
			ns:        "testns",
			name:      "test-binding",
			imageRepo: "proj/taxi",
			tlsVerify: "false",
			want: triggersv1.TriggerBinding{
				TypeMeta: triggerBindingTypeMeta,
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-binding",
					Namespace: "testns",
				},
				Spec: triggersv1.TriggerBindingSpec{
					Params: []pipelinev1.Param{
						{
							Name: "imageRepo",
							Value: pipelinev1.ArrayOrString{
								StringVal: "proj/taxi",
								Type:      pipelinev1.ParamTypeString,
							},
						},
						{
							Name: "tlsVerify",
							Value: pipelinev1.ArrayOrString{
								StringVal: "false",
								Type:      pipelinev1.ParamTypeString,
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got := CreateImageRepoBinding(tt.ns, tt.name, tt.imageRepo, tt.tlsVerify)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("CreateImageRepoBinding() failed:\n%s", diff)
			}
		})
	}
}
