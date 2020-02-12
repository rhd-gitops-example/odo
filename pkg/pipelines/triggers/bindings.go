package triggers

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetTriggerBindings returns a slice of trigger bindings
func GetTriggerBindings() []triggersv1.TriggerBinding {
	return []triggersv1.TriggerBinding{
		createDevCDDeployBinding(),
		createDevCIBuildBinding(),
		createStageCDDeployBinding(),
		createStageCIDryRunBinding(),
	}
}

func createDevCDDeployBinding() triggersv1.TriggerBinding {
	return triggersv1.TriggerBinding{
		TypeMeta:   createTriggerBindingeMeta(),
		ObjectMeta: createObjectMeta("dev-cd-deploy-from-master-binding"),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.head_commit.id)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
			},
		},
	}
}

func createDevCIBuildBinding() triggersv1.TriggerBinding {
	return triggersv1.TriggerBinding{
		TypeMeta:   createTriggerBindingeMeta(),
		ObjectMeta: createObjectMeta("dev-ci-build-from-pr-binding"),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.pull_request.head.ref)"),
				createBindingParam("gitsha", "$(body.pull_request.head.sha)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
				createBindingParam("fullname", "$(body.repository.full_name)"),
			},
		},
	}
}

func createStageCDDeployBinding() triggersv1.TriggerBinding {
	return triggersv1.TriggerBinding{
		TypeMeta:   createTriggerBindingeMeta(),
		ObjectMeta: createObjectMeta("stage-cd-deploy-from-push-binding"),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.ref)"),
				createBindingParam("gitsha", "$(body.commits.0.id)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
			},
		},
	}
}

func createStageCIDryRunBinding() triggersv1.TriggerBinding {
	return triggersv1.TriggerBinding{
		TypeMeta:   createTriggerBindingeMeta(),
		ObjectMeta: createObjectMeta("stage-ci-dryrun-from-pr-binding"),
		Spec: triggersv1.TriggerBindingSpec{
			Params: []pipelinev1.Param{
				createBindingParam("gitref", "$(body.pull_request.head.ref)"),
				createBindingParam("gitrepositoryurl", "$(body.repository.clone_url)"),
			},
		},
	}
}

func createTriggerBindingeMeta() v1.TypeMeta {
	return v1.TypeMeta{
		Kind:       "TriggerBinding",
		APIVersion: "tekton.dev/v1alpha1",
	}
}

func createObjectMeta(name string) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name: name,
	}
}

func createBindingParam(name string, value string) pipelinev1.Param {
	return pipelinev1.Param{
		Name: name,
		Value: pipelinev1.ArrayOrString{
			StringVal: value,
		},
	}
}
