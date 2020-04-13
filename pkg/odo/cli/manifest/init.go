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
	// InitRecommendedCommandName the recommended command name
	InitRecommendedCommandName = "init"
)

var (
	initExample = ktemplates.Examples(`
	# Initialize OpenShift GitOps manifest
	%[1]s 
	`)

	initLongDesc  = ktemplates.LongDesc(`Initialize GitOps manifest`)
	initShortDesc = `Initialize manifest`
)

// InitParameters encapsulates the parameters for the odo pipelines init command.
type InitParameters struct {
	gitOpsRepo          string // repo to store Gitops resources e.g. org/repo
	gitOpsWebhookSecret string // used to create Github's shared webhook secret for gitops repo
	output              string // path to add Gitops resources
	prefix              string // used to generate the environments in a shared cluster
	skipChecks          bool   // skip Tekton installation checks
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewInitParameters bootstraps a InitParameters instance.
func NewInitParameters() *InitParameters {
	return &InitParameters{}
}

// Complete completes InitParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *InitParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if io.prefix != "" && !strings.HasSuffix(io.prefix, "-") {
		io.prefix = io.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the InitParameters.
func (io *InitParameters) Validate() error {
	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(strings.Split(io.gitOpsRepo, "/")) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", io.gitOpsRepo)
	}
	return nil
}

// Run runs the project bootstrap command.
func (io *InitParameters) Run() error {
	options := manifest.InitParameters{
		GitOpsWebhookSecret: io.gitOpsWebhookSecret,
		GitOpsRepo:          io.gitOpsRepo,
		Output:              io.output,
		Prefix:              io.prefix,
		SkipChecks:          io.skipChecks,
	}
	return manifest.Init(&options)
}

// NewCmdInit creates the project init command.
func NewCmdInit(name, fullName string) *cobra.Command {
	o := NewInitParameters()

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
	initCmd.MarkFlagRequired("gitops-repo")
	initCmd.Flags().StringVar(&o.gitOpsWebhookSecret, "gitops-webhook-secret", "", "provide the GitHub webhook secret for GitOps repository")
	initCmd.MarkFlagRequired("gitops-webhook-secret")
	initCmd.Flags().StringVar(&o.output, "output", ".", "folder path to add GitOps resources")
	initCmd.MarkFlagRequired("output")
	initCmd.Flags().StringVarP(&o.prefix, "prefix", "p", "", "add a prefix to the environment names")
	return initCmd
}
