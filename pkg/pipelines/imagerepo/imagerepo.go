package imagerepo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/namespaces"
	"github.com/openshift/odo/pkg/pipelines/roles"

	"github.com/openshift/odo/pkg/pipelines/config"

	res "github.com/openshift/odo/pkg/pipelines/resources"

	corev1 "k8s.io/api/core/v1"
)

// ValidateImageRepo validates the input image repo.  It determines if it is
// for internal registry and prepend internal registry hostname if neccessary.
func ValidateImageRepo(imageRepo, registryURL string) (bool, string, error) {
	components := strings.Split(imageRepo, "/")

	// repo url has minimum of 2 components
	if len(components) < 2 {
		return false, "", imageRepoValidationErrors(imageRepo)
	}

	for _, v := range components {
		// check for empty components
		if strings.TrimSpace(v) == "" {
			return false, "", imageRepoValidationErrors(imageRepo)
		}
		// check for white spaces
		if len(v) > len(strings.TrimSpace(v)) {
			return false, "", imageRepoValidationErrors(imageRepo)
		}
	}

	if len(components) == 2 {
		if components[0] == "docker.io" || components[0] == "quay.io" {
			// we recognize docker.io and quay.io.  It is missing one component
			return false, "", imageRepoValidationErrors(imageRepo)
		}
		// We have format like <project>/<app> which is an internal registry.
		// We prepend the internal registry hostname.
		return true, registryURL + "/" + imageRepo, nil
	}

	// Check the first component to see if it is an internal registry
	if len(components) == 3 {
		return components[0] == registryURL, imageRepo, nil
	}

	// > 3 components.  invalid repo
	return false, "", imageRepoValidationErrors(imageRepo)
}

func imageRepoValidationErrors(imageRepo string) error {
	return fmt.Errorf("failed to parse image repo:%s, expected image repository in the form <registry>/<username>/<repository> or <project>/<app> for internal registry", imageRepo)
}

func CreateInternalRegistryResources(m *config.Manifest, cicdEnv *config.Environment, sa *corev1.ServiceAccount, imageRepo string) (res.Resources, error) {

	// Provide access to service account for using internal registry
	namespace := strings.Split(imageRepo, "/")[1]

	namespaceExists, err := namespaces.ExistsInManifestOrCluster(m, namespace)
	if err != nil {
		return nil, err
	}

	resources := res.Resources{}

	if !namespaceExists {
		namespacePath := filepath.Join(config.PathForEnvironment(cicdEnv), "base", "pipelines", "01-namespaces", fmt.Sprintf("%s.yaml", namespace))
		resources[namespacePath] = namespaces.Create(namespace)
	}

	resources = res.Merge(createInternalRegistryRoleBinding(cicdEnv, namespace, sa), resources)
	return resources, nil
}

func createInternalRegistryRoleBinding(cicdEnv *config.Environment, ns string, sa *corev1.ServiceAccount) res.Resources {
	roleBindingName := fmt.Sprintf("internal-registry-%s-binding", ns)
	roleBindingPath := filepath.Join(config.PathForEnvironment(cicdEnv), "base", "pipelines", "02-rolebindings", fmt.Sprintf("%s.yaml", roleBindingName))
	return res.Resources{roleBindingPath: roles.CreateRoleBinding(meta.NamespacedName(ns, roleBindingName), sa, "ClusterRole", "edit")}
}
