package templates

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func TestCreateDevCDDeployTemplate(t *testing.T) {
	validDevCDTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   createTypeMeta("TriggerTemplate", "tekton.dev/v1alpha1"),
		ObjectMeta: createObjectMeta("dev-cd-deploy-from-master-Template"),
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
					RawMessage: createDevCDResourcetemplate(),
				},
			},
		},
	}

	template := createDevCDDeployTemplate()
	if diff := cmp.Diff(validDevCDTemplate, template); diff != "" {
		t.Fatalf("CreateDevCDDeployTemplate failed:\n%s", diff)
	}
}

func TestCreatedevCIBuildPRTemplate(t *testing.T) {
	validdevCIPRTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   createTypeMeta("TriggerTemplate", "tekton.dev/v1alpha1"),
		ObjectMeta: createObjectMeta("dev-ci-build-from-pr-template"),

		Spec: triggersv1.TriggerTemplateSpec{
			Params: []pipelinev1.ParamSpec{
				pipelinev1.ParamSpec{
					Name:        "gitref",
					Description: "The git branch for this PR",
				},
				pipelinev1.ParamSpec{
					Name:        "gitsha",
					Description: "the specific commit SHA.",
				},
				pipelinev1.ParamSpec{
					Name:        "gitrepositoryurl",
					Description: "The git repository url",
				},

				pipelinev1.ParamSpec{
					Name:        "fullname",
					Description: "The GitHub repository for this PullRequest.",
				},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				triggersv1.TriggerResourceTemplate{
					RawMessage: createDevCIResourceTemplate(),
				},
			},
		},
	}
	template := createDevCIBuildPRTemplate()
	if diff := cmp.Diff(validdevCIPRTemplate, template); diff != "" {
		t.Fatalf("CreatedevCIBuildPRTemplate failed:\n%s", diff)
	}
}

func TestCreateStageCDPushTemplate(t *testing.T) {
	ValidStageCDPushTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   createTypeMeta("TriggerTemplate", "tekton.dev/v1alpha1"),
		ObjectMeta: createObjectMeta("stage-cd-deploy-from-push-template"),
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
					RawMessage: createStageCDResourceTemplate(),
				},
			},
		},
	}
	template := createStageCDPushTemplate()
	if diff := cmp.Diff(ValidStageCDPushTemplate, template); diff != "" {
		t.Fatalf("createStageCDPushTemplate failed:\n%s", diff)
	}
}

func TestCreateStageCIdryrunptemplate(t *testing.T) {
	validStageCIdryrunTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   createTypeMeta("TriggerTemplate", "tekton.dev/v1alpha1"),
		ObjectMeta: createObjectMeta("stage-ci-dryrun-from-pr-template"),

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
					RawMessage: createStageCIResourceTemplate(),
				},
			},
		},
	}
	template := createStageCIdryrunptemplate()
	if diff := cmp.Diff(validStageCIdryrunTemplate, template); diff != "" {
		t.Fatalf("createStageCIdryrunptemplate failed:\n%s", diff)
	}

}
