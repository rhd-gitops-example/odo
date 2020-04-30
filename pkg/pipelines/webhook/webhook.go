package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

type webhook struct {
	*scm.Client
	repoName string
}

// New creates a new webhook object
func New(gitRepoURL, token string) (*webhook, error) {
	parsedURL, err := url.Parse(gitRepoURL)
	if err != nil {
		return nil, err
	}

	driverName, err := getDriverName(parsedURL)
	if err != nil {
		return nil, err
	}

	client, err := factory.NewClient(driverName, "", token)
	if err != nil {
		return nil, err
	}

	repoName, err := getRepoName(parsedURL)
	if err != nil {
		return nil, err
	}

	return &webhook{repoName: repoName, Client: client}, nil
}

func (w *webhook) list(listenerURL string) ([]string, error) {
	hooks, _, err := w.Client.Repositories.ListHooks(context.Background(), w.repoName, scm.ListOptions{})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, hook := range hooks {
		if hook.Target == listenerURL {
			ids = append(ids, hook.ID)
		}
	}

	return ids, nil
}

func (w *webhook) delete(listenerURL string, ids []string) error {
	for _, id := range ids {
		_, err := w.Client.Repositories.DeleteHook(context.Background(), w.repoName, id)
		if err != nil {
			return fmt.Errorf("failed to delete webhook id %s: %w", id, err)
		}
	}
	return nil
}

func (w *webhook) add(listenerURL, secret string) error {
	in := &scm.HookInput{
		Target: listenerURL,
		Secret: secret,
		Events: scm.HookEvents{
			PullRequest: true,
			Push:        true,
		},
	}

	_, _, err := w.Client.Repositories.CreateHook(context.Background(), w.repoName, in)
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

func getRepoName(u *url.URL) (string, error) {

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
