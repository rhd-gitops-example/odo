package pipelines

import (
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	InitialiseRecommendedCommandName = "initialise"
)

var (
	initialiseExample = ktemplates.Examples(`
	# Initialise OpenShift pipelines in a cluster
	%[1]s 
	`)

	initialiseLongDesc  = ktemplates.LongDesc(`Initialise GitOps CI/CD Pipelines`)
	initialiseShortDesc = `Initialise pipelines`
)

// InitialiseParameters encapsulates the paratmeters for the odo pipelines bootstrap
// command.
type InitialiseParameters struct {
	gitOpsRepo               string
	output                   string
	appGithubSecret          string
	appGithubRepo            string
	appImageRepo             string
	prefix                   string
	internalRegistryHostname string
	dockerConfigJSONFileName string
	githubToken              string
	skipChecks               bool
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewInitialiseParameters bootstraps a BootstrapParameters instance.
func NewInitialiseParameters() *InitialiseParameters {
	return &InitialiseParameters{}
}

// Complete completes BootstrapParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *InitialiseParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if io.prefix != "" && !strings.HasSuffix(io.prefix, "-") {
		io.prefix = io.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the BootstrapParameters.
func (io *InitialiseParameters) Validate() error {
	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(strings.Split(io.gitOpsRepo, "/")) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", io.gitOpsRepo)
	}
	return nil
}

// Run runs the project bootstrap command.
func (io *InitialiseParameters) Run() error {
	options := pipelines.InitialiseParameters{
		GitOpsRepo: io.gitOpsRepo,
		Prefix:     io.prefix,
		Output:     io.output,
	}

	return pipelines.Initialise(&options)
}

// NewCmdInitialise creates the project initialise command.
func NewCmdInitialise(name, fullName string) *cobra.Command {
	o := NewInitialiseParameters()

	initialiseCmd := &cobra.Command{
		Use:     name,
		Short:   initialiseShortDesc,
		Long:    initialiseLongDesc,
		Example: fmt.Sprintf(initialiseExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	initialiseCmd.Flags().StringVarP(&o.prefix, "prefix", "p", "", "add a prefix to the environment names")
	initialiseCmd.Flags().StringVarP(&o.githubToken, "status-tracker-token", "", "", "provide the Github token")
	initialiseCmd.Flags().StringVar(&o.dockerConfigJSONFileName, "dockerconfigjson", "", "Docker configuration json filename")
	initialiseCmd.Flags().StringVar(&o.gitOpsRepo, "gitops-repo", "", "git repository in this form <username>/<repository>")
	initialiseCmd.MarkFlagRequired("gitops-repo")
	initialiseCmd.Flags().StringVar(&o.appImageRepo, "app-image-repo", "", "image repository in this form <registry>/<username>/<repository> or <project>/<app> for internal registry")
	initialiseCmd.Flags().StringVar(&o.output, "output", "output", "folder path to add resources")
	initialiseCmd.MarkFlagRequired("output")
	initialiseCmd.Flags().BoolVarP(&o.skipChecks, "skip-checks", "b", false, "skip Tekton installation checks")
	initialiseCmd.Flags().StringVar(&o.internalRegistryHostname, "internal-registry-hostname", "image-registry.openshift-image-registry.svc:5000", "internal image registry hostname")

	return initialiseCmd
}
