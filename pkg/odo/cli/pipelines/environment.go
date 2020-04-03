package pipelines

import (
	"fmt"
	"strings"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// EnvRecommendedCommandName the recommended command name
	EnvRecommendedCommandName = "env"
)

var (
	envExample = ktemplates.Examples(`
	# Add a new environment to GitOps
	%[1]s 
	`)

	envLongDesc  = ktemplates.LongDesc(`Add a new environment to the GitOps repository`)
	envShortDesc = `Add a new environment`
)

// EnvParameters encapsulates the parameters for the odo pipelines init command.
type EnvParameters struct {
	envName string
	output string
	prefix string
	gitOpsRepo string
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewEnvParameters bootstraps a InitParameters instance.
func NewEnvParameters() *EnvParameters {
	return &EnvParameters{}
}

// Complete completes EnvParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (eo *EnvParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if eo.prefix != "" && !strings.HasSuffix(eo.prefix, "-") {
		eo.prefix = eo.prefix + "-"
	}
	return nil
}

// Validate validates the parameters of the EnvParameters.
func (eo *EnvParameters) Validate() error {
	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(strings.Split(eo.gitOpsRepo, "/")) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", eo.gitOpsRepo)
	}
	return nil
}

// Run runs the project bootstrap command.
func (eo *EnvParameters) Run() error {
	options:=pipelines.EnvParameters{
		GitOpsRepo: eo.gitOpsRepo,
		EnvName: eo.envName,
		Output: eo.output,
		Prefix: eo.prefix,
	}

	return pipelines.Env(&options)
}

// NewCmdEnv creates the project env command.
func NewCmdEnv(name, fullName string) *cobra.Command {
	o := NewEnvParameters()

	envCmd := &cobra.Command{
		Use:     name,
		Short:   envShortDesc,
		Long:    envLongDesc,
		Example: fmt.Sprintf(envExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	envCmd.Flags().StringVar(&o.gitOpsRepo, "gitops-repo", "", "CI/CD pipelines configuration Git repository in this form <username>/<repository>")
	envCmd.MarkFlagRequired("gitops-repo")
	envCmd.Flags().StringVar(&o.envName, "env-name", "", "name of the environment/namespace")
	envCmd.MarkFlagRequired("env-name")
	envCmd.Flags().StringVar(&o.output, "output", ".", "folder path to add Gitops resources")
	envCmd.Flags().StringVarP(&o.prefix, "prefix", "p", "", "add a prefix to the environment names")

	return envCmd
}
