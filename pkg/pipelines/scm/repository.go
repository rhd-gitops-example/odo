package scm

func NewRepository(rawURL string) (Repository, error) {
	repoType, err := GetDriverName(rawURL)
	if err != nil {
		return nil, err
	}
	switch repoType {
	case "github":
		return NewGithubRepository(rawURL)
	}
	return nil, nil
}
