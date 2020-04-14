package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/spf13/afero"
)

func Build(fs afero.Fs, m *config.Manifest) (map[string]interface{}, error) {
	files := make(map[string]interface{})
	eb := &envBuilder{fs: fs, files: files}
	err := m.Walk(eb)
	return eb.files, err
}

type envBuilder struct {
	files map[string]interface{}
	fs    afero.Fs
}

type Kustomization struct {
	Resources []string `yaml:"resources"`
	Bases     []string `yaml:"bases"`
}

func listFiles(fs afero.Fs, base string) ([]string, error) {
	files := []string{}
	err := afero.Walk(fs, base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, strings.TrimPrefix(path, base+"/"))
		return nil
	})
	return files, err
}

func (b *envBuilder) Environment(env *config.Environment) error {
	basePath := filepath.Join(config.PathForEnvironment(env), "base")
	envFiles := map[string]interface{}{}
	filename := filepath.Join(basePath, fmt.Sprintf("%s-environment.yaml", env.Name))
	envFiles[filename] = CreateNamespace(env.Name)
	kustomizedFilenames, err := listFiles(b.fs, basePath)
	if err != nil {
		return fmt.Errorf("failed to list initial files: %s", err)
	}

	for k, _ := range envFiles {
		kustomizedFilenames = append(kustomizedFilenames, filepath.Base(k))
	}
	envFiles[filepath.Join(basePath, "kustomization.yaml")] = &Kustomization{Resources: kustomizedFilenames}

	overlaysPath := filepath.Join(config.PathForEnvironment(env), "overlays")
	relPath, err := filepath.Rel(overlaysPath, basePath)
	if err != nil {
		return err
	}
	envFiles[filepath.Join(overlaysPath, "kustomization.yaml")] = &Kustomization{Bases: []string{relPath}}
	b.files = merge(envFiles, b.files)
	return nil
}
