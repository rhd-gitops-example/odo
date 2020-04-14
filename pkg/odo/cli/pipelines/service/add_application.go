package service

import (
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// AddServiceRecommendedCommandName the recommended command name
	AddServiceRecommendedCommandName = "add"
)

var (
	AddExample = ktemplates.Examples(`
  # Add applications to OpenShift pipelines in a cluster
  %[1]s
  `)

	AddLongDesc  = ktemplates.LongDesc(`Add applications with GitOps CI/CD Pipelines`)
	AddShortDesc = `Add bootstrapped Application repo pipelines`
)

// AddParameter encapsulates the parameters for the odo pipelines services add command.
type AddParameters struct {
	gitopsRepo           string
	gitopsWebhookSecret  string
	output               string
	prefix               string
	appName              string
	serviceWebhookSecret string
	serviceImageRepo     string
	envName              string
	serviceGitRepo       string
	skipChecks           bool

	*genericclioptions.Context
}

// NewAddParameters() bootstraps a AddParameters instance.
func NewAddParameters() *AddParameters {
	return &AddParameters{}
}

// Complete completes AddParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *AddParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if io.prefix != "" && !strings.HasSuffix(io.prefix, "-") {
		io.prefix = io.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the AddParameters
func (io *AddParameters) Validate() error {
	if len(strings.Split(io.gitopsRepo, "/")) != 2 || len(strings.Split(io.serviceGitRepo, "/")) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", io.gitopsRepo)
	}

	return nil
}

// Run runs the project bootstrap command.
func (io *AddParameters) Run() error {
	options := pipelines.AddParameters{
		GitopsRepo:           io.gitopsRepo,
		GitopsWebhookSecret:  io.gitopsWebhookSecret,
		Output:               io.output,
		AppName:              io.appName,
		ServiceWebhookSecret: io.serviceWebhookSecret,
		ServiceImageRepo:     io.serviceImageRepo,
		EnvName:              io.envName,
		ServicesGitRepo:      io.serviceGitRepo,
		SkipChecks:           io.skipChecks,
		Prefix:               io.prefix,
	}

	return pipelines.Add_Application(&options)
}

// NewCmdAddService creates the project add service command.
func NewCmdAddService(name, fullName string) *cobra.Command {
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
	addCmd.Flags().StringVar(&o.gitopsWebhookSecret, "gitops-webhook-secret", "", "provide the GitHub webhook secret for gitops repository")
	addCmd.Flags().StringVar(&o.output, "output", "", "folder path to add Gitops resources")
	addCmd.Flags().StringVar(&o.prefix, "prefix", "", "add a prefix to the environment names")
	addCmd.Flags().StringVar(&o.appName, "app-name", "", "CI/CD pipelines configuration Git repository in this form <username>/<repository>")
	addCmd.Flags().StringVar(&o.serviceWebhookSecret, "service-webhook-secret", "", "Provide the webhook secret of the app git repository")
	addCmd.Flags().StringVar(&o.serviceImageRepo, "service-image-repo", "", "Image repository name in form <username>/<repository>")
	addCmd.Flags().StringVar(&o.envName, "env-name", "", "Add the name of the environment(namespace) to which the pipelines should be bootstrapped")
	addCmd.Flags().StringVar(&o.serviceGitRepo, "services-git-repo", "", "Add the docker auth.json file path")
	addCmd.Flags().BoolVarP(&o.skipChecks, "skip-checks", "b", true, "skip Tekton installation checks")
	addCmd.MarkFlagRequired("gitops-repo")
	addCmd.MarkFlagRequired("app-name")
	addCmd.MarkFlagRequired("service-webhook-secret")
	addCmd.MarkFlagRequired("env-name")
	addCmd.MarkFlagRequired("service-dockercfgjson")

	return addCmd
}
