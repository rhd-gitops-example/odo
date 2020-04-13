package pipelines

import (
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
	EnvName string
	Output  string
	Prefix  string
}

// Env will bootstrap a new environment directory
func Env(o *EnvParameters) error {

	envName := addPrefix(o.Prefix, o.EnvName)
	envPath := getEnvPath(o.Output, o.EnvName, o.Prefix)

	// check if the gitops dir exists
	exists, err := isExisting(o.Output)
	if !exists {
		return err
	}

	// check if the environment dir already exists
	exists, err = isExisting(envPath)
	if exists {
		return err
	}

	err = addKustomize("resources", []string{envNamespace, envRoleBinding}, filepath.Join(envPath, "base", kustomize))
	if err != nil {
		return err
	}

	err = addKustomize("bases", []string{"../base"}, filepath.Join(envPath, "overlays", kustomize))
	if err != nil {
		return err
	}

	if err = addEnvResources(o.Prefix, envPath, envName); err != nil {
		return err
	}

	return nil
}

func addEnvResources(prefix, envPath, envName string) error {

	namespaces := namespaceNames(prefix)

	outputs := map[string]interface{}{}
	basePath := filepath.Join(envPath, "base")

	outputs[envNamespace] = CreateNamespace(envName)

	sa := roles.CreateServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))

	outputs[envRoleBinding] = roles.CreateRoleBinding(meta.NamespacedName(envName, roleBindingName), sa, "ClusterRole", "edit")
	_, err := writeResources(basePath, outputs)
	if err != nil {
		return err
	}
	return nil
}

func getEnvPath(gitopsPath, envName, prefix string) string {
	return filepath.Join(gitopsPath, envsDir, addPrefix(prefix, envName))
}
