package pipelines

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/cli/pipelines/service"
	"github.com/spf13/cobra"

	odoutil "github.com/openshift/odo/pkg/odo/util"
)

// RecommendedCommandName is the recommended pipelines command name.
const RecommendedCommandName = "pipelines"

// NewCmdComponent implements the component odo command
func NewCmdComponent(name, fullName string) *cobra.Command {

	initCmd := NewCmdInit(InitRecommendedCommandName, odoutil.GetFullName(fullName, InitRecommendedCommandName))
	ServiceCmd := service.NewCmdAddService(service.ServiceRecommendedCommandName, odoutil.GetFullName(fullName, service.ServiceRecommendedCommandName))
	var pipelinesCmd = &cobra.Command{
		Use:   name,
		Short: "Manage pipelines",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, InitRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	pipelinesCmd.Flags().AddFlagSet(initCmd.Flags())
	pipelinesCmd.AddCommand(initCmd)
	pipelinesCmd.AddCommand(ServiceCmd)

	pipelinesCmd.Annotations = map[string]string{"command": "main"}
	pipelinesCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return pipelinesCmd
}
