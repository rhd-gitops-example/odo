package webhook

import(
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/manifest"
	"github.com/openshift/odo/pkg/manifest/ioutils"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// createWebhookRecommendedCommandName the recommended command name
	createWebhookRecommendedCommandName = "create"
)

var (
	createWebhookExample = ktemplates.Examples(`
	# create a new webhook to GitOps
	%[1]s 
	`)

	createWebhookLongDesc  = ktemplates.LongDesc(`reate a new webhook to the GitOps repository`)
	createWebhookShortDesc = `reate a new webhook.`
)

type createWebhookParameter struct {
	*scm.Client
	repoName string
	// generic context options common to all commands
	*genericclioptions.Context

}

// NewAddEnvParameters bootstraps a AddEnvParameters instance.
func createWebhookParameters() *createWebhookParameters {
	return &createWebhookParameters{}
}

// Complete completes AddEnvParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated webhook url nicer to read.
func (w *createWebhookParameters) Complete(gitRepoURL string, cmd *cobra.Command, args []string) error {
	if w.prefix != "" && !strings.HasSuffix(eo.prefix, "-") {
		w.prefix = w.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the createWebhookParameters.
func (w *createWebhhokParameters) Validate() error {
	return nil
}

// Run runs the project bootstrap command.
func (w *createWebhookParameters) Run() error {
	options := manifest.createWebhookParameters{
		*scm.Client,
		repoName : eo.
	}

	return manifest.new(gitRepoURL, token)
}

// NewCmdAddEnv creates the project add environment command.
func NewCmdcreateWebhook(gitRepoURL, token string) *cobra.Command {
	c := createWebhookParameters()

	createWebhookCmd := &cobra.Command{
		Use:     gitRepoURL,
		Short:   addEnvShortDesc,
		Long:    addEnvLongDesc,
		Example: fmt.Sprintf(createWebhookExample,token),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(c, cmd, args)
		},
	}

	createWebhookCmd.Flags().StringVar(&c.gitRepoURL, "git_Repo_URL","", "Url to create new webhook")
	createWebhookCmd.MarkFlagRequired("git_Repo_URL")
	//createWebhookCmd.Flags().StringVar(&c.sc)
	return createWebhookCmd
}
