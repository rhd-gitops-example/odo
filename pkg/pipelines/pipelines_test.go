package pipelines

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

func TestCreateDevCIPipeline(t *testing.T) {
	DevCIpipeline := createDevCIPipeline()

	want := &pipelinev1.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: devCIPipelineName,
		},
		Spec: pipelinev1.PipelineSpec{
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					Name: "REPO",
					Type: "string",
				},
				pipelinev1.ParamSpec{
					Name: "COMMIT_SHA",
					Type: "string",
				},
			},
			Resources: []pipelinev1.PipelineDeclaredResource{
				pipelinev1.PipelineDeclaredResource{
					Name: "source-repo",
					Type: "git",
				},
				pipelinev1.PipelineDeclaredResource{
					Name: "runtime-image",
					Type: "image",
				},
			},
			Tasks: []pipelinev1.PipelineTask{
				pipelinev1.PipelineTask{
					Name: "create-pending-status",
					TaskRef: &pipelinev1.TaskRef{
						Name: "create-github-status-task",
					},
					Params: []pipelinev1.Param{
						{Name: "REPO", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "$(params.REPO)"}},
						{Name: "COMMIT_SHA", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "$(params.COMMIT_SHA)"}},
						{Name: "STATE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "pending"}},
						{Name: "DESCRIPTION", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "Starting dev-ci-pipeline"}},
						{Name: "CONTEXT", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "dev-ci-pipeline"}},
					},
				},

				pipelinev1.PipelineTask{
					Name: "build-image",
					TaskRef: &pipelinev1.TaskRef{
						Name: "buildah-task",
					},
					RunAfter: []string{"create-pending-status"},
					Resources: &pipelinev1.PipelineTaskResources{
						Inputs: []pipelinev1.PipelineTaskInputResource{
							{Name: "source",
								Resource: "source-repo"},
						},
						Outputs: []pipelinev1.PipelineTaskOutputResource{
							{Name: "image",
								Resource: "runtime-image"},
						},
					},
				},

				pipelinev1.PipelineTask{
					Name: "create-success-status",
					TaskRef: &pipelinev1.TaskRef{
						Name: "create-github-status-task",
					},
					RunAfter: []string{"build-image"},
					Params: []pipelinev1.Param{
						{Name: "REPO", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "$(params.REPO)"}},
						{Name: "COMMIT_SHA", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "$(params.COMMIT_SHA)"}},
						{Name: "STATE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "success"}},
						{Name: "DESCRIPTION", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "Completed dev-ci-pipeline"}},
						{Name: "CONTEXT", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "dev-ci-pipeline"}},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, DevCIpipeline); diff != "" {
		t.Fatalf("TestCreateDevCIPipeline() failed got\n%s", diff)
	}
}

func TestCreateStageCIPipeline(t *testing.T) {
	stageCIpipeline := createStageCIPipeline("Testprefix")
	want := &pipelinev1.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "stage-ci-pipeline",
		},
		Spec: pipelinev1.PipelineSpec{
			Resources: []pipelinev1.PipelineDeclaredResource{
				pipelinev1.PipelineDeclaredResource{
					Name: "source-repo",
					Type: "git",
				},
			},

			Tasks: []pipelinev1.PipelineTask{
				pipelinev1.PipelineTask{
					Name: "apply-source",
					TaskRef: &pipelinev1.TaskRef{
						Name: "deploy-from-source-task",
					},
					Resources: &pipelinev1.PipelineTaskResources{
						Inputs: []pipelinev1.PipelineTaskInputResource{
							{Name: "source",
								Resource: "source-repo"},
						},
					},
					Params: []pipelinev1.Param{
						{Name: "NAMESPACE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "Testprefix-stage-environment"}},
						{Name: "DRYRUN", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "true"}},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, stageCIpipeline); diff != "" {
		t.Fatalf("TestcreateStageCIPipeline() failed got\n%s", diff)
	}

}

func TestCreateDevCDPipeline(t *testing.T) {
	DevCDpipeline := createDevCDPipeline("Testprefix", "usr/path/")

	want := &pipelinev1.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "dev-cd-pipeline",
		},

		Spec: pipelinev1.PipelineSpec{
			Resources: []pipelinev1.PipelineDeclaredResource{
				pipelinev1.PipelineDeclaredResource{
					Name: "source-repo",
					Type: "git",
				},
				pipelinev1.PipelineDeclaredResource{
					Name: "runtime-image",
					Type: "image",
				},
			},
			Tasks: []pipelinev1.PipelineTask{

				pipelinev1.PipelineTask{
					Name: "build-image",
					TaskRef: &pipelinev1.TaskRef{
						Name: "buildah-task",
					},

					Resources: &pipelinev1.PipelineTaskResources{
						Inputs: []pipelinev1.PipelineTaskInputResource{
							{Name: "source",
								Resource: "source-repo"},
						},
						Outputs: []pipelinev1.PipelineTaskOutputResource{
							{Name: "image",
								Resource: "runtime-image"},
						},
					},
				},

				pipelinev1.PipelineTask{
					Name: "deploy-image",
					TaskRef: &pipelinev1.TaskRef{
						Name: "deploy-using-kubectl-task",
					},
					RunAfter: []string{"build-image"},
					Resources: &pipelinev1.PipelineTaskResources{
						Inputs: []pipelinev1.PipelineTaskInputResource{
							{Name: "source", Resource: "source-repo"},
							{Name: "image", Resource: "runtime-image"},
						},
					},
					Params: []pipelinev1.Param{
						{Name: "PATHTODEPLOYMENT", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "usr/path/"}},
						{Name: "YAMLPATHTOIMAGE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "spec.template.spec.containers[0].image"}},
						{Name: "NAMESPACE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "Testprefix-dev-environment"}},
					},
				},
			},
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					Name: "REPO",
					Type: "string",
				},
				pipelinev1.ParamSpec{
					Name: "COMMIT_SHA",
					Type: "string",
				},
			},
		},
	}

	if diff := cmp.Diff(want, DevCDpipeline); diff != "" {
		t.Fatalf("TestCreateDevCDPipeline() failed got\n%s", diff)
	}
}

func TestCreateStageCDPipeline(t *testing.T) {
	stageCDpipeline := createStageCDPipeline("Testprefix")
	want := &pipelinev1.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "stage-cd-pipeline",
		},
		Spec: pipelinev1.PipelineSpec{
			Resources: []pipelinev1.PipelineDeclaredResource{
				pipelinev1.PipelineDeclaredResource{
					Name: "source-repo",
					Type: "git",
				},
			},

			Tasks: []pipelinev1.PipelineTask{
				pipelinev1.PipelineTask{
					Name: "apply-source",
					TaskRef: &pipelinev1.TaskRef{
						Name: "deploy-from-source-task",
					},
					Resources: &pipelinev1.PipelineTaskResources{
						Inputs: []pipelinev1.PipelineTaskInputResource{
							{Name: "source",
								Resource: "source-repo"},
						},
					},
					Params: []pipelinev1.Param{
						{Name: "NAMESPACE", Value: pipelinev1.ArrayOrString{Type: "string", StringVal: "Testprefix-stage-environment"}},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, stageCDpipeline); diff != "" {
		t.Fatalf("TestcreateStageCSPipeline() failed got\n%s", diff)
	}

}
