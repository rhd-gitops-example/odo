package pipelines

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	kustomize = "kustomization.yaml"
)

// InitialiseParameters is a struct that provides flags for initialize command
type InitialiseParameters struct {
	GitOpsRepo               string
	Output                   string
	AppGithubSecret          string
	AppGitRepo               string
	AppImageRepo             string
	Prefix                   string
	InternalRegistryHostname string
	DockerConfigJSONFileName string
	SkipChecks               bool
}

// repo represents gitops directory struct
type repo struct {
	path     string
	isDir    bool
	contents []*repo
}

// validate the existing gitops dir
func (r *repo) validate() error {
	if exists, err := isExisting(r.path); !exists {
		return err
	}
	for _, a := range r.contents {
		if err := a.validate(); err != nil {
			return err
		}
	}
	return nil
}

// create the required gitops files and dirs
func (r *repo) create() error {
	if r.isDir {
		if err := createDir(r.path); err != nil {
			return err
		}
	} else {
		if err := createFile(r.path); err != nil {
			return err
		}
	}
	for _, a := range r.contents {
		if err := a.create(); err != nil {
			return err
		}
	}
	return nil
}

// getRepoStruct() returns gitops dir structure
func getRepoStruct(path string) *repo {
	envPath := filepath.Join(path, "env")
	appsPath := filepath.Join(path, "apps")
	basePath := filepath.Join(envPath, "base")
	return &repo{
		path:  path,
		isDir: true,
		contents: []*repo{
			&repo{
				path:  envPath,
				isDir: true,
				contents: []*repo{
					&repo{
						path:  basePath,
						isDir: true,
						contents: []*repo{
							&repo{
								path:  filepath.Join(basePath, kustomize),
								isDir: false,
							},
						},
					},
				},
			},
			&repo{
				path:  appsPath,
				isDir: true,
			},
		},
	}
}

// Initialise function will initialise the gitops directory
func Initialise(o *InitialiseParameters) error {

	if !o.SkipChecks {
		installed, err := checkTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	repo := getRepoStruct(o.Output)

	// check if the dir already exists
	exists, _ := isExisting(o.Output)
	if !exists {
		if err := repo.create(); err != nil {
			return err
		}
	}

	// validate existing dir
	if err := repo.validate(); err != nil {
		return err
	}
	fmt.Println("gitops initialised")
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

// create and invoke a Tekton Checker
func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}
