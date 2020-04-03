package pipelines

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/openshift/odo/pkg/pipelines/roles"
)

const (
	envNamespace   = "namespace.yaml"
	envRoleBinding = "rolebinding.yaml"
)

// EnvParameters encapsulates parameters for add env command
type EnvParameters struct {
	EnvName    string
	Output     string
	Prefix     string
	GitOpsRepo string
}

// Env will bootstrap a new environment directory
func Env(o *EnvParameters) error {

	namespaces := namespaceNames(o.Prefix)

	gitopsName := getGitopsRepoName(o.GitOpsRepo)
	gitopsPath := filepath.Join(o.Output, gitopsName)

	// check if the gitops dir exists
	exists, _ := isExisting(gitopsPath)
	if !exists {
		return fmt.Errorf("%s doesn't exist at %s", gitopsName, o.Output)
	}

	envName := addPrefix(o.Prefix, o.EnvName)
	envPath := getEnvPath(gitopsPath, o.EnvName, o.Prefix)

	err := addKustomize("bases", []string{envNamespace, envRoleBinding}, filepath.Join(envPath, "base", kustomize))
	if err != nil {
		return err
	}

	err = addKustomize("bases", []string{"../base"}, filepath.Join(envPath, "overlays", kustomize))
	if err != nil {
		return err
	}

	outputs := map[string]interface{}{}
	basePath := filepath.Join(envPath, "base")

	outputs[filepath.Join(basePath, envNamespace)] = createNamespace(envName)

	sa := roles.CreateServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))

	outputs[filepath.Join(basePath, envRoleBinding)] = roles.CreateRoleBinding(meta.NamespacedName(envName, roleBindingName), sa, "ClusterRole", roles.ClusterRoleName)

	for path, value := range outputs {
		err := marshalItemsToFile(path, list(value))
		if err != nil {
			return err
		}
	}

	return nil
}

func addEnvBase() {

}

func getEnvPath(gitopsPath, envName, prefix string) string {
	return filepath.Join(gitopsPath, envsDir, addPrefix(prefix, envName))
}
