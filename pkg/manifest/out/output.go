package out

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

		// Add all output items
		AddAll(map[string]interface{})

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

// AddAll adds all items from src to Items map
func (o *BaseOutput) AddAll(src map[string]interface{}) {
	for k, v := range src {
		o.Items[k] = v
	}
}
