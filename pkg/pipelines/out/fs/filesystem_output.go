package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/out"
)

type output struct {
	out.BaseOutput
	rootDir string
}

// New creates an Output that outputs to filesystem
func New(repo, outputFolder string) (out.Output, error) {
	gitopsName := getGitopsRepoName(repo)
	gitopsPath := filepath.Join(outputFolder, gitopsName)

	// check if the gitops dir already exists
	exists, _ := isExisting(gitopsPath)
	if exists {
		return nil, fmt.Errorf("%s already exists at %s", gitopsName, gitopsPath)
	}

	return &output{
		BaseOutput: out.New(),
		rootDir:    gitopsPath,
	}, nil
}

// Write all outpout items to filesystem
func (o *output) Write() error {
	for filename, item := range o.Items {
		err := writeToFile(filepath.Join(o.rootDir, filename), item)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeToFile(filename string, item interface{}) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return out.Marshal(f, item)
}

func getGitopsRepoName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func isExisting(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
