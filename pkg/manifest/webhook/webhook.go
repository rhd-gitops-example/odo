package webhook

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

// AddWebHook add web hook
func AddWebHook(gitRepoURL, token, listenerURL, secret string) error {

	parsedURL, err := url.Parse(gitRepoURL)
	if err != nil {
		return err
	}

	driverName, err := getDriverName(parsedURL)
	if err != nil {
		return err
	}

	repo, err := getRepo(parsedURL)
	if err != nil {
		return err
	}

	client, err := factory.NewClient(driverName, "", token)
	if err != nil {
		return err
	}

	in := &scm.HookInput{
		Target: listenerURL,
		Secret: secret,
		Events: scm.HookEvents{
			PullRequest: true,
			Push:        true,
		},
	}

	_, _, err = client.Repositories.CreateHook(context.Background(), repo, in)

	return err
}

func getDriverName(u *url.URL) (string, error) {
	if s := strings.TrimSuffix(u.Host, ".com"); s != u.Host {
		return strings.ToLower(s), nil
	}

	if s := strings.TrimSuffix(u.Host, ".org"); s != u.Host {
		return strings.ToLower(s), nil
	}

	return "", errors.New("unknown Git server: " + u.Host)
}

func getRepo(u *url.URL) (string, error) {

	var components []string

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
