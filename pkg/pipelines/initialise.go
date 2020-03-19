package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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

	gitopsPath := filepath.Join(o.Output, getGitopsRepoName(o.GitOpsRepo))
	gitopsFolder := getGitopsFolder(gitopsPath, o.Prefix)

	// check if the gitops dir already exists
	exists, _ := isExisting(gitopsPath)
	if !exists {
		if err := gitopsFolder.create(); err != nil {
			return err
		}
	}

	// validate existing gitops dir
	if err := gitopsFolder.validate(); err != nil {
		return err
	}

	// create the required resources here

	fmt.Println("Gitops initialised")
	return nil
}

// create and invoke a Tekton Checker
func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}

func getGitopsRepoName(repo string) string {
	return strings.Split(repo, "/")[1]
}
