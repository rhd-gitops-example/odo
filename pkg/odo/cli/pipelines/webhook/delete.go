package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const deleteRecommendedCommandName = "delete"

var (
	deleteExample = ktemplates.Examples(`	# Delete a Git repository webhook 
	%[1]s`)
)

type deleteOptions struct {
	accessToken string
	isCICD      bool
	pipelines   string
	serviceName string
	*genericclioptions.Context
}

func newDeleteOptions() *deleteOptions {
	return &deleteOptions{}
}

// Complete completes deleteOptions after they've been created
func (o *deleteOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	return
}

// Validate validates the deleteOptions based on completed values
func (o *deleteOptions) Validate() (err error) {

	if o.isCICD {
		if o.serviceName != "" {
			return fmt.Errorf("Only one of --cicd or --service-name can be specified")
		}
	} else {
		if o.serviceName == "" {
			return fmt.Errorf("One of --cicd or --service-name must be specified")
		}
	}

	return nil
}

// Run contains the logic for the odo command
func (o *deleteOptions) Run() (err error) {
	return
}

// NewCmdDelete creates a new "delete" command
func NewCmdDelete(name, fullName string) *cobra.Command {
	o := newDeleteOptions()
	command := &cobra.Command{
		Use:     name,
		Short:   "Delete webhooks.",
		Long:    "Delete all Git repository webhooks that trigger event to CI/CD Pipeline Event Listeners.",
		Example: fmt.Sprintf(deleteExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	// pipeline option
	command.Flags().StringVar(&o.pipelines, "pipelines", "pipelines.yaml", "path to pipelines file")

	// access-token option
	command.Flags().StringVar(&o.accessToken, "access-token", "", "access token to be used to create Git repository webhook")
	command.MarkFlagRequired("access-token")

	// cicd option
	command.Flags().BoolVar(&o.isCICD, "cicd", false, "provide this flag if the target Git repository is a CI/CD configuration repository")

	// service-name option
	command.Flags().StringVar(&o.serviceName, "service-name", "", "provide service-name if the target Git repository is a service's source repository")

	return command
}
