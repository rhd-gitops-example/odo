package pipelines

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"sigs.k8s.io/yaml"
)

// Bootstrap is the main driver for getting OpenShift pipelines for GitOps
// configured with a basic configuration.
func Bootstrap(quayUsername, baseRepo, prefix string) error {
	outputs := make([]interface{}, 0)

	tokenPath, err := pathToDownloadedFile("token")
	if err != nil {
		return fmt.Errorf("failed to generate path to file: %w", err)
	}
	f, err := os.Open(tokenPath)
	if err != nil {
		return err
	}
	defer f.Close()

	githubAuth, err := CreateOpaqueSecret("github-auth", f)
	if err != nil {
		return err
	}
	outputs = append(outputs, githubAuth)

	authJsonPath, err := pathToDownloadedFile(quayUsername + "-auth.json")
	if err != nil {
		return fmt.Errorf("failed to generate path to file: %w", err)
	}

	f, err = os.Open(authJsonPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dockerSecret, err := CreateDockerConfigSecret("regcred", f)
	if err != nil {
		return err
	}
	outputs = append(outputs, dockerSecret)

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
