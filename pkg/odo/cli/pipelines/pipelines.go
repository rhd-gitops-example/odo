package pipelines

import (
	"fmt"

	"github.com/spf13/cobra"

	odoutil "github.com/openshift/odo/pkg/odo/util"
)

// RecommendedCommandName is the recommended pipelines command name.
const RecommendedCommandName = "pipelines"

// NewCmdComponent implements the component odo command
func NewCmdComponent(name, fullName string) *cobra.Command {

	initialiseCmd := NewCmdInitialise(InitialiseRecommendedCommandName, odoutil.GetFullName(fullName, InitialiseRecommendedCommandName))
	var pipelinesCmd = &cobra.Command{
		Use:   name,
		Short: "Manage pipelines",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, InitialiseRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	pipelinesCmd.Flags().AddFlagSet(initialiseCmd.Flags())
	pipelinesCmd.AddCommand(initialiseCmd)

	pipelinesCmd.Annotations = map[string]string{"command": "main"}
	pipelinesCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return pipelinesCmd
}
