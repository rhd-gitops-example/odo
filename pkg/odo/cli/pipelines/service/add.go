package service

import (
	"fmt"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
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
	options
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

	err := pipelines.AddService(&pipelines.AddServiceParameters{
		AppName:                  o.appName,
		EnvName:                  o.envName,
		GitRepoURL:               o.gitRepoURL,
		ImageRepo:                o.imageRepo,
		InternalRegistryHostname: o.internalRegistryHostname,
		PipelinesFilePath:        o.pipelinesFilePath,
		ServiceName:              o.serviceName,
		WebhookSecret:            o.webhookSecret,
	}, ioutils.NewFilesystem())

	if err != nil {
		return err
	}
	log.Successf("Created Service %s sucessfully at environment %s.", o.serviceName, o.envName)
	return nil

}

func newCmdAdd(name, fullName string) *cobra.Command {
	o := &AddOptions{}

	cmd := &cobra.Command{
		Use:     name,
		Short:   addShortDesc,
		Long:    addLongDesc,
		Example: fmt.Sprintf(addExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	o.setFlags(cmd)
	return cmd
}
