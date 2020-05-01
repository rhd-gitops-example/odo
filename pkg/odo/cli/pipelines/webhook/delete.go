package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	backend "github.com/openshift/odo/pkg/pipelines/webhook"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const deleteRecommendedCommandName = "delete"

var (
	deleteExample = ktemplates.Examples(`	# Delete a Git repository webhook 
	%[1]s`)
)

type deleteOptions struct {
	options
}

func newDeleteOptions() *deleteOptions {
	return &deleteOptions{}
}

// Run contains the logic for the odo command
func (o *deleteOptions) Run() (err error) {
	appName, serviceName := o.getAppServiceNames()
	return backend.Delete(o.accessToken, o.pipelines, appName, serviceName, o.isCICD, o.isInsecure)
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

	o.setFlags(command)
	return command
}
