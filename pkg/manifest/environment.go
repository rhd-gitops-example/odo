package manifest

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/manifest/ioutils"
	"github.com/openshift/odo/pkg/manifest/meta"
	"github.com/openshift/odo/pkg/manifest/out/fs"
	"github.com/openshift/odo/pkg/manifest/roles"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"github.com/spf13/afero"
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

	envPath := getEnvPath(o.EnvName, o.Prefix)

	output, err := fs.New(o.Output, &afero.Afero{Fs: afero.NewOsFs()}, func() error {
		// check if the gitops dir exists
		exists, err := ioutils.IsExisting(o.Output)
		if !exists {
			return err
		}
		// check if the environment dir already exists
		exists, err = ioutils.IsExisting(envPath)
		if exists {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create output : %w", err)
	}

	// base kustomization
	output.Add(filepath.Join(envPath, "base", Kustomize), yaml.Kustomization("resources", []string{envNamespace, envRoleBinding}))

	// overlay kustomization
	output.Add(filepath.Join(envPath, "overlays", Kustomize), yaml.Kustomization("bases", []string{"../base"}))

	envName := AddPrefix(o.Prefix, o.EnvName)
	output.AddAll(envResources(o.Prefix, envPath, envName))

	return output.Write()
}

func envResources(prefix, envPath, envName string) map[string]interface{} {

	namespaces := NamespaceNames(prefix)

	outputs := map[string]interface{}{}
	basePath := filepath.Join(envPath, "base")

	outputs[filepath.Join(basePath, envNamespace)] = CreateNamespace(envName)

	sa := roles.CreateServiceAccount(meta.NamespacedName(namespaces["cicd"], saName))
	outputs[filepath.Join(basePath, envRoleBinding)] = roles.CreateRoleBinding(meta.NamespacedName(envName, roleBindingName), sa, "ClusterRole", "edit")

	return outputs
}

func getEnvPath(envName, prefix string) string {
	return filepath.Join(EnvsDir, AddPrefix(prefix, envName))
}
