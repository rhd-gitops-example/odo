package project

import (
	"fmt"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/project"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const setRecommendedCommandName = "set"

var (
	setExample = ktemplates.Examples(`
	# Set the active project
	%[1]s myproject
	`)

	setLongDesc = ktemplates.LongDesc(`Set the active project`)

	setShortDesc = `Set the current active project`
)

// ProjectSetOptions encapsulates the options for the odo project set command
type ProjectSetOptions struct {

	// if supplied then only print the project name
	projectShortFlag bool

	// the name of the project that needs to be set as active
	projectName string

	// generic context options common to all commands
	*genericclioptions.Context
}

// NewProjectSetOptions creates a ProjectSetOptions instance
func NewProjectSetOptions() *ProjectSetOptions {
	return &ProjectSetOptions{}
}

// Complete completes ProjectSetOptions after they've been created
func (pso *ProjectSetOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	pso.Context = genericclioptions.NewContext(cmd)
	pso.projectName = args[0]

	return
}

// Validate validates the parameters of the ProjectSetOptions
func (pso *ProjectSetOptions) Validate() (err error) {

	exists, err := project.Exists(pso.Context, pso.projectName)

	if !exists {
		return fmt.Errorf("The project %s does not exist", pso.projectName)
	}

	return
}

// Run runs the project set command
func (pso *ProjectSetOptions) Run() (err error) {
	current := pso.Project
	err = project.SetCurrent(pso.Context, pso.projectName)
	if err != nil {
		return err
	}
	if pso.projectShortFlag {
		fmt.Print(pso.projectName)
	} else {
		if current == pso.projectName {
			log.Infof("Already on project : %v", pso.projectName)
		} else {
			log.Infof("Switched to project : %v", pso.projectName)
		}
	}
	return
}

// NewCmdProjectSet creates the project set command
func NewCmdProjectSet(name, fullName string) *cobra.Command {
	o := NewProjectSetOptions()

	projectSetCmd := &cobra.Command{
		Use:     name,
		Short:   setShortDesc,
		Long:    setLongDesc,
		Example: fmt.Sprintf(setExample, fullName),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	projectSetCmd.Flags().BoolVarP(&o.projectShortFlag, "short", "q", false, "If true, display only the project name")

	return projectSetCmd
}
