package service

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubernetes/pkg/kubectl/util/templates"
)

const (
	// AddEnvRecommendedCommandName the recommended command name
	AddServiceRecommendedCommandName = "add"
)

var (
	addEnvExample = ktemplates.Examples(`
	# Add a new environment to GitOps
	%[1]s 
	`)

	addEnvLongDesc  = ktemplates.LongDesc(`Add a new environment to the GitOps repository`)
	addEnvShortDesc = `Add a new environment`
)

// AddEnvParameters encapsulates the parameters for the odo pipelines init command.
type AddServiceParameters struct {
	serviceGitRepo       string
	serviceWebhookSecret string
	appName              string
	envName              string
	output               string
	manifest             string
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewAddEnvParameters bootstraps a AddEnvParameters instance.
func NewAddServiceParameters() *AddServiceParameters {
	return &AddServiceParameters{}
}

// Complete completes AddEnvParameters after they've been created.
//
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (eo *AddServiceParameters) Complete(name string, cmd *cobra.Command, args []string) error {

	return nil
}

// Validate validates the parameters of the EnvParameters.
func (eo *AddServiceParameters) Validate() error {
	return nil
}

// Run runs the project bootstrap command.
func (eo *AddServiceParameters) Run() error {
	options := pipelines.ServiceParameters{
		ServiceGitRepo:       eo.serviceGitRepo,
		ServiceWebhookSecret: eo.serviceWebhookSecret,
		AppName:              eo.appName,
		EnvName:              eo.envName,
		Output:               eo.output,
		Manifest:             eo.manifest,
	}

	return pipelines.AddService(&options, ioutils.NewFilesystem())
}

// NewCmdAddEnv creates the project add environment command.
func NewCmdAddService(name, fullName string) *cobra.Command {
	o := NewAddServiceParameters()

	addServiceCmd := &cobra.Command{
		Use:     name,
		Short:   addEnvShortDesc,
		Long:    addEnvLongDesc,
		Example: fmt.Sprintf(addEnvExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	addServiceCmd.Flags().StringVar(&o.serviceGitRepo, "service-git-repo", "", "folder/path to add Gitops resources")
	addServiceCmd.Flags().StringVar(&o.serviceWebhookSecret, "service-webhook-secret", "", "folder/path to add Gitops resources")
	addServiceCmd.Flags().StringVar(&o.appName, "app-name", "", "folder/path to add Gitops resources")
	addServiceCmd.Flags().StringVar(&o.envName, "env-name", "", "name of the environment/namespace")
	addServiceCmd.Flags().StringVar(&o.output, "output", ".", "folder/path to add Gitops resources")
	addServiceCmd.Flags().StringVar(&o.manifest, "manifest", "", "folder/path to add Gitops resources")

	addServiceCmd.MarkFlagRequired("service-git-repo")
	addServiceCmd.MarkFlagRequired("service-webhook-secret")
	addServiceCmd.MarkFlagRequired("app-name")
	addServiceCmd.MarkFlagRequired("env-name")
	return addServiceCmd
}
