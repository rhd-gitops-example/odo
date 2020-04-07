package environment

import (
	"github.com/spf13/cobra"

	odoutil "github.com/openshift/odo/pkg/odo/util"
)

// EnvRecommendedCommandName is the recommended environment command name.
const EnvRecommendedCommandName = "environment"

// NewCmdEnv implements the component odo command
func NewCmdEnv(name, fullName string) *cobra.Command {

	addEnvCmd := NewCmdAddEnv(AddEnvRecommendedCommandName, odoutil.GetFullName(fullName, AddEnvRecommendedCommandName))
	var envCmd = &cobra.Command{
		Use:   name,
		Short: "Add a new component to GitOps",
		//Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
		//fullName, EnvRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	envCmd.Flags().AddFlagSet(addEnvCmd.Flags())
	envCmd.AddCommand(addEnvCmd)

	envCmd.Annotations = map[string]string{"command": "main"}
	envCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return envCmd
}
