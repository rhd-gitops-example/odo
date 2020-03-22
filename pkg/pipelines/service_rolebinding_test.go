package pipelines

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/pipelines/meta"
	corev1 "k8s.io/api/core/v1"
	v1rbac "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRoleBinding(t *testing.T) {
	want := &v1rbac.RoleBinding{
		TypeMeta: roleBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []v1rbac.Subject{
			v1rbac.Subject{
				Kind:      "ServiceAccount",
				Name:      "pipeline",
				Namespace: "testing",
			},
		},
		RoleRef: v1rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		},
	}
	sa := &corev1.ServiceAccount{
		TypeMeta:   serviceAccountTypeMeta,
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("testing", "pipeline")),
	}
	roleBindingTask := createRoleBinding(
		meta.NamespacedName("", roleBindingName),
		sa, "Role", roleName)
	if diff := cmp.Diff(want, roleBindingTask); diff != "" {
		t.Errorf("TestRoleBinding() failed:\n%s", diff)
	}

}

func TestRoleBindingForSubjects(t *testing.T) {
	want := &v1rbac.RoleBinding{
		TypeMeta: roleBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: "testns",
		},
		Subjects: []v1rbac.Subject{
			v1rbac.Subject{
				Kind:      "ServiceAccount",
				Name:      "pipeline",
				Namespace: "testing",
			},
			v1rbac.Subject{
				Kind:      "ServiceAccount",
				Name:      "pipeline",
				Namespace: "testing2",
			},
		},
		RoleRef: v1rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		},
	}

	roleBinding := createRoleBindingForSubjects(meta.NamespacedName("testns", roleBindingName), "Role", roleName,
		[]v1rbac.Subject{v1rbac.Subject{Kind: "ServiceAccount", Name: "pipeline", Namespace: "testing"},
			v1rbac.Subject{Kind: "ServiceAccount", Name: "pipeline", Namespace: "testing2"},
		})

	if diff := cmp.Diff(want, roleBinding); diff != "" {
		t.Errorf("TestRoleBindingForSubjects() failed:\n%s", diff)
	}

}

func TestCreateRole(t *testing.T) {
	want := &v1rbac.Role{
		TypeMeta: roleTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name: roleName,
		},
		Rules: []v1rbac.PolicyRule{
			v1rbac.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"namespace"},
				Verbs:     []string{"patch"},
			},
		},
	}
	roleTask := createRole(meta.NamespacedName("", roleName), rules)
	if diff := cmp.Diff(roleTask, want); diff != "" {
		t.Errorf("TestCreateRole() failed:\n%s", diff)
	}
}

func TestServiceAccount(t *testing.T) {
	want := &corev1.ServiceAccount{
		TypeMeta: serviceAccountTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name: "pipeline",
		},
		Secrets: []corev1.ObjectReference{
			corev1.ObjectReference{
				Name: "regcred",
			},
		},
	}
	servicetask := createServiceAccount(meta.NamespacedName("", "pipeline"))
	servicetask = addSecretToSA(servicetask, "regcred")
	if diff := cmp.Diff(servicetask, want); diff != "" {
		t.Errorf("TestServiceAccount() failed:\n%s", diff)
	}
}

func TestAddSecretToSA(t *testing.T) {
	validSA := &corev1.ServiceAccount{
		TypeMeta: serviceAccountTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name: "pipeline",
		},
	}
	validSecrets := []corev1.ObjectReference{
		corev1.ObjectReference{
			Name: "regcred",
		},
	}
	sa := addSecretToSA(validSA, "regcred")
	if diff := cmp.Diff(sa.Secrets, validSecrets); diff != "" {
		t.Errorf("addSecretToSA() failed:\n%s", diff)
	}
}
