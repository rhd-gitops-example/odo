package triggers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCreateDevCDDeployTemplate(t *testing.T) {
	validDevCDTemplate := triggersv1.TriggerTemplate{

		TypeMeta: v1.TypeMeta{
			Kind:       "TriggerTemplate",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: "dev-cd-deploy-from-master-template",
		},
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					Name:        "gitref",
					Description: "The git revision",
					Default: &pipelinev1.ArrayOrString{
						StringVal: "master",
					},
				},
				pipelinev1.ParamSpec{
					Name:        "gitrepositoryurl",
					Description: "The git repository url",
				},
			},

			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawExtension: createResourcetemplate(),
				},
			},
		},
	}

	template := CreateDevCDDeployTemplate()
	if diff := cmp.Diff(validDevCDTemplate, template); diff != "" {
		t.Fatalf("CreateDevCDDeployTemplate failed:\n%s", diff)
	}
}

func createResourcetemplate() runtime.RawExtension {
	return runtime.RawExtension{
		Raw: byte(`{
			"apiVersion": "tekton.dev/v1alpha1"
			"kind": "PipelineRun",
			"metadata":
			 { "name": "dev-cd-pipeline-run-$(uid)"},
			"spec":
			  {{"serviceAccountName": "demo-sa"}
			  "pipelineRef":
				{"name": "dev-cd-pipeline"}
			  "resources":
				{"name": "source-repo"}
				  "resourceSpec":
					{"type": "git"}
					"params":
						{"name": "revision"
						"value": "$(params.gitref)"
						"name": "url"
					  	"value": "$(params.gitrepositoryurl)"}
				"name": "runtime-image"
				  "resourceSpec":
					{"type": "image"}
					"params":
						{"name": url
						value: REPLACE_IMAGE:$(params.gitref)}}
		
		
		}`),
	}
}
