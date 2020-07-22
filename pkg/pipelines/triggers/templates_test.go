package triggers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/odo/pkg/pipelines/meta"
)

const (
	serviceAccName = "pipeline"
)

func TestCreateDevCDDeployTemplate(t *testing.T) {
	validDevCDTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   triggerTemplateTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testns", "app-cd-template")),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []triggersv1.ParamSpec{
				{
					Name:        GitCommitID,
					Description: "The specific commit SHA.",
				},
				{
					Name:        "gitrepositoryurl",
					Description: "The git repository url",
				},
			},

			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				{
					RawExtension: runtime.RawExtension{
						Raw: createDevCDResourceTemplate(serviceAccName),
					},
				},
			},
		},
	}

	template := CreateDevCDDeployTemplate("testns", serviceAccName)
	if diff := cmp.Diff(validDevCDTemplate, template); diff != "" {
		t.Fatalf("CreateDevCDDeployTemplate failed:\n%s", diff)
	}
}

func TestCreateDevCIBuildPRTemplate(t *testing.T) {
	validdevCIPRTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   triggerTemplateTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testns", "app-ci-template")),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []triggersv1.ParamSpec{
				{
					Name:        GitRef,
					Description: "The git branch for this PR.",
				},
				{
					Name:        GitCommitID,
					Description: "the specific commit SHA.",
				},
				{
					Name:        GitCommitDate,
					Description: "The date at which the commit was made",
				},
				{
					Name:        GitCommitAuthor,
					Description: "The name of the github user handle that made the commit",
				},
				{
					Name:        GitCommitMessage,
					Description: "The commit message",
				},
				{
					Name:        "gitrepositoryurl",
					Description: "The git repository URL.",
				},
				{
					Name:        "fullname",
					Description: "The GitHub repository for this PullRequest.",
				},
				{
					Name:        "imageRepo",
					Description: "The repository to push built images to.",
				},
				{
					Name:        "tlsVerify",
					Description: "Enable image repostiory TLS certification verification.",
				},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				{
					RawExtension: runtime.RawExtension{
						Raw: createDevCIResourceTemplate(serviceAccName),
					},
				},
			},
		},
	}
	template := CreateDevCIBuildPRTemplate("testns", serviceAccName)
	if diff := cmp.Diff(validdevCIPRTemplate, template); diff != "" {
		t.Fatalf("CreatedevCIBuildPRTemplate failed:\n%s", diff)
	}
}

func TestCreateCDPushTemplate(t *testing.T) {
	ValidStageCDPushTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   triggerTemplateTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testns", "cd-deploy-from-push-template")),
		Spec: triggersv1.TriggerTemplateSpec{
			Params: []triggersv1.ParamSpec{
				{
					Name:        GitRef,
					Description: "The git revision",
					Default:     strPtr("master"),
				},
				{
					Name:        GitCommitDate,
					Description: "The date at which the commit was made",
				},
				{
					Name:        GitCommitAuthor,
					Description: "The name of the github user handle that made the commit",
				},
				{
					Name:        GitCommitMessage,
					Description: "The commit message",
				},
				{
					Name:        "gitrepositoryurl",
					Description: "The git repository url",
				},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				{
					RawExtension: runtime.RawExtension{
						Raw: createCDResourceTemplate(serviceAccName),
					},
				},
			},
		},
	}
	template := CreateCDPushTemplate("testns", serviceAccName)
	if diff := cmp.Diff(ValidStageCDPushTemplate, template); diff != "" {
		t.Fatalf("createCDPushTemplate failed:\n%s", diff)
	}
}

func TestCreateCIDryRunTemplate(t *testing.T) {
	validStageCIDryRunTemplate := triggersv1.TriggerTemplate{
		TypeMeta:   triggerTemplateTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testns", "ci-dryrun-from-push-template")),

		Spec: triggersv1.TriggerTemplateSpec{
			Params: []triggersv1.ParamSpec{
				{Name: GitRef, Description: "The git revision", Default: strPtr("master")},
				{Name: "io.openshift.build.commit.id", Description: "The specific commit SHA"},
				{Name: "gitrepositoryurl", Description: "The git repository url"},
			},
			ResourceTemplates: []triggersv1.TriggerResourceTemplate{
				{
					RawExtension: runtime.RawExtension{
						Raw: createCIResourceTemplate(serviceAccName),
					},
				},
			},
		},
	}
	template := CreateCIDryRunTemplate("testns", serviceAccName)
	if diff := cmp.Diff(validStageCIDryRunTemplate, template); diff != "" {
		t.Fatalf("createCIdryrunptemplate failed:\n%s", diff)
	}

}

func TestStatusTrackerAnnotations(t *testing.T) {
	om := meta.ObjectMeta(meta.NamespacedName("test-ns", "name"))
	statusTrackerAnnotations("my-pipeline", "description")(&om)

	want := v1.ObjectMeta{
		Name:      "name",
		Namespace: "test-ns",
		Annotations: map[string]string{
			"tekton.dev/git-status":         "true",
			"tekton.dev/status-context":     "my-pipeline",
			"tekton.dev/status-description": "description",
		},
	}

	if diff := cmp.Diff(want, om); diff != "" {
		t.Fatalf("statusTrackerAnnotations failed:\n%s", diff)
	}
}
