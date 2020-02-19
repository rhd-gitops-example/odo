package pipelines

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/tasks"
	"sigs.k8s.io/yaml"
)

// BootstrapOptions is a struct that provides the optional flags
type BootstrapOptions struct {
	Prefix string
}

// Bootstrap is the main driver for getting OpenShift pipelines for GitOps
// configured with a basic configuration.
func Bootstrap(quayUsername, baseRepo, githubToken, quayIOAuthFilename string, o *BootstrapOptions) error {

	// First, check for Tekton.  We proceed only if Tekton is installed
	installed, err := checkTektonInstall()
	if err != nil {
		return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
	}
	if !installed {
		return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
	}

	outputs := make([]interface{}, 0)

	githubAuth, err := createOpaqueSecret("github-auth", githubToken)
	if err != nil {
		return err
	}
	outputs = append(outputs, githubAuth)

	authJSONPath, err := homedir.Expand(quayIOAuthFilename)
	if err != nil {
		return fmt.Errorf("failed to expand the path %s: %w", authJSONPath, err)
	}

	f, err := os.Open(authJSONPath)
	if err != nil {
		return fmt.Errorf("failed to expand the path %s: %w", authJSONPath, err)
	}
	defer f.Close()

	dockerSecret, err := createDockerConfigSecret("regcred", f)
	if err != nil {
		return err
	}
	outputs = append(outputs, dockerSecret)

	tasks := tasks.Generate(githubAuth.GetName())
	for _, task := range tasks {
		outputs = append(outputs, task)
	}

	eventListener := eventlisteners.Generate(baseRepo)
	outputs = append(outputs, eventListener)

	route := routes.Generate()
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

// create and invoke a Tekton Checker
func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}
