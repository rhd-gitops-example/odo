package pipelines

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/openshift/odo/pkg/manifest"
	"github.com/openshift/odo/pkg/manifest/ioutils"
	"github.com/openshift/odo/pkg/manifest/out/fs"
	pl "github.com/openshift/odo/pkg/manifest/pipelines"

	"github.com/openshift/odo/pkg/manifest/yaml"
)

// InitParameters is a struct that provides flags for initialise command
type InitParameters struct {
	GitOpsRepo          string
	GitOpsWebhookSecret string
	Output              string
	Prefix              string
	SkipChecks          bool
}

// Init function will initialise the gitops directory
func Init(o *InitParameters) error {

	if !o.SkipChecks {
		installed, err := pl.CheckTektonInstall()
		if err != nil {
			return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
		}
		if !installed {
			return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
		}
	}

	output, err := fs.New(o.Output, func() error {
		// check if the gitops dir already exists
		exists, err := ioutils.IsExisting(o.Output)
		if exists {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create output : %w", err)
	}

	files, err := manifest.CreateResources(o.Prefix, o.GitOpsRepo, o.GitOpsWebhookSecret)
	if err != nil {
		return err
	}

	output.AddAll(files)

	pipelinesPath := manifest.GetPipelinesDir(o.Output, o.Prefix)

	fileNames := output.GetPaths()
	sort.Strings(fileNames)

	// kustomize file should refer all the pipeline resources
	output.Add(filepath.Join(pipelinesPath, manifest.Kustomize), yaml.Kustomization("resources", fileNames))

	// bases
	output.Add(filepath.Join(getCICDDir(o.Output, o.Prefix), manifest.BaseDir, manifest.Kustomize), yaml.Kustomization("bases", []string{"./pipelines"}))

	// Add overlays
	output.Add(filepath.Join(getCICDDir(o.Output, o.Prefix), "overlays", manifest.Kustomize), yaml.Kustomization("bases", []string{"../base"}))

	return output.Write()
}

func getCICDDir(path, prefix string) string {
	return filepath.Join(path, manifest.EnvsDir, manifest.AddPrefix(prefix, manifest.CICDDir))
}
