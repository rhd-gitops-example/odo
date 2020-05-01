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

// Create creates a new webhook on the target Git Repository
// names is a {envName, appName, seviceName} tuple.
func Create(accessToken, pipelines string, names []string, isCICD, isInsecure bool) error {

	manifest, err := config.ParseFile(ioutils.NewFilesystem(), pipelines)
	if err != nil {
		return fmt.Errorf("failed to parse pipelines: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return err
	}

	repoURL := getRepoURL(manifest, isCICD, names)
	if repoURL == "" {
		return errors.New("failed to find Git repostory URL in manifest")
	}

	resources, err := newResources()
	if err := manifest.Validate(); err != nil {
		return err
	}

	cicdNamepace := getCICDNamespace(manifest)
	if cicdNamepace == "" {
		return errors.New("failed to find CICD namespace in manifest")
	}

	listenerURL, err := getListenerURL(resources, cicdNamepace, isInsecure)
	if err != nil {
		return fmt.Errorf("failed to get event listener URL: %w", err)
	}

	secret, err := getWebhookSecret(resources, cicdNamepace, isCICD, names)
	if err != nil {
		return fmt.Errorf("failed to get webhook secret: %w", err)
	}

	repo, err := git.NewRepository(repoURL, accessToken)
	if err != nil {
		return err
	}

	return repo.CreateWehoook(listenerURL, secret)
}

// Delete deletes webhooks on the target Git Repository that match the listener address
// names is a {envName, appName, seviceName} tuple.
func Delete(accessToken, pipelines string, names []string, isCICD, isInsecure bool) error {

	return nil
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
