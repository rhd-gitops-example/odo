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
	AddRecommendedCommandName = "addApplication"
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
	gitopsRepo           string // repo to store Gitops resources e.g. org/repo
	  
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

	

	return addCmd
}
