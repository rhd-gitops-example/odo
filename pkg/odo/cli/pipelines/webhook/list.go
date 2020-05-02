package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	backend "github.com/openshift/odo/pkg/pipelines/webhook"
	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const listRecommendedCommandName = "list"

var (
	listExample = ktemplates.Examples(`	# List Git repository webhooks 
	%[1]s`)
)

type listOptions struct {
	options
}

func newListOptions() *listOptions {
	return &listOptions{}
}

// Run contains the logic for the odo command
func (o *listOptions) Run() (err error) {
	listeners, err := backend.List(o.accessToken, o.pipelines, o.getAppServiceNames(), o.isCICD, o.isInsecure)
	if err != nil {
		for _, listener := range listeners {
			fmt.Printf(" XXXXXXXXXXxx %s\n", listener)
		}
	}
	return nil
}

// NewCmdList creates a new "list" command
func NewCmdList(name, fullName string) *cobra.Command {
	o := newListOptions()
	command := &cobra.Command{
		Use:     name,
		Short:   "List existing webhooks.",
		Long:    "List existing Git repository webhooks that triggers CI/CD pipeline runs.",
		Example: fmt.Sprintf(createExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	o.setFlags(command)
	return command
}
