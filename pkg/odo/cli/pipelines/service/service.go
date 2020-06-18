package service

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	odoutil "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)

// RecommendedCommandName is the recommended environment command name.
const RecommendedCommandName = "service"

// AddOptions encapsulates the parameters for service add command
type AddOptions struct {
	appName                  string
	envName                  string
	gitRepoURL               string
	imageRepo                string
	internalRegistryHostname string
	pipelinesFilePath        string
	serviceName              string
	webhookSecret            string

	// generic context options common to all commands
	*genericclioptions.Context
}

// NewCmd creates a new environment command
func NewCmd(name, fullName string) *cobra.Command {

	addCmd := newCmdAdd(addRecommendedCommandName, odoutil.GetFullName(fullName, addRecommendedCommandName))

	var cmd = &cobra.Command{

		Use:   name,
		Short: "Manage services in an environment",
		Long:  "Manage services in a GitOps environment where service source repositories are synchronized",
		Example: fmt.Sprintf("%s\n%s\n\n  See sub-commands individually for more examples",
			fullName, addRecommendedCommandName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// cmd.PersistentFlags().StringVar(&o.gitRepoURL, "git-repo-url", "", "source Git repository URL")
	// cmd.PersistentFlags().StringVar(&o.webhookSecret, "webhook-secret", "", "source Git repository webhook secret")
	// cmd.PersistentFlags().StringVar(&o.appName, "app-name", "", "the name of the application where the service will be added")
	// cmd.PersistentFlags().StringVar(&o.serviceName, "service-name", "", "the name of the service to be added")
	// cmd.PersistentFlags().StringVar(&o.envName, "env-name", "", "the name of the environment where the service will be added")
	// cmd.PersistentFlags().StringVar(&o.imageRepo, "image-repo", "", "used to push built images")
	// cmd.PersistentFlags().StringVar(&o.internalRegistryHostname, "internal-registry-hostname", "image-registry.openshift-image-registry.svc:5000", "internal image registry hostname")
	// cmd.PersistentFlags().StringVar(&o.pipelinesFilePath, "pipelines-file", "pipelines.yaml", "path to pipelines file")
	// cmd.Flags().AddFlagSet(addCmd.Flags())
	// // required flags
	// _ = cmd.MarkFlagRequired("service-name")
	// _ = cmd.MarkFlagRequired("app-name")
	// _ = cmd.MarkFlagRequired("env-name")
	cmd.AddCommand(addCmd)

	cmd.Annotations = map[string]string{"command": "main"}
	cmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	return cmd
}
