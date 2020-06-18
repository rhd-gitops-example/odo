package environment

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	odoutil "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)

// EnvRecommendedCommandName is the recommended environment command name.
const EnvRecommendedCommandName = "environment"

//AddEnvParameters struct for the list of persistent flags
type AddEnvParameters struct {
	envName       string
	output        string
	pipelinesFile string
	cluster       string
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewCmdEnv create a new environment command
func NewCmdEnv(name, fullName string) *cobra.Command {
	o := NewAddEnvParameters()

	addEnvCmd := NewCmdAddEnv(AddEnvRecommendedCommandName, odoutil.GetFullName(fullName, AddEnvRecommendedCommandName), o)

	var envCmd = &cobra.Command{
		Use:   name,
		Short: "Manage an environment in GitOps",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, AddEnvRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	envCmd.Flags().AddFlagSet(addEnvCmd.Flags())
	envCmd.AddCommand(addEnvCmd)
	envCmd.PersistentFlags().StringVar(&o.pipelinesFile, "pipelines-file", "pipelines.yaml", "The value of the pipelines path")
	envCmd.PersistentFlags().StringVar(&o.envName, "env-name", "", "The name of the environment to be mentioned")
	envCmd.MarkFlagRequired("env-name")

	envCmd.Annotations = map[string]string{"command": "main"}
	envCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return envCmd
}
