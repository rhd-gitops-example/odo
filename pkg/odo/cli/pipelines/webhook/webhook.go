package webhook

import (
	"fmt"

	odoutil "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended webhook command name.
const RecommendedCommandName = "webhook"

// NewCmdWebhook create a new webhook command
func NewCmdWebhook(name, fullName string) *cobra.Command {
	createCmd := NewCmdCreate(createRecommendedCommandName, odoutil.GetFullName(fullName, createRecommendedCommandName))
	deleteCmd := NewCmdDelete(deleteRecommendedCommandName, odoutil.GetFullName(fullName, deleteRecommendedCommandName))

	var webhookCmd = &cobra.Command{
		Use:   name,
		Short: "Manage Git repository webhooks",
		Long:  "Add/Delete Git repository webhooks that trigger CI/CD pipeline runs.",
		Example: fmt.Sprintf("%s\n%s\n%s\n\n  See sub-commands individually for more examples",
			fullName,
			createRecommendedCommandName,
			deleteRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	webhookCmd.AddCommand(createCmd)
	webhookCmd.AddCommand(deleteCmd)

	webhookCmd.Annotations = map[string]string{"command": "main"}
	webhookCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return webhookCmd
}
