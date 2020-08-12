package pipelines

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/ui"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/utility"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/pipelines"
	"github.com/openshift/odo/pkg/pipelines/ioutils"
	"github.com/openshift/odo/pkg/pipelines/namespaces"
	"github.com/openshift/odo/pkg/pipelines/scm"
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
//
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
		io.SealedSecretsService.Name = ui.EnterSealedSecretService()
		io.SealedSecretsService.Namespace = ui.EnterSealedSecretNamespace()
	}

	io.GitOpsRepoURL = ui.EnterGitRepo()
	io.GitOpsRepoURL = utility.AddGitSuffixIfNecessary(io.GitOpsRepoURL)
	_, err := scm.NewRepository(io.GitOpsRepoURL)
	if err != nil {
		return err
	}
	option := ui.SelectOptionImageRepository()
	if option == "Openshift Internal repository" {
		io.InternalRegistryHostname = ui.EnterInternalRegistry()
		io.ImageRepo = ui.EnterImageRepoInternalRegistry()

	} else {
		io.DockerConfigJSONFilename = ui.EnterDockercfg()
		fs := ioutils.NewFilesystem()
		_, err = pipelines.CheckFileExists(fs, io.DockerConfigJSONFilename)
		if err != nil {
			return err
		}
		io.ImageRepo = ui.EnterImageRepoExternalRepository()
	}
	io.GitOpsWebhookSecret = ui.EnterGitWebhookSecret()
	if ui.CheckSecretLength(io.GitOpsWebhookSecret) {
		return fmt.Errorf("The GitOps Webhook Secret length should 16 or more ")
	}
	io.SealedSecretsService.Name = ui.EnterSealedSecretService()
	io.SealedSecretsService.Namespace = ui.EnterSealedSecretNamespace()
	io.Prefix = ui.EnterPrefix()
	io.Prefix = utility.MaybeCompletePrefix(io.Prefix)
	io.ServiceRepoURL = utility.AddGitSuffixIfNecessary(io.ServiceRepoURL)
	_, err = scm.NewRepository(io.ServiceRepoURL)
	if err != nil {
		return err
	}
	io.ServiceWebhookSecret = ui.EnterServiceWebhookSecret()
	if ui.CheckSecretLength(io.ServiceWebhookSecret) {
		return fmt.Errorf("The GitOps Webhook Secret length should 16 or more ")
	}
	commitStatusTrackerCheck := ui.SelectOptionCommitStatusTracker()
	if commitStatusTrackerCheck == "yes" {
		io.StatusTrackerAccessToken = ui.EnterStatusTrackerAccessToken()
	}
	io.OutputPath = ui.EnterOutputPath(io.GitOpsRepoURL)
	exists, _ := ioutils.IsExisting(ioutils.NewFilesystem(), filepath.Join(io.OutputPath, "pipelines.yaml"))
	if exists {
		selectOverwriteOption := ui.SelectOptionOverwrite()
		if selectOverwriteOption == "no" {
			io.Overwrite = false
			return fmt.Errorf("Cannot create GitOps configuration since file exists at %s", io.OutputPath)
		}
	}
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
