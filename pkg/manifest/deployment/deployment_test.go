package deployment

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/manifest/meta"
)

func TestCreate(t *testing.T) {
	component := "my-component"
	image := "quay.io/testing/testing"
	d := Create("testing", component, image)

	want := &appsv1.Deployment{
		TypeMeta:   meta.TypeMeta("Deployment", "apps/v1"),
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testing", component)),
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr32(1),
			Selector: labelSelector("name", component),
			Template: podTemplate(component, image),
		},
	}

	if diff := cmp.Diff(want, d); diff != "" {
		t.Fatalf("deployment diff:\n%s", diff)
	}
}

func TestDefaultPodTemplate(t *testing.T) {
	component := "test-svc"
	image := "quay.io/example/example"
	want := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"name": component,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			Containers: []corev1.Container{
				{
					Name:            component,
					Image:           image,
					ImagePullPolicy: corev1.PullAlways,
				},
			},
		},
	}

	spec := podTemplate(component, image)

	if diff := cmp.Diff(want, spec); diff != "" {
		t.Fatalf("labelTemplate diff: %s", diff)
	}
}
