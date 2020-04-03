package pipelines

import (
	"github.com/spf13/cobra"

	"fmt"
	odoutil "github.com/openshift/odo/pkg/odo/util"
)

// RecommendedCommandName is the recommended pipelines command name.
const AddRecommendedCommandName = "add"

// NewCmdComponent implements the component odo command
func NewCmdAdd(name, fullName string) *cobra.Command {

	// envCmd := NewCmdInit(InitRecommendedCommandName, odoutil.GetFullName(fullName, InitRecommendedCommandName))
	envCmd := NewCmdEnv(EnvRecommendedCommandName, odoutil.GetFullName(fullName, EnvRecommendedCommandName))
	var addCmd = &cobra.Command{
		Use:   name,
		Short: "Add a new component to GitOps",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, EnvRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	addCmd.Flags().AddFlagSet(envCmd.Flags())
	addCmd.AddCommand(envCmd)

	addCmd.Annotations = map[string]string{"command": "main"}
	addCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return addCmd
}
