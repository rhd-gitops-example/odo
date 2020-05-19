package utility

import (
	"errors"
	"net/url"
	"strings"
)

func GetDriverName(rawURL string) (string, error) {

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if s := strings.TrimSuffix(u.Host, ".com"); s != u.Host {
		return strings.ToLower(s), nil
	}

	if s := strings.TrimSuffix(u.Host, ".org"); s != u.Host {
		return strings.ToLower(s), nil
	}

	return "", errors.New("unknown Git server: " + u.Host)
}

func GetRepoName(rawURL string) (string, error) {

	var components []string

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	for _, s := range strings.Split(u.Path, "/") {
		if s != "" {
			components = append(components, s)
		}
	}

	if len(components) != 2 {
		return "", errors.New("failed to get Git repo: " + u.Path)
	}

	components[1] = strings.TrimSuffix(components[1], ".git")

	for _, s := range components {
		if strings.Index(s, ".") != -1 {
			return "", errors.New("failed to get Git repo: " + u.Path)
		}
	}

	return components[0] + "/" + components[1], nil
}
