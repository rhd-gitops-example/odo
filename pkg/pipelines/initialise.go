package pipelines

import (
	"errors"
	"fmt"
	"os"
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

	exists, err := isExistingDir(o.Output)
	if err != nil {
		return err
	}
	if !exists {
		// Clone GitOps repo
		cloneRepository(o.GitOpsRepo, o.Output)

	}
	fmt.Println("gitops initialised")
	return nil
}

func cloneRepository(repo string, output string) error {

	// _, err := git.PlainClone(output, false, &git.CloneOptions{
	// 	URL:      repo,
	// 	Progress: os.Stdout,
	// })
	return nil
}

func isExistingDir(path string) (bool, error) {
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
