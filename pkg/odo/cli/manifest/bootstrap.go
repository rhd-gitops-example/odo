package manifest

import (
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/manifest"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// BootstrapRecommendedCommandName the recommended command name
	BootstrapRecommendedCommandName = "bootstrap"
)

var (
	bootstrapExample = ktemplates.Examples(`
    # Bootstrap OpenShift manifest.
    %[1]s 
    `)

	bootstrapLongDesc  = ktemplates.LongDesc(`Bootstrap GitOps CI/CD Manifest`)
	bootstrapShortDesc = `Bootstrap manifest with a starter configuration`
)

// BootstrapParameters encapsulates the parameters for the odo pipelines init command.
type BootstrapParameters struct {
	gitOpsRepo               string // repo to store Gitops resources e.g. org/repo
	gitOpsWebhookSecret      string // used to create Github's shared webhook secret for gitops repo
	output                   string // path to add Gitops resources
	prefix                   string // used to generate the environments in a shared cluster
	skipChecks               bool   // skip Tekton installation checks
	imageRepo                string
	internalRegistryHostname string
	deploymentPath           string
	dockercfgjson            string
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewBootstrapParameters bootstraps a BootstrapParameters instance.
func NewBootstrapParameters() *BootstrapParameters {
	return &BootstrapParameters{}
}

// Complete completes BootstrapParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *BootstrapParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if io.prefix != "" && !strings.HasSuffix(io.prefix, "-") {
		io.prefix = io.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the BootstrapParameters.
func (io *BootstrapParameters) Validate() error {
	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(strings.Split(io.gitOpsRepo, "/")) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", io.gitOpsRepo)
	}
	return nil
}

// Run runs the project bootstrap command.
func (io *BootstrapParameters) Run() error {
	options := manifest.BootstrapParameters{
		GitOpsWebhookSecret:      io.gitOpsWebhookSecret,
		GitOpsRepo:               io.gitOpsRepo,
		Output:                   io.output,
		Prefix:                   io.prefix,
		SkipChecks:               io.skipChecks,
		DeploymentPath:           io.deploymentPath,
		ImageRepo:                io.imageRepo,
		InternalRegistryHostname: io.internalRegistryHostname,
		DockerConfigJSONFilename: io.dockercfgjson,
	}

	return manifest.Bootstrap(&options)
}

// NewCmdBootstrap creates the project init command.
func NewCmdBootstrap(name, fullName string) *cobra.Command {
	o := NewBootstrapParameters()

	initCmd := &cobra.Command{
		Use:     name,
		Short:   initShortDesc,
		Long:    initLongDesc,
		Example: fmt.Sprintf(initExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	initCmd.Flags().StringVar(&o.gitOpsRepo, "gitops-repo", "", "CI/CD pipelines configuration Git repository in this form <username>/<repository>")
	initCmd.Flags().StringVar(&o.gitOpsWebhookSecret, "gitops-webhook-secret", "", "provide the GitHub webhook secret for gitops repository")
	initCmd.Flags().StringVar(&o.dockercfgjson, "dockercfgjson", "", "provide the dockercfg json path")
	initCmd.Flags().StringVar(&o.output, "output", ".", "folder path to add Gitops resources")
	initCmd.MarkFlagRequired("output")
	initCmd.Flags().StringVarP(&o.prefix, "prefix", "p", "", "add a prefix to the environment names")
	initCmd.Flags().BoolVarP(&o.skipChecks, "skip-checks", "b", false, "skip Tekton installation checks")
	initCmd.Flags().StringVar(&o.imageRepo, "image-repo", "", "image repository in this form <registry>/<username>/<repository> or <project>/<app> for internal registry")
	initCmd.Flags().StringVar(&o.deploymentPath, "deployment-path", "deploy", "deployment folder path name")
	initCmd.Flags().StringVar(&o.internalRegistryHostname, "internal-registry-hostname", "image-registry.openshift-image-registry.svc:5000", "internal image registry hostname")
	initCmd.MarkFlagRequired("gitops-repo")
	initCmd.MarkFlagRequired("gitops-webhook-secret")
	initCmd.MarkFlagRequired("dockercfgjson")
	initCmd.MarkFlagRequired("image-repo")

	return initCmd
}
