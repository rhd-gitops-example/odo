package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	backend "github.com/openshift/odo/pkg/pipelines/webhook"
	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const createRecommendedCommandName = "create"

var (
	createExample = ktemplates.Examples(`	# Create a new Git repository webhook 
	%[1]s`)
)

type createOptions struct {
	isInsecure bool
	options
}

func newCreateOptions() *createOptions {
	return &createOptions{}
}

// Run contains the logic for the odo command
func (o *createOptions) Run() (err error) {
	return backend.Create(o.accessToken, o.pipelines, o.getAppServiceNames(), o.isCICD, o.isInsecure)
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

	o.setFlags(command)
	// insecure option
	command.Flags().BoolVar(&o.isInsecure, "insecure", false, "provide this flag if the Event Listenr external HTTP endpoint does not use TLS")
	return command
}
