package scm

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

// GetAllBindings returns a list of supported binding types
func GetAllBindings() []string {
	return []string{
		githubPRBindingName,
		githubPushBindingName,
		gitlabPRBindingName,
		gitlabPushBindingName,
	}
}
