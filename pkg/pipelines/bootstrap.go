package pipelines

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"sigs.k8s.io/yaml"
)

var DefaultTokenFilename string = "~/Downloads/token"
var DefaultQuayIOAuthFilename string = "~/Downloads/<username>-auth.json"

// BootstrapOptions is a struct that provides the optional flags
type BootstrapOptions struct {
	Prefix             string
	TokenFilename      string
	QuayIOAuthFilename string
}

// Bootstrap is the main driver for getting OpenShift pipelines for GitOps
// configured with a basic configuration.
func Bootstrap(quayUsername, baseRepo string, o *BootstrapOptions) error {

	// First, check for Tekton.  We proceed only if Tekton is installed
	installed, err := checkTektonInstall()
	if err != nil {
		return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
	}
	if !installed {
		return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
	}

	outputs := make([]interface{}, 0)

	tokenPath, err := getTokenFileName(o)
	if err != nil {
		return fmt.Errorf("failed to generate token path to file: %w", err)
	}
	f, err := os.Open(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to open path to TokenFileName: %w", err)
	}
	defer f.Close()

	githubAuth, err := createOpaqueSecret("github-auth", f)
	if err != nil {
		return err
	}
	outputs = append(outputs, githubAuth)

	authJSONPath, err := getQuayIOAuthFileName(quayUsername, o)
	if err != nil {
		return fmt.Errorf("failed to generate path to file: %w", err)
	}

	f, err = os.Open(authJSONPath)
	if err != nil {
		return fmt.Errorf("failed to open path authJSON : %w", err)
	}
	defer f.Close()

	dockerSecret, err := createDockerConfigSecret("regcred", f)
	if err != nil {
		return err
	}
	outputs = append(outputs, dockerSecret)

	eventListener := eventlisteners.GenerateEventListener(baseRepo)
	outputs = append(outputs, eventListener)

	route := routes.GenerateRoute()
	outputs = append(outputs, route)

	for _, r := range outputs {
		data, err := yaml.Marshal(r)
		if err != nil {
			return err
		}
		fmt.Printf("%s---\n", data)
	}
	return nil
}

func pathToDownloadedFile(fname string) (string, error) {
	return homedir.Expand(path.Join("~/Downloads/", fname))
}

//to get the github token file name
func getTokenFileName(o *BootstrapOptions) (string, error) {
	if o.TokenFilename == DefaultTokenFilename {
		return pathToDownloadedFile("token")
	}
	return homedir.Expand(o.TokenFilename)
}

// to get the quay file name
func getQuayIOAuthFileName(quayUsername string, o *BootstrapOptions) (string, error) {

	if o.QuayIOAuthFilename == DefaultQuayIOAuthFilename {
		return pathToDownloadedFile(quayUsername + "-auth.json")
	}
	return homedir.Expand(o.QuayIOAuthFilename)

}

// create and invoke a Tekton Checker
func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}
