package webhook

import (
	"fmt"

	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/spf13/cobra"
)

type options struct {
	accessToken string
	isCICD      bool
	pipelines   string
	serviceName string
	isInsecure  bool
	*genericclioptions.Context
}

// Complete completes createOptions after they've been created
func (o *options) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	return nil
}

// Validate validates the createOptions based on completed values
func (o *options) Validate() (err error) {

	if o.isCICD {
		if o.serviceName != "" {
			return fmt.Errorf("Only one of --cicd or --service-name can be specified")
		}
	} else {
		if o.serviceName == "" {
			return fmt.Errorf("One of --cicd or --service-name must be specified")
		}
	}

	return nil
}

func (o *options) setFlags(command *cobra.Command) {

	// pipeline option
	command.Flags().StringVar(&o.pipelines, "pipelines", "pipelines.yaml", "path to pipelines file")

	// access-token option
	command.Flags().StringVar(&o.accessToken, "access-token", "", "access token to be used to create Git repository webhook")
	command.MarkFlagRequired("access-token")

	// cicd option
	command.Flags().BoolVar(&o.isCICD, "cicd", false, "provide this flag if the target Git repository is a CI/CD configuration repository")

	// service-name option
	command.Flags().StringVar(&o.serviceName, "service-name", "", "provide service-name if the target Git repository is a service's source repository")

	// insecure option
	command.Flags().BoolVar(&o.isInsecure, "insecure", false, "provide this flag if the Event Listenr external HTTP endpoint does not use TLS")

}
