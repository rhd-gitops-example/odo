package pipelines

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

var baseNames = map[string]string{
	"dev":   "dev-environment",
	"stage": "stage-environment",
	"cicd":  "cicd-environment",
}

func createNamespaces(prefix string) []*corev1.Namespace {
	ns := []*corev1.Namespace{}
	for _, n := range namespaceNames(prefix) {
		ns = append(ns, createNamespace(n))
	}
	return ns
}

func namespaceNames(prefix string) map[string]string {
	prefixedNames := make(map[string]string)
	for k, v := range baseNames {
		prefixedNames[k] = fmt.Sprintf("%s%s", prefix, v)
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
