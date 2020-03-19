package pipelines

import (
	"os"
	"path/filepath"
)

var (
	kustomize = "kustomization.yaml"
)

// folder represents the directory struct
type folder struct {
	path    string
	isDir   bool
	folders []*folder
}

// validate the existing gitops dir
func (f *folder) validate() error {
	if exists, err := isExisting(f.path); !exists {
		return err
	}
	for _, folder := range f.folders {
		if err := folder.validate(); err != nil {
			return err
		}
	}
	return nil
}

// create the required gitops files and dirs
func (f *folder) create() error {
	if f.isDir {
		if err := createDir(f.path); err != nil {
			return err
		}
	} else {
		if err := createFile(f.path); err != nil {
			return err
		}
	}
	for _, dir := range f.folders {
		if err := dir.create(); err != nil {
			return err
		}
	}
	return nil
}

// getGitopsFolder() returns an initialised gitops folder struct
func getGitopsFolder(path, prefix string) *folder {
	envs := filepath.Join(path, "envs")
	apps := filepath.Join(path, "apps")
	base := filepath.Join(envs, "base")
	cicd := filepath.Join(envs, addPrefix(prefix, "cicd-environment"))
	pipeline := filepath.Join(cicd, "pipelines")
	return &folder{
		path:  path,
		isDir: true,
		folders: []*folder{
			&folder{
				path:  envs,
				isDir: true,
				folders: []*folder{
					&folder{
						path:    base,
						isDir:   true,
						folders: []*folder{addKustomize(base)},
					},
					&folder{
						path:  cicd,
						isDir: true,
						folders: []*folder{
							&folder{
								path:  pipeline,
								isDir: true,
								folders: []*folder{
									addKustomize(pipeline),
								},
							},
							addKustomize(cicd),
						},
					},
				},
			},
			&folder{path: apps, isDir: true},
		},
	}
}

func addKustomize(path string) *folder {
	return &folder{
		path:    filepath.Join(path, kustomize),
		isDir:   false,
		folders: []*folder{},
	}
}

func addPrefix(prefix, name string) string {
	if prefix != "" {
		return prefix + name
	}
	return name
}

func createDir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

func createFile(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}
	return nil
}

func isExisting(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
