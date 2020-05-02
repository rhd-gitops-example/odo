package webhook

import (
	"errors"
	"fmt"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/git"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/secrets"
)

type webhook struct {
	clusterResource *resources
	repository      *git.Repository
	gitRepoURL      string
	cicdNamepace    string
	listenerURL     string
	accessToken     string
	names           []string
	isCICD          bool
}

func newWebhookInfo(accessToken, pipelinesFile string, names []string, isCICD, isInsecure bool) (*webhook, error) {
	manifest, err := config.ParseFile(ioutils.NewFilesystem(), pipelinesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pipelines: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	gitRepoURL := getRepoURL(manifest, isCICD, names)
	if gitRepoURL == "" {
		return nil, errors.New("failed to find Git repostory URL in manifest")
	}

	cicdNamepace := getCICDNamespace(manifest)
	if cicdNamepace == "" {
		return nil, errors.New("failed to find CICD namespace in manifest")
	}

	clusterResources, err := newResources()
	if err != nil {
		return nil, err
	}

	repository, err := git.NewRepository(gitRepoURL, accessToken)
	if err != nil {
		return nil, err
	}

	listenerURL, err := getListenerURL(clusterResources, cicdNamepace, isInsecure)
	if err != nil {
		return nil, fmt.Errorf("failed to get event listener URL: %w", err)
	}

	return &webhook{clusterResources, repository, gitRepoURL, cicdNamepace, listenerURL, accessToken, names, isCICD}, nil
}

// Create creates a new webhook on the target Git Repository
// names is a {envName, appName, seviceName} tuple.
func Create(accessToken, pipelinesFile string, names []string, isCICD, isInsecure bool) error {
	webhook, err := newWebhookInfo(accessToken, pipelinesFile, names, isCICD, isInsecure)
	if err != nil {
		return err
	}

	exists, err := webhook.exists()
	if err != nil {
		return err
	}

	if exists {
		return errors.New("webhook already exists")
	}

	return webhook.create()
}

// Delete deletes webhooks on the target Git Repository that match the listener address
// names is a {envName, appName, seviceName} tuple.
func Delete(accessToken, pipelinesFile string, names []string, isCICD, isInsecure bool) error {

	webhook, err := newWebhookInfo(accessToken, pipelinesFile, names, isCICD, isInsecure)
	if err != nil {
		return err
	}

	ids, err := webhook.list()
	if err != nil {
		return err
	}

	return webhook.delete(ids)

}

// List returns an array of webhook IDs for the target Git repository/listeners
func List(accessToken, pipelinesFile string, names []string, isCICD, isInsecure bool) ([]string, error) {

	webhook, err := newWebhookInfo(accessToken, pipelinesFile, names, isCICD, isInsecure)
	if err != nil {
		return nil, err
	}

	return webhook.list()

}
func (w *webhook) exists() (bool, error) {
	ids, err := w.repository.ListWebhooks(w.listenerURL)
	if err != nil {
		return false, err
	}
	return len(ids) > 0, nil
}

func (w *webhook) list() ([]string, error) {
	return w.repository.ListWebhooks(w.listenerURL)
}

func (w *webhook) delete(ids []string) error {
	return w.repository.DeleteWebhooks(w.listenerURL, ids)
}

func (w *webhook) create() error {
	secret, err := getWebhookSecret(w.clusterResource, w.cicdNamepace, w.isCICD, w.names)
	if err != nil {
		return fmt.Errorf("failed to get webhook secret: %w", err)
	}

	return w.repository.CreateWehoook(w.listenerURL, secret)
}

// Get Git repository URL whether it is CICD configuration or service source repository
// Return "" if not found
func getRepoURL(manifest *config.Manifest, isCICD bool, names []string) string {
	if isCICD {
		return manifest.GitOpsURL
	}
	return getSourceRepoURL(manifest, names)
}

// Get serrice source repository URL.  Return "" if not found
func getSourceRepoURL(manifest *config.Manifest, names []string) string {
	for _, env := range manifest.Environments {
		if env.Name == names[0] {
			for _, app := range env.Apps {
				if app.Name == names[1] {
					for _, svc := range app.Services {
						if svc.Name == names[2] {
							return svc.SourceURL
						}
					}
					break
				}
			}
			break
		}
	}
	return ""
}

// get CICD namespace.  Return "" if not found
func getCICDNamespace(manifest *config.Manifest) string {
	for _, env := range manifest.Environments {
		if env.IsCICD {
			return env.Name
		}
	}
	return ""
}

func getListenerURL(r *resources, cicdNamespace string, isInsecure bool) (string, error) {
	host, port, err := r.getListenerAddress(cicdNamespace, routes.GitOpsWebhookEventListenerRouteName)
	if err != nil {
		return "", err
	}
	return buildURL(host, port, isInsecure), nil
}

func buildURL(host, port string, isInsecure bool) string {
	scheme := "http"
	if !isInsecure {
		scheme = scheme + "s"
	}
	url := scheme + "://" + host
	if isInsecure {
		if port != "80" {
			url = url + ":" + port
		}
	} else {
		if port != "443" {
			url = url + ":" + port
		}
	}
	return url
}

func getWebhookSecret(r *resources, namespace string, isCICD bool, names []string) (string, error) {
	var secretName string
	if isCICD {
		secretName = eventlisteners.GitOpsWebhookSecret
	} else {
		// currently, use the app name to create webhook secret name.
		// also currently, service webhook secret are in CICI namespace
		secretName = secrets.MakeSerivceWebhookSecretName(names[1])
	}
	return r.getWebhookSecret(namespace, secretName, eventlisteners.WebhookSecretKey)
}