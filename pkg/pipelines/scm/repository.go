package scm

var (
	supportedTriggerBindings []string
)

// NewRepository returns a suitable Repository instance
// based on the driver name (github,gitlab,etc)
func NewRepository(rawURL string) (Repository, error) {
	repoType, err := getDriverName(rawURL)
	if err != nil {
		return nil, err
	}
	switch repoType {
	case "github":
		return NewGitHubRepository(rawURL)
	case "gitlab":
		return NewGitLabRepository(rawURL)
	}
	return nil, invalidRepoTypeError(rawURL)
}

// IsRepositoryBinding checks if a binding
// is a supported binding type
func IsRepositoryBinding(b string) bool {
	for _, binding := range supportedTriggerBindings {
		if binding == b {
			return true
		}
	}
	return false
}
