package service

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	addRecommendedCommandName = "add"
)

var (
	addExample = ktemplates.Examples(`	Add a Service to an environment in GitOps 
	%[1]s`)

	addLongDesc  = ktemplates.LongDesc(`Add a Service to an environment in GitOps`)
	addShortDesc = `Add a new service`
)

// AddOptions encapsulates the parameters for service add command
type AddOptions struct {
	appName       string
	envName       string
	gitRepoURL    string
	manifest      string
	webhookSecret string
	serviceName   string

	// generic context options common to all commands
	*genericclioptions.Context
}

// Complete is called when the command is completed
func (o *AddOptions) Complete(name string, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate validates the parameters of the EnvParameters.
func (o *AddOptions) Validate() error {
	return nil
}

// Run runs the project bootstrap command.
func (o *AddOptions) Run() error {

	return pipelines.AddService(o.gitRepoURL, o.appName, o.serviceName, o.webhookSecret, o.envName, o.manifest, ioutils.NewFilesystem())
}

func newCmdAdd(name, fullName string) *cobra.Command {
	o := &AddOptions{}

	cmd := &cobra.Command{
		Use:     name,
		Short:   addExample,
		Long:    addExample,
		Example: fmt.Sprintf(addExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	cmd.Flags().StringVar(&o.gitRepoURL, "git-repo-url", "", "folder/path to add Gitops resources")
	cmd.Flags().StringVar(&o.webhookSecret, "webhook-secret", "", "folder/path to add Gitops resources")
	cmd.Flags().StringVar(&o.appName, "app-name", "", "the name of the application where the service will be added")
	cmd.Flags().StringVar(&o.serviceName, "service-name", "", "the name of the application where the service will be added")
	cmd.Flags().StringVar(&o.envName, "env-name", "", "the name of the environment where the service will be added")
	cmd.Flags().StringVar(&o.manifest, "manifest", "pipelines.yaml", "path to manifest file")

	// required flags
	cmd.MarkFlagRequired("service-name")
	cmd.MarkFlagRequired("app-name")
	cmd.MarkFlagRequired("env-name")
	return cmd
}
