package pipelines

import (
	"os"
	"path/filepath"
)

func createRepositoryLayout(path string) error {
	envPath := filepath.Join(path, "env")
	appsPath := filepath.Join(path, "apps")
	basePath := filepath.Join(envPath, "base")
	if err := createDir(path); err != nil {
		return err
	}
	if err := createDir(envPath); err != nil {
		return err
	}
	if err := createDir(appsPath); err != nil {
		return err
	}
	if err := createDir(basePath); err != nil {
		return err
	}
	if err := createFile(filepath.Join(basePath, kustomize)); err != nil {
		return err
	}
	return nil
}

// check if the existing dir complies to the gitops repo layout
func validateRepositoryLayout(path string) error {

	envPath := filepath.Join(path, "env")
	appsPath := filepath.Join(path, "apps")
	basePath := filepath.Join(envPath, "base")
	if exists, err := isExisting(path); !exists {
		return err
	}
	if exists, err := isExisting(envPath); !exists {
		return err
	}
	if exists, err := isExisting(appsPath); !exists {
		return err
	}
	if exists, err := isExisting(filepath.Join(envPath, "base")); !exists {
		return err
	}
	if exists, err := isExisting(filepath.Join(basePath, kustomize)); !exists {
		return err
	}
	return nil
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
