package webhook
import(
	"fmt"

	odoutil "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)
// WebhookRecommendedCommandName is the recommended environment command name.
const AddWebhookRecommendedCommandName = "webhook"

func NewCmdWebhook(gitRepoURL, token string) *cobra.command{
	addWebhookCmd := NewCmdAddWebhook(AddWebhookRecommendedCommandName,odoutil.gitRepoURL(gitRepoURL,AddWebhookRecommendedCommandName))

	var webhookcmd = &cobra.Command{
		use: gitRepoURL,
		Short: "create a webhook URL"
		Examples: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			gitRepoURL, AddWebhookRecommendedCommandName),
		Run: func(cmd *cobra.Command, args[]string){

		},

	} 

	webhookcmd.Flags().AddFlagSet(addWebhookCmd.Flags())
	webhookcmd.AddCommand(addWebhookCmd)

	webhookcmd.Annotations = map[string] string{"command":"main"}
	webhookcmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return webhookcmd
}