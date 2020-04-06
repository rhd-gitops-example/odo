package out

import (
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

type (
	// BaseOutput base structure
	BaseOutput struct {
		// key: path of the resource
		// value: YAML content of the resource
		Items map[string]interface{}

		Output
	}

	// Output inferface defines output operations
	Output interface {
		// Add an output item
		Add(path string, data interface{})

		// Write out all outpout items
		Write() error

		// Get output item path names
		GetPaths() []string
	}
)

// New creates a new BaseOotput
func New() BaseOutput {
	return BaseOutput{Items: make(map[string]interface{})}
}

// Add - adds an outoupt item
func (o *BaseOutput) Add(path string, data interface{}) {
	o.Items[path] = data
}

// GetPaths returns all output item path names
func (o *BaseOutput) GetPaths() []string {
	paths := make([]string, len(o.Items))

	i := 0
	for path := range o.Items {
		paths[i] = path
		i++
	}

	return paths
}

// Marshal marshals object to given writer
func Marshal(writer io.Writer, object interface{}) error {
	data, err := yaml.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = fmt.Fprintf(writer, "%s", data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	return nil
}
