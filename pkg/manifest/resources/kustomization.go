package resources

// Kustomization is a structural representation of the Kustomize file format.
type Kustomization struct {
	Resources []string `yaml:"resources,omitempty"`
	Bases     []string `yaml:"bases,omitempty"`
}
