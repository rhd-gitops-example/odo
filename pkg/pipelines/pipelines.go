package pipelines

import (
	"github.com/openshift/odo/pkg/pipelines/meta"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

const devCIPipelineName = "dev-ci-pipeline"

// createCIPipeline creates the dev-ci-pipeline.
func createDevCIPipeline(nss map[string]string) *pipelinev1.Pipeline {
	typeMeta := meta.TypeMeta("Pipeline", "tekton.dev/v1alpha1")
	objectMeta := meta.CreateObjectMeta(nss["cicd"], devCIPipelineName)
	return &pipelinev1.Pipeline{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Spec: pipelinev1.PipelineSpec{
			Params: []pipelinev1.ParamSpec{
				createParamSpec("REPO", "string"),
				createParamSpec("COMMIT_SHA", "string"),
			},

			Resources: []pipelinev1.PipelineDeclaredResource{
				createPipelineDeclaredResource("source-repo", "git"),
				createPipelineDeclaredResource("runtime-image", "image"),
			},

			Tasks: []pipelinev1.PipelineTask{
				createGitHubStatusTask("create-pending-status", "", "pending", "Starting dev-ci-pipeline"),
				createBuildImageTask("build-image"),
				createGitHubStatusTask("create-success-status", "build-image", "success", "Completed dev-ci-pipeline"),
			},
		},
	}
}

func createStageCIPipeline(prefix string, nss map[string]string) *pipelinev1.Pipeline {
	return &pipelinev1.Pipeline{
		TypeMeta:   meta.TypeMeta("Pipeline", "tekton.dev/v1alpha1"),
		ObjectMeta: meta.CreateObjectMeta(nss["cicd"], "stage-ci-pipeline"),
		Spec: pipelinev1.PipelineSpec{

			Resources: []pipelinev1.PipelineDeclaredResource{
				createPipelineDeclaredResource("source-repo", "git"),
			},

			Tasks: []pipelinev1.PipelineTask{
				createStageCIPipelineTask("apply-source", prefix),
			},
		},
	}
}
func createDevCDPipeline(prefix, deploymentPath string, nss map[string]string) *pipelinev1.Pipeline {
	return &pipelinev1.Pipeline{
		TypeMeta:   meta.TypeMeta("Pipeline", "tekton.dev/v1alpha1"),
		ObjectMeta: meta.CreateObjectMeta(nss["cicd"], "dev-cd-pipeline"),
		Spec: pipelinev1.PipelineSpec{
			Resources: []pipelinev1.PipelineDeclaredResource{
				createPipelineDeclaredResource("source-repo", "git"),
				createPipelineDeclaredResource("runtime-image", "image"),
			},
			Params: []pipelinev1.ParamSpec{
				createParamSpec("REPO", "string"),
				createParamSpec("COMMIT_SHA", "string"),
			},

			Tasks: []pipelinev1.PipelineTask{
				createDevCDBuildImageTask("build-image"),
				createDevCDDeployImageTask("deploy-image", prefix, deploymentPath),
			},
		},
	}
}

func createStageCDPipeline(prefix string, nss map[string]string) *pipelinev1.Pipeline {
	return &pipelinev1.Pipeline{
		TypeMeta:   meta.TypeMeta("Pipeline", "tekton.dev/v1alpha1"),
		ObjectMeta: meta.CreateObjectMeta(nss["cicd"], "stage-cd-pipeline"),
		Spec: pipelinev1.PipelineSpec{

			Resources: []pipelinev1.PipelineDeclaredResource{
				createPipelineDeclaredResource("source-repo", "git"),
			},

			Tasks: []pipelinev1.PipelineTask{
				createStageCDPipelineTask("apply-source", prefix),
			},
		},
	}
}

func createParamSpec(name string, paramType pipelinev1.ParamType) pipelinev1.ParamSpec {
	return pipelinev1.ParamSpec{Name: name, Type: paramType}
}

func createPipelineDeclaredResource(name string, resourceType string) pipelinev1.PipelineDeclaredResource {
	return pipelinev1.PipelineDeclaredResource{Name: name, Type: resourceType}
}

func createGitHubStatusTask(name, runAfter, state, description string) pipelinev1.PipelineTask {
	t := pipelinev1.PipelineTask{
		Name:    name,
		TaskRef: createTaskRef("create-github-status-task"),
		Params: []pipelinev1.Param{
			createTaskParam("REPO", "$(params.REPO)"),
			createTaskParam("COMMIT_SHA", "$(params.COMMIT_SHA)"),
			createTaskParam("STATE", state),
			createTaskParam("DESCRIPTION", description),
			createTaskParam("CONTEXT", "dev-ci-pipeline"),
		},
	}
	if runAfter != "" {
		t.RunAfter = []string{runAfter}
	}
	return t
}

func createBuildImageTask(name string) pipelinev1.PipelineTask {
	return pipelinev1.PipelineTask{
		Name:     name,
		TaskRef:  createTaskRef("buildah-task"),
		RunAfter: []string{"create-pending-status"},
		Resources: &pipelinev1.PipelineTaskResources{
			Inputs:  []pipelinev1.PipelineTaskInputResource{createInputTaskResource("source", "source-repo")},
			Outputs: []pipelinev1.PipelineTaskOutputResource{createOutputTaskResource("image", "runtime-image")},
		},
	}

}

func createDevCDBuildImageTask(name string) pipelinev1.PipelineTask {
	return pipelinev1.PipelineTask{
		Name:    name,
		TaskRef: createTaskRef("buildah-task"),
		Resources: &pipelinev1.PipelineTaskResources{
			Inputs:  []pipelinev1.PipelineTaskInputResource{createInputTaskResource("source", "source-repo")},
			Outputs: []pipelinev1.PipelineTaskOutputResource{createOutputTaskResource("image", "runtime-image")},
		},
	}
}

func createDevCDDeployImageTask(name, prefix, deploymentPath string) pipelinev1.PipelineTask {
	return pipelinev1.PipelineTask{
		Name:     name,
		TaskRef:  createTaskRef("deploy-using-kubectl-task"),
		RunAfter: []string{"build-image"},
		Resources: &pipelinev1.PipelineTaskResources{
			Inputs: []pipelinev1.PipelineTaskInputResource{
				createInputTaskResource("source", "source-repo"),
				createInputTaskResource("image", "runtime-image"),
			},
		},
		Params: []pipelinev1.Param{
			createTaskParam("PATHTODEPLOYMENT", deploymentPath),
			createTaskParam("YAMLPATHTOIMAGE", "spec.template.spec.containers[0].image"),
			createTaskParam("NAMESPACE", prefix+"-dev-environment"),
		},
	}
}

func createStageCIPipelineTask(name, prefix string) pipelinev1.PipelineTask {
	return pipelinev1.PipelineTask{
		Name:    name,
		TaskRef: createTaskRef("deploy-from-source-task"),
		Resources: &pipelinev1.PipelineTaskResources{
			Inputs: []pipelinev1.PipelineTaskInputResource{createInputTaskResource("source", "source-repo")},
		},
		Params: []pipelinev1.Param{
			createTaskParam("NAMESPACE", prefix+"-stage-environment"),
			createTaskParam("DRYRUN", "true"),
		},
	}

}

func createStageCDPipelineTask(name, prefix string) pipelinev1.PipelineTask {
	return pipelinev1.PipelineTask{
		Name:    name,
		TaskRef: createTaskRef("deploy-from-source-task"),
		Resources: &pipelinev1.PipelineTaskResources{
			Inputs: []pipelinev1.PipelineTaskInputResource{createInputTaskResource("source", "source-repo")},
		},
		Params: []pipelinev1.Param{
			createTaskParam("NAMESPACE", prefix+"-stage-environment"),
		},
	}

}

func createTaskParam(name, value string) pipelinev1.Param {
	return pipelinev1.Param{
		Name: name,

		Value: pipelinev1.ArrayOrString{
			Type:      pipelinev1.ParamTypeString,
			StringVal: value,
		},
	}
}

func createTaskRef(name string) *pipelinev1.TaskRef {
	return &pipelinev1.TaskRef{
		Name: name,
	}
}

func createInputTaskResource(name string, resource string) pipelinev1.PipelineTaskInputResource {
	return pipelinev1.PipelineTaskInputResource{
		Name:     name,
		Resource: resource,
	}

}

func createOutputTaskResource(name string, resource string) pipelinev1.PipelineTaskOutputResource {
	return pipelinev1.PipelineTaskOutputResource{
		Name:     name,
		Resource: resource,
	}
}
