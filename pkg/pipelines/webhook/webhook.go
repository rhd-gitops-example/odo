package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
)

// Create create a new webhook on the target Git Repository
func Create(accessToken, pipelines, appName, serviceName string, isCICD, isInsecure bool) error {

	manifest, err := config.ParseFile(ioutils.NewFilesystem(), pipelines)
	if err != nil {
		return fmt.Errorf("failed to parse pipelines: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return err
	}

	_, nil := getRepoURL(manifest, isCICD, appName, serviceName)
	return nil
}

// Delete deletes webhooks on the target Git Repository that match the listener address
func Delete(accessToken, pipelines, appName, serviceName string, isCICD, isInsecure bool) error {

	return nil
}

func getRepoURL(manifest *config.Manifest, isCICD bool, appName, serviceName string) (string, error) {
	var repoURL string
	if isCICD {
		repoURL = manifest.GitOpsURL
	} else {
		repoURL = getSourceRepoURL(manifest, serviceName)
	}
	return repoURL, nil
}

func getSourceRepoURL(manifest *config.Manifest, serviceName string) string {
	//for manifest.Environments

	return ""
}
