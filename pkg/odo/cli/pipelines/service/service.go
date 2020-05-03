package service

import (
	"fmt"

	odoutil "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)

// EnvRecommendedCommandName is the recommended environment command name.
const ServiceRecommendedCommandName = "service"

// NewCmdEnv create a new environment command
func NewCmdService(name, fullName string) *cobra.Command {

	addServiceCmd := NewCmdAddService(AddServiceRecommendedCommandName, odoutil.GetFullName(fullName, AddServiceRecommendedCommandName))

	var serviceCmd = &cobra.Command{
		Use:   name,
		Short: "Manage an environment in GitOps",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, AddServiceRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	serviceCmd.Flags().AddFlagSet(addServiceCmd.Flags())
	serviceCmd.AddCommand(addServiceCmd)

	serviceCmd.Annotations = map[string]string{"command": "main"}
	serviceCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return serviceCmd
}
