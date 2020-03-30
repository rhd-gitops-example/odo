package pipelines

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// InitRecommendedCommandName the recommended command name
	AddRecommendedCommandName = "add-service"
)

var (
	AddExample = ktemplates.Examples(`
	# Add applications to OpenShift pipelines in a cluster
	%[1]s 
	`)

	AddLongDesc  = ktemplates.LongDesc(`Add applications with GitOps CI/CD Pipelines`)
	AddShortDesc = `Add bootstrapped Application repo pipelines`
)

// InitParameters encapsulates the parameters for the odo pipelines init command.
type AddParameters struct {
	gitopsRepo string
	gitopsWebhookSecret string
	output string
	prefix string
	appGitRepo           string
	appWebhookSecret string
	appImageRepo string
	envName string
	dockerCfgJson string
	skipChecks bool
  
	*genericclioptions.Context
}

// NewInitParameters bootstraps a InitParameters instance.
func NewAddParameters() *AddParameters {
	return &AddParameters{}
}

// Complete completes InitParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *AddParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	
	return nil
}

// Validate validates the parameters of the InitParameters.
func (io *AddParameters) Validate() error {

	return nil
}

// Run runs the project bootstrap command.
func (io *AddParameters) Run() error {
	options := pipelines.AddParameters{
		GitopsRepo: io.gitopsRepo,
		GitopsWebhookSecret : io.gitopsWebhookSecret,
		Output: io.output,
		AppGitRepo: io.appGitRepo,
		AppWebhookSecret: io.appWebhookSecret,
		AppImageRepo: io.appImageRepo,
		EnvName: io.envName,
		DockerCfgJson: io.dockerCfgJson,
		SkipChecks: io.skipChecks,
	}

	return pipelines.Add(&options)
}

// NewCmdInit creates the project init command.
func NewCmdAdd(name, fullName string) *cobra.Command {
	o := NewAddParameters()

	addCmd := &cobra.Command{
		Use:     name,
		Short:   AddShortDesc,
		Long:    AddLongDesc,
		Example: fmt.Sprintf(AddExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	addCmd.Flags().StringVar(&o.gitopsRepo, "gitops-repo", "", "CI/CD pipelines configuration Git repository in this form <username>/<repository>")
	addCmd.MarkFlagRequired("gitops-repo")
	addCmd.Flags().StringVar(&o.gitopsWebhookSecret, "gitops-webhook-secret", "", "provide the GitHub webhook secret for gitops repository")
	addCmd.MarkFlagRequired("gitops-webhook-secret")
	addCmd.Flags().StringVar(&o.output, "output", "", "folder path to add Gitops resources")
	addCmd.Flags().StringVar(&o.prefix, "prefix", "", "add a prefix to the environment names")
	addCmd.Flags().StringVar(&o.appGitRepo, "app-git-repo", "", "CI/CD pipelines configuration Git repository in this form <username>/<repository>")
	addCmd.MarkFlagRequired("app-git-repo")
	addCmd.Flags().StringVar(&o.appWebhookSecret, "app-webhook-secret", "", "Provide the webhook secret of the app git repository")
	addCmd.MarkFlagRequired("app-webhook-secret")
	addCmd.Flags().StringVar(&o.appImageRepo, "app-image-repo", "", "Image repository name in form <username>/<repository>")
	addCmd.MarkFlagRequired("app-image-repo")
	addCmd.Flags().StringVar(&o.envName, "env-name", "", "Add the name of the environment(namespace) to which the pipelines should be bootstrapped")
	addCmd.MarkFlagRequired("env-name")
	addCmd.Flags().StringVar(&o.dockerCfgJson, "dockercfgjson", "", "Add the docker auth.json file path")
	addCmd.Flags().BoolVarP(&o.skipChecks, "skip-checks", "b", false, "skip Tekton installation checks")


	

	return addCmd
}
