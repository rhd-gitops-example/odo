package pipelines

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/utility"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	// BootstrapRecommendedCommandName the recommended command name
	BootstrapRecommendedCommandName = "bootstrap"
)

type drivers []string

var supportedDrivers = drivers{
	"github",
	"gitlab",
}

func (d drivers) supported(s string) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}

var (
	bootstrapExample = ktemplates.Examples(`
    # Bootstrap OpenShift pipelines.
    %[1]s 
    `)

	bootstrapLongDesc  = ktemplates.LongDesc(`Bootstrap GitOps CI/CD Manifest`)
	bootstrapShortDesc = `Bootstrap pipelines with a starter configuration`
)

// BootstrapParameters encapsulates the parameters for the odo pipelines init command.
type BootstrapParameters struct {
	*pipelines.BootstrapOptions
	// generic context options common to all commands
	*genericclioptions.Context
}

// NewBootstrapParameters bootstraps a BootstrapParameters instance.
func NewBootstrapParameters() *BootstrapParameters {
	return &BootstrapParameters{
		BootstrapOptions: &pipelines.BootstrapOptions{},
	}
}

// Complete completes BootstrapParameters after they've been created.
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *BootstrapParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	io.Prefix = utility.MaybeCompletePrefix(io.Prefix)
	io.GitOpsRepoURL = utility.AddGitSuffixIfNecessary(io.GitOpsRepoURL)
	io.ServiceRepoURL = utility.AddGitSuffixIfNecessary(io.ServiceRepoURL)

	if io.PrivateRepoDriver != "" {
		host, err := hostFromURL(io.GitOpsRepoURL)
		if err != nil {
			return err
		}
		identifier := factory.NewDriverIdentifier(factory.Mapping(host, io.PrivateRepoDriver))
		factory.DefaultIdentifier = identifier
	}
	return nil
}

// Validate validates the parameters of the BootstrapParameters.
func (io *BootstrapParameters) Validate() error {
	gr, err := url.Parse(io.GitOpsRepoURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %s: %w", io.GitOpsRepoURL, err)
	}

	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(utility.RemoveEmptyStrings(strings.Split(gr.Path, "/"))) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", strings.Trim(gr.Path, ".git"))
	}

	if io.PrivateRepoDriver != "" {
		if !supportedDrivers.supported(io.PrivateRepoDriver) {
			return fmt.Errorf("invalid driver type: %q", io.PrivateRepoDriver)
		}
	}
	return nil
}

// Run runs the project bootstrap command.
func (io *BootstrapParameters) Run() error {
	err := pipelines.Bootstrap(io.BootstrapOptions, ioutils.NewFilesystem())
	if err != nil {
		return err
	}
	log.Success("Bootstrapped GitOps sucessfully.")
	return nil
}

// NewCmdBootstrap creates the project init command.
func NewCmdBootstrap(name, fullName string) *cobra.Command {
	o := NewBootstrapParameters()

	bootstrapCmd := &cobra.Command{
		Use:     name,
		Short:   bootstrapShortDesc,
		Long:    bootstrapLongDesc,
		Example: fmt.Sprintf(bootstrapExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	bootstrapCmd.Flags().StringVar(&o.GitOpsRepoURL, "gitops-repo-url", "", "Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git")
	bootstrapCmd.Flags().StringVar(&o.GitOpsWebhookSecret, "gitops-webhook-secret", "", "Provide a secret that we can use to authenticate incoming hooks from your Git hosting service for the GitOps repository. (if not provided, it will be auto-generated)")
	bootstrapCmd.Flags().StringVar(&o.OutputPath, "output", ".", "Path to write GitOps resources")
	bootstrapCmd.Flags().StringVarP(&o.Prefix, "prefix", "p", "", "Add a prefix to the environment names(Dev, stage,prod,cicd etc.) to distinguish and identify individual environments")
	bootstrapCmd.Flags().StringVar(&o.DockerConfigJSONFilename, "dockercfgjson", "~/.docker/config.json", "Filepath to config.json which authenticates the image push to the desired image registry ")
	bootstrapCmd.Flags().StringVar(&o.InternalRegistryHostname, "image-repo-internal-registry-hostname", "image-registry.openshift-image-registry.svc:5000", "Host-name for internal image registry e.g. docker-registry.default.svc.cluster.local:5000, used if you are pushing your images to the internal image registry")
	bootstrapCmd.Flags().StringVar(&o.ImageRepo, "image-repo", "", "Image repository of the form <registry>/<username>/<repository> or <project>/<app> which is used to push newly built images")
	bootstrapCmd.Flags().StringVar(&o.SealedSecretsService.Namespace, "sealed-secrets-ns", "kube-system", "Namespace in which the Sealed Secrets operator is installed, automatically generated secrets are encrypted with this operator")
	bootstrapCmd.Flags().StringVar(&o.SealedSecretsService.Name, "sealed-secrets-svc", "sealed-secrets-controller", "Name of the Sealed Secrets service that encrypts secrets")
	bootstrapCmd.Flags().StringVar(&o.StatusTrackerAccessToken, "status-tracker-access-token", "", "Used to authenticate requests to push commit-statuses to your Git hosting service")
	bootstrapCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "Overwrites previously existing GitOps configuration (if any)")

	bootstrapCmd.MarkFlagRequired("gitops-repo-url")
	bootstrapCmd.Flags().StringVar(&o.ServiceRepoURL, "service-repo-url", "", "Provide the URL for your Service repository e.g. https://github.com/organisation/service.git")
	bootstrapCmd.Flags().StringVar(&o.ServiceWebhookSecret, "service-webhook-secret", "", "Provide a secret that we can use to authenticate incoming hooks from your Git hosting service for the Service repository. (if not provided, it will be auto-generated)")
	bootstrapCmd.Flags().StringVar(&o.PrivateRepoDriver, "private-repo-driver", "", "If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab")

	bootstrapCmd.MarkFlagRequired("service-repo-url")
	return bootstrapCmd
}
