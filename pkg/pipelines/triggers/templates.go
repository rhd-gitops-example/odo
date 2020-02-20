package triggers

import (
	"encoding/json"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDevCDDeployTemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("dev-cd-deploy-from-master-Template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{

				createTemplateParamSpecDefault("gitref", "The git revision", "master"),
				createTemplateParamSpec("gitrepositoryurl", "The git repository url"),
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createDevCDResourcetemplate(),
				},
			},
		},
	}
}

func createdevCIBuildPRTemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("dev-ci-build-from-pr-template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{

				createTemplateParamSpec("gitref", "The git branch for this PR"),
				createTemplateParamSpec("gitsha", "the specific commit SHA."),
				createTemplateParamSpec("gitrepositoryurl", "The git repository url"),
				createTemplateParamSpec("fullname", "The GitHub repository for this PullRequest."),
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createDevCIResourceTemplate(),
				},
			},
		},
	}

}

func createStageCDPushTemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("stage-cd-deploy-from-push-template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{

				createTemplateParamSpecDefault("gitref", "The git revision", "master"),
				createTemplateParamSpec("gitrepositoryurl", "The git repository url"),
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createStageCDResourceTemplate(),
				},
			},
		},
	}
}

func createStageCIdryrunptemplate() triggersv1.TriggerTemplate {
	return triggersv1.TriggerTemplate{
		TypeMeta:   createTriggerTemplateMeta(),
		ObjectMeta: createObjectMeta("stage-ci-dryrun-from-pr-template"),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{

				createTemplateParamSpecDefault("gitref", "The git revision", "master"),
				createTemplateParamSpec("gitrepositoryurl", "The git repository url"),
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createStageCIResourceTemplate(),
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

func createTemplateParamSpecDefault(name string, description string, value string) pipelinev1.ParamSpec {
	return pipelinev1.ParamSpec{
		Name:        name,
		Description: description,
		Default: &pipelinev1.ArrayOrString{
			StringVal: value,
		},
	}
}

func createTemplateParamSpec(name string, description string) pipelinev1.ParamSpec {
	return pipelinev1.ParamSpec{
		Name:        name,
		Description: description,
	}

}

func createDevCDResourcetemplate() []byte {
	byteTemplate, _ := json.Marshal(createDevCDPipelineRun())
	return []byte(string(byteTemplate))

}
func createDevCIResourceTemplate() []byte {
	byteTemplateCI, _ := json.Marshal(createDevCIPipelineRun())
	return []byte(string(byteTemplateCI))

}
func createStageCDResourceTemplate() []byte {
	byteStageCD, _ := json.Marshal(createStageCDPipelineRun())
	return []byte(string(byteStageCD))
}
func createStageCIResourceTemplate() []byte {
	byteStageCI, _ := json.Marshal(createStageCIPipelineRun())
	return []byte(string(byteStageCI))
}
