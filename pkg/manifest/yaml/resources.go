package yaml

import (
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

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

// Kustomization creates kustomization yaml content
// File content is a map of the given name and values
func Kustomization(name string, values []string) []interface{} {
	return []interface{}{map[string]interface{}{name: values}}
}
