package pipelines

import (
	"errors"
	"fmt"
)

// InitParameters is a struct that provides flags for initialise command
type AddParameters struct {
	GitopsRepo          string
	GitopsWebhookSecret string
	Output              string
	Prefix              string
	AppGitRepo          string
	AppWebhookSecret    string
	AppImageRepo        string
	EnvName             string
	DockerCfgJson       string
	SkipChecks          bool
}

// Init function will initialise the gitops directory
func Add(o *AddParameters) error {
	if !o.SkipChecks {
		installed, err := checkTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}
	return nil
}
