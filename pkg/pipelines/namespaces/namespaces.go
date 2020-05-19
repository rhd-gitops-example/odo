package namespaces

import (
	"fmt"

	"github.com/openshift/odo/pkg/pipelines/clientconfig"
	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	namespaceBaseNames = map[string]string{
		"dev":   "dev",
		"stage": "stage",
		"cicd":  "cicd",
	}

	namespaceTypeMeta = meta.TypeMeta("Namespace", "v1")
)

// Namespaces create namespaces for the given names.
func Namespaces(names []string) []*corev1.Namespace {
	ns := []*corev1.Namespace{}
	for _, n := range names {
		ns = append(ns, Create(n))
	}
	return ns
}

// NamesWithPrefix returns namespaces of all environments based on the prefix,
// and using the set of predefined names: dev, stage, cicd.
func NamesWithPrefix(prefix string) map[string]string {
	prefixedNames := make(map[string]string)
	for k, v := range namespaceBaseNames {
		prefixedNames[k] = fmt.Sprintf("%s%s", prefix, v)
	}
	return prefixedNames
}

// Create creates a Namespace value from a string.
func Create(name string) *corev1.Namespace {
	ns := &corev1.Namespace{
		TypeMeta: namespaceTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return ns
}

// GetClientSet creates and returns a new Kubernetes clientset.
func GetClientSet() (*kubernetes.Clientset, error) {
	clientConfig, err := clientconfig.GetRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get client config due to %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get APIs client due to %w", err)
	}
	return clientSet, nil
}

// Exists returns true if the given namespace exists
func Exists(clientSet kubernetes.Interface, name string) (bool, error) {
	_, err := clientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ExistsInManifestOrCluster returns true if the namespace exists in either in the manfiest or cluster
// This function looks in the manifest first.  If manifest is nil, it will only check in the Cluster
func ExistsInManifestOrCluster(m *config.Manifest, name string) (bool, error) {
	if m != nil {
		for _, env := range m.Environments {
			if env.Name == name {
				return true, nil
			}
		}
	}

	clientSet, err := GetClientSet()
	if err != nil {
		return false, err
	}

	return Exists(clientSet, name)
}
