package service

import (
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"

	backend "github.com/openshift/odo/pkg/pipelines/webhook"
)

type options struct {
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

// Complete completes options after they've been created
func (o *options) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	return nil

}

// Validate validates the options based on completed values
func (o *options) Validate() (err error) {
	return nil
}

func (o *options) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.gitRepoURL, "git-repo-url", "", "source Git repository URL")
	cmd.Flags().StringVar(&o.webhookSecret, "webhook-secret", "", "source Git repository webhook secret")
	cmd.Flags().StringVar(&o.appName, "app-name", "", "the name of the application where the service will be added")
	cmd.Flags().StringVar(&o.serviceName, "service-name", "", "the name of the service to be added")
	cmd.Flags().StringVar(&o.envName, "env-name", "", "the name of the environment where the service will be added")
	cmd.Flags().StringVar(&o.imageRepo, "image-repo", "", "used to push built images")
	cmd.Flags().StringVar(&o.internalRegistryHostname, "internal-registry-hostname", "image-registry.openshift-image-registry.svc:5000", "internal image registry hostname")
	cmd.Flags().StringVar(&o.pipelinesFilePath, "pipelines-file", "pipelines.yaml", "path to pipelines file")

	// required flags
	_ = cmd.MarkFlagRequired("service-name")
	_ = cmd.MarkFlagRequired("app-name")
	_ = cmd.MarkFlagRequired("env-name")

}

func (o *options) getAppServiceNames() *backend.QualifiedServiceName {

	return &backend.QualifiedServiceName{
		EnvironmentName: o.envName,
		ServiceName:     o.serviceName,
	}
}
