package pipelines

import (
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

const devCIPipelineName = "dev-ci-pipeline"

func CreateCIPipeline() *pipelinev1.Pipeline {
	typeMeta := createTypeMeta("Pipeline", "tekton.dev/v1alpha1")
	objectMeta := createObjectMeta(devCIPipelineName)
	return &pipelinev1.Pipeline{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Spec: pipelinev1.PipelineSpec{
			Params: []pipelinev1.ParamSpec{
				createParamSpec("REPO", "string"),
				createParamSpec("COMMIT_SHA", "string"),
			},
		},
	}
}

func createParamSpec(name string, paramType pipelinev1.ParamType) pipelinev1.ParamSpec {
	return pipelinev1.ParamSpec{Name: name, Type: paramType}
}
