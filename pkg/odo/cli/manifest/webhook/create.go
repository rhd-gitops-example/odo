package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const createRecommendedCommandName = "create"

var (
	createExample = ktemplates.Examples(`	# Create a new Git repository webhook 
	%[1]s`)
)

// CreateOptions are opttions can be passed to webhook create command
type CreateOptions struct {
	*genericclioptions.Context
}

func newCreateOptions() *CreateOptions {
	return &CreateOptions{}
}

// Complete completes CreateOptions after they've been created
func (o *CreateOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	o.Context = genericclioptions.NewContext(cmd)
	return
}

// Validate validates the CreateOptions based on completed values
func (o *CreateOptions) Validate() (err error) {
	return nil
}

// Run contains the logic for the odo command
func (o *CreateOptions) Run() (err error) {
	return
}

// NewCmdCreate creates a new "create" command
func NewCmdCreate(name, fullName string) *cobra.Command {
	o := newCreateOptions()
	command := &cobra.Command{
		Use:     name,
		Short:   "Create a new webhook.",
		Long:    "Create a new Git repository webhook that triggers CI/CD pipeline runs.",
		Example: fmt.Sprintf(createExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	return command
}
