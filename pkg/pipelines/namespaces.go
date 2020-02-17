package pipelines

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

var baseNames = []string{
	"dev-environment",
	"stage-environment",
	"cicd-environment",
}

func createNamespaces(prefix string) []*corev1.Namespace {
	ns := []*corev1.Namespace{}
	for _, n := range namespaceNames(prefix) {
		ns = append(ns, createNamespace(n))
	}
	return ns
}

func namespaceNames(prefix string) []string {
	prefixedNames := make([]string, len(baseNames))
	for i, n := range baseNames {
		prefixedNames[i] = fmt.Sprintf("%s%s", prefix, n)
	}
	return prefixedNames
}

func createNamespace(name string) *corev1.Namespace {
	ns := &corev1.Namespace{
		TypeMeta:   createTypeMeta("Namespace", "v1"),
		ObjectMeta: createObjectMeta(name),
	}
	return ns
}
