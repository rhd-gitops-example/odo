package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/odo/pkg/manifest/out"
	"github.com/openshift/odo/pkg/manifest/yaml"
)

type output struct {
	out.BaseOutput
	rootDir string
}

// New creates an Output that outputs to filesystem
func New(outputFolder string, validation func() error) (out.Output, error) {
	if err := validation(); err != nil {
		return nil, err
	}

	return &output{
		BaseOutput: out.New(),
		rootDir:    outputFolder,
	}, nil
}

// Write all output items to filesystem.
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
	return yaml.Marshal(f, item)
}
