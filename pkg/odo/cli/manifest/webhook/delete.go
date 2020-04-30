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

// DeleteOptions are opttions can be passed to webhook delete command
type DeleteOptions struct {
	*genericclioptions.Context
}

func newDeleteOptions() *DeleteOptions {
	return &DeleteOptions{}
}

// Complete completes DeleteOptions after they've been created
func (o *DeleteOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	o.Context = genericclioptions.NewContext(cmd)
	return
}

// Validate validates the DeleteOptions based on completed values
func (o *DeleteOptions) Validate() (err error) {
	return nil
}

// Run contains the logic for the odo command
func (o *DeleteOptions) Run() (err error) {
	return
}

// NewCmdDelete creates a new "delete" command
func NewCmdDelete(name, fullName string) *cobra.Command {
	command := &cobra.Command{
		Use:     name,
		Short:   "Delete webhooks.",
		Long:    "Delete all Git repository webhooks that trigger event to CI/CD Pipeline Event Listeners.",
		Example: fmt.Sprintf(deleteExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(newDeleteOptions(), cmd, args)
		},
	}

	return command
}
