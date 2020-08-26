package pipelines

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/ui"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/utility"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/namespaces"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	// WizardRecommendedCommandName the recommended command name
	WizardRecommendedCommandName = "wizard"

	sealedSecretsName   = "sealed-secrets-controller"
	sealedSecretsNS     = "kube-system"
	argoCDNS            = "argocd"
	argoCDOperatorName  = "argocd-operator"
	argoCDServerName    = "argocd-server"
	pipelinesOperatorNS = "openshift-operators"
)

var (
	WizardExample = ktemplates.Examples(`
    # Wizard OpenShift pipelines.
    %[1]s 
    `)

	WizardLongDesc  = ktemplates.LongDesc(`Wizard GitOps CI/CD Manifest`)
	WizardShortDesc = `Wizard pipelines with a starter configuration`
)

// WizardParameters encapsulates the parameters for the odo pipelines init command.
type WizardParameters struct {
	*pipelines.BootstrapOptions
	// generic context options common to all commands
	*genericclioptions.Context
}

type status interface {
	WarningStatus(status string)
	Start(status string, debug bool)
	End(status bool)
}

// NewWizardParameters Wizards a WizardParameters instance.
func NewWizardParameters() *WizardParameters {
	return &WizardParameters{
		BootstrapOptions: &pipelines.BootstrapOptions{},
	}
}

// Complete completes WizardParameters after they've been created.
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *WizardParameters) Complete(name string, cmd *cobra.Command, args []string) error {

	clientSet, err := namespaces.GetClientSet()
	if err != nil {
		return err
	}

	err = checkBootstrapDependencies(io, clientSet, log.NewStatus(os.Stdout))
	if err != nil {
		return err
	}

	// ask for sealed secrets only when default is absent
	if io.SealedSecretsService == (types.NamespacedName{}) {
		io.SealedSecretsService.Name = ui.EnterSealedSecretService(&io.SealedSecretsService)

	}

	io.GitOpsRepoURL = ui.EnterGitRepo()
	io.GitOpsRepoURL = utility.AddGitSuffixIfNecessary(io.GitOpsRepoURL)
	if !isKnownDriver(io.GitOpsRepoURL) {
		io.PrivateRepoDriver = ui.SelectPrivateRepoDriver()
		host, err := hostFromURL(io.GitOpsRepoURL)
		if err != nil {
			return fmt.Errorf("failed to parse the gitops url: %w", err)
		}
		identifier := factory.NewDriverIdentifier(factory.Mapping(host, io.PrivateRepoDriver))
		factory.DefaultIdentifier = identifier
	}

	option := ui.SelectOptionImageRepository()
	if option == "Openshift Internal repository" {
		io.InternalRegistryHostname = ui.EnterInternalRegistry()
		io.ImageRepo = ui.EnterImageRepoInternalRegistry()
	} else {
		io.DockerConfigJSONFilename = ui.EnterDockercfg()
		fs := ioutils.NewFilesystem()
		_, err := pipelines.CheckFileExists(fs, io.DockerConfigJSONFilename)
		if err != nil {
			return err
		}
		io.ImageRepo = ui.EnterImageRepoExternalRepository()
	}
	io.GitOpsWebhookSecret = ui.EnterGitWebhookSecret()
	io.Prefix = ui.EnterPrefix()
	io.ServiceRepoURL = ui.EnterServiceRepoURL()
	io.Prefix = utility.MaybeCompletePrefix(io.Prefix)
	io.ServiceRepoURL = utility.AddGitSuffixIfNecessary(io.ServiceRepoURL)
	io.ServiceWebhookSecret = ui.EnterServiceWebhookSecret()
	commitStatusTrackerCheck := ui.SelectOptionCommitStatusTracker()
	if commitStatusTrackerCheck == "yes" {
		io.StatusTrackerAccessToken = ui.EnterStatusTrackerAccessToken(io.ServiceRepoURL)
	}
	io.OutputPath = ui.EnterOutputPath()
	io.Overwrite = true
	return nil
}

func checkBootstrapDependencies(io *WizardParameters, kubeClient kubernetes.Interface, spinner status) error {
	var errs []error
	client := utility.NewClient(kubeClient)
	log.Progressf("\nChecking dependencies\n")

	spinner.Start("Checking if Sealed Secrets is installed with the default configuration", false)
	err := client.CheckIfSealedSecretsExists(types.NamespacedName{Namespace: sealedSecretsNS, Name: sealedSecretsName})
	setSpinnerStatus(spinner, "Please install Sealed Secrets from https://github.com/bitnami-labs/sealed-secrets/releases", err)
	if err == nil {
		io.SealedSecretsService.Name = sealedSecretsName
		io.SealedSecretsService.Namespace = sealedSecretsNS
	} else if !errors.IsNotFound(err) {
		return clusterErr(err.Error())
	}

	spinner.Start("Checking if ArgoCD Operator is installed with the default configuration", false)
	err = client.CheckIfArgoCDExists(argoCDNS)
	setSpinnerStatus(spinner, "Please install ArgoCD operator from OperatorHub", err)
	if err != nil {
		if !errors.IsNotFound(err) {
			return clusterErr(err.Error())
		}
		errs = append(errs, err)
	}

	spinner.Start("Checking if OpenShift Pipelines Operator is installed with the default configuration", false)
	err = client.CheckIfPipelinesExists(pipelinesOperatorNS)
	setSpinnerStatus(spinner, "Please install OpenShift Pipelines operator from OperatorHub", err)
	if err != nil {
		if !errors.IsNotFound(err) {
			return clusterErr(err.Error())
		}
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("Failed to satisfy the required dependencies")
	}
	return nil
}

func setSpinnerStatus(spinner status, warningMsg string, err error) {
	if err != nil {
		if errors.IsNotFound(err) {
			spinner.WarningStatus(warningMsg)
		}
		spinner.End(false)
		return
	}
	spinner.End(true)
}

// Validate validates the parameters of the WizardParameters.
func (io *WizardParameters) Validate() error {
	gr, err := url.Parse(io.GitOpsRepoURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %s: %w", io.GitOpsRepoURL, err)
	}

	// TODO: this won't work with GitLab as the repo can have more path elements.
	if len(utility.RemoveEmptyStrings(strings.Split(gr.Path, "/"))) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", strings.Trim(gr.Path, ".git"))
	}

	return nil
}

// Run runs the project Wizard command.
func (io *WizardParameters) Run() error {
	if io.ServiceRepoURL != "" {
		err := pipelines.Bootstrap(io.BootstrapOptions, ioutils.NewFilesystem())
		if err != nil {
			return err
		}
		log.Success("Bootstrapped GitOps sucessfully.")
	}
	return nil
}

// NewCmdWizard creates the project init command.
func NewCmdWizard(name, fullName string) *cobra.Command {
	o := NewWizardParameters()

	wizardCmd := &cobra.Command{
		Use:     name,
		Short:   WizardShortDesc,
		Long:    WizardLongDesc,
		Example: fmt.Sprintf(WizardExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	return wizardCmd
}

func clusterErr(errMsg string) error {
	return fmt.Errorf("Couldn't connect to cluster: %s", errMsg)
}

//returns the username/reponame from the url
func repoFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Path, "/")
	return strings.TrimSuffix(parts[len(parts)-2], ".git") + "/" + strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
}

func isKnownDriver(repoURL string) bool {
	host, err := hostFromURL(repoURL)
	if err != nil {
		return false
	}
	_, err = factory.DefaultIdentifier.Identify(host)
	if err == nil {
		return true
	}
	return false
}

func hostFromURL(s string) (string, error) {
	p, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return strings.ToLower(p.Host), nil
}
