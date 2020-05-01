package webhook

import (
	"errors"
	"fmt"

	"github.com/openshift/odo/pkg/pipelines/config"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
)

// Create create a new webhook on the target Git Repository
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

	return nil
}

// Delete deletes webhooks on the target Git Repository that match the listener address
// names is a {envName, appName, seviceName} tuple.
func Delete(accessToken, pipelines string, names []string, isCICD, isInsecure bool) error {

	return nil
}

func getRepoURL(manifest *config.Manifest, isCICD bool, names []string) string {

	if isCICD {
		return manifest.GitOpsURL
	}

	return getSourceRepoURL(manifest, names)

}

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
