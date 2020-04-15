package fs

import (
	"fmt"
	"path/filepath"

	"github.com/openshift/odo/pkg/manifest/out"
	"github.com/openshift/odo/pkg/manifest/yaml"
	"github.com/spf13/afero"
)

type output struct {
	out.BaseOutput
	rootDir string
	filesys *afero.Afero
}

// New creates an Output that outputs to filesystem
func New(outputFolder string, aferofs *afero.Afero, validation func() error) (out.Output, error) {
	if err := validation(); err != nil {
		return nil, err
	}

	return &output{
		BaseOutput: out.New(),
		rootDir:    outputFolder,
		filesys:    aferofs,
	}, nil
}

// Write all outpout items to filesystem
func (o *output) Write() error {
	for filename, item := range o.Items {
		err := o.writeToFile(filepath.Join(o.rootDir, filename), item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *output) writeToFile(filename string, item interface{}) error {
	err := o.filesys.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := o.filesys.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return yaml.Marshal(f, item)
}
