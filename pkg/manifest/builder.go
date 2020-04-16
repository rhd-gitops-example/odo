package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openshift/odo/pkg/manifest/config"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"github.com/spf13/afero"
)

func buildEnvironments(fs afero.Fs, m *config.Manifest) (res.Resources, error) {
	files := make(res.Resources)
	eb := &envBuilder{fs: fs, files: files}
	err := m.Walk(eb)
	return eb.files, err
}

type envBuilder struct {
	files res.Resources
	fs    afero.Fs
}

func listFiles(fs afero.Fs, base string) (map[string]bool, error) {
	files := map[string]bool{}
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
		files[strings.TrimPrefix(path, base+"/")] = true
		return nil
	})
	return files, err
}

func (b *envBuilder) Environment(env *config.Environment) error {
	if env.IsCICD || env.IsArgoCD {
		return nil
	}
	basePath := filepath.Join(config.PathForEnvironment(env), "base")
	envFiles := filesForEnvironment(basePath, env)
	kustomizedFilenames, err := listFiles(b.fs, basePath)
	if err != nil {
		return fmt.Errorf("failed to list initial files for %s: %s", basePath, err)
	}
	for k, _ := range envFiles {
		kustomizedFilenames[filepath.Base(k)] = true
	}
	envFiles[filepath.Join(basePath, "kustomization.yaml")] = &res.Kustomization{Resources: extractFilenames(kustomizedFilenames)}

	overlaysPath := filepath.Join(config.PathForEnvironment(env), "overlays")
	relPath, err := filepath.Rel(overlaysPath, basePath)
	if err != nil {
		return err
	}
	envFiles[filepath.Join(overlaysPath, "kustomization.yaml")] = &res.Kustomization{Bases: []string{relPath}}
	b.files = res.Merge(envFiles, b.files)
	return nil
}

func filesForEnvironment(basePath string, env *config.Environment) res.Resources {
	envFiles := res.Resources{}
	filename := filepath.Join(basePath, fmt.Sprintf("%s-environment.yaml", env.Name))
	envFiles[filename] = CreateNamespace(env.Name)
	return envFiles
}

func extractFilenames(f map[string]bool) []string {
	names := []string{}
	for k, _ := range f {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
