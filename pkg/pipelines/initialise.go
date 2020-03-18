package pipelines

import (
	"errors"
	"fmt"
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

	// check if the dir already exists
	exists, _ := isExisting(o.Output)
	if !exists {
		if err := createRepositoryLayout(o.Output); err != nil {
			return err
		}
	}

	// validate existing dir
	if err := validateRepositoryLayout(o.Output); err != nil {
		return err
	}
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
