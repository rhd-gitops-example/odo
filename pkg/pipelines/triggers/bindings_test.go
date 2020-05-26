package triggers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	imageRepoBinding := triggersv1.TriggerBinding{
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
			},
		},
	}
	binding := CreateImageRepoBinding("testns", "test-binding", "quay.io/user/testing")
	if diff := cmp.Diff(imageRepoBinding, binding); diff != "" {
		t.Fatalf("CreateImageRepoBinding() failed:\n%s", diff)
	}
}
