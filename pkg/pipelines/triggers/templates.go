package triggers

import (
	"encoding/json"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDevCDDeployTemplate() *triggersv1.TriggerTemplate {
	return &triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("dev-cd-deploy-from-master-Template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					createTemplateParamSpec("gitref", "The git revision", "master"),
					createTemplateParamSpec("gitrepositoryurl", "the git repository url"),
				},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createTemplateResource(),
				},
			},
		},
	}
}

func createTriggerTemplateMeta() v1.TypeMeta {
	return v1.TypeMeta{
		Kind:       "TriggerTemplate",
		APIVersion: "tekton.dev/v1alpha1",
	}
}

func createTemplateParamSpec(name string, description string, value string) triggersv1.ParamSpec {
	return pipelinev1.ParamSpec{
		Name:        name,
		Description: description,
		Default: &pipelinev1.ArrayOrString{
			StringVal: value,
		},
	}
}

func createTemplateResource() json.RawMessage {
	return json.RawMessage(
		[]byte(`{
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
		
		
		}`))
}

func CreatedevCIBuildPRTemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("dev-cd-deploy-from-master-Template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					createTemplateParamSpec("gitref", "The git revision", "master"),
					createTemplateParamSpec("gitrepositoryurl", "the git repository url"),
				},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createTemplateResource(),
				},
			},
		},
	}

}

func validdevCIdryrunTemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{}
}
