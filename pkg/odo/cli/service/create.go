package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/openshift/odo/pkg/log"
	"github.com/openshift/odo/pkg/odo/cli/service/ui"
	commonui "github.com/openshift/odo/pkg/odo/cli/ui"
	"github.com/openshift/odo/pkg/odo/genericclioptions"
	"github.com/openshift/odo/pkg/odo/util/completion"
	"github.com/openshift/odo/pkg/odo/util/experimental"
	"github.com/openshift/odo/pkg/odo/util/validation"
	svc "github.com/openshift/odo/pkg/service"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	scv1beta1 "github.com/kubernetes-sigs/service-catalog/pkg/apis/servicecatalog/v1beta1"
	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"k8s.io/klog"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	createRecommendedCommandName = "create"
	equivalentTemplate           = "{{.CmdFullName}} {{.ServiceType}}" +
		"{{if .ServiceName}} {{.ServiceName}}{{end}}" +
		" --app {{.Application}}" +
		" --project {{.Project}}" +
		"{{if .Plan}} --plan {{.Plan}}{{end}}" +
		"{{range $key, $value := .ParametersMap}} -p {{$key}}={{$value}}{{end}}"
)

var (
	createExample = ktemplates.Examples(`
    # Create new postgresql service from service catalog using dev plan and name my-postgresql-db.
    %[1]s dh-postgresql-apb my-postgresql-db --plan dev -p postgresql_user=luke -p postgresql_password=secret`)

	createOperatorExample = ktemplates.Examples(`
	# Create new EtcdCluster service from etcdoperator.v0.9.4 operator.
	%[1]s etcdoperator.v0.9.4/EtcdCluster`)

	createShortDesc = `Create a new service from service catalog using the plan defined and deploy it on OpenShift.`

	createLongDesc = ktemplates.LongDesc(`
Create a new service from service catalog using the plan defined and deploy it on OpenShift.

A --plan must be passed along with the service type. Parameters to configure the service are passed as key=value pairs.

For a full list of service types, use: 'odo catalog list services'`)

	createShortDescExperimental = `Create a new service from Operator Hub or Service Catalog and deploy it on OpenShift.`

	createLongDescExperimental = ktemplates.LongDesc(`
Create a new service from Operator Hub or Service Catalog and deploy it on OpenShift.

When creating a service using Operator Hub, provide a service name along with Operator name.

When creating a service using Service Catalog, a --plan must be passed along with the service type. Parameters to configure the service are passed as key=value pairs.

For a full list of service types, use: 'odo catalog list services'`)
)

// ServiceCreateOptions encapsulates the options for the odo service create command
type ServiceCreateOptions struct {
	// parameters hold the user-provided values for service class parameters via flags (populated by cobra)
	parameters []string
	// Plan is the selected service plan
	Plan string
	// ServiceType corresponds to the service class name
	ServiceType string
	// ServiceName is how the service will be named and known by odo
	ServiceName string
	// ParametersMap is populated from the flag-provided values (parameters) and/or the interactive mode and is the expected format by the business logic
	ParametersMap map[string]string
	// interactive specifies whether the command operates in interactive mode or not
	interactive bool
	// outputCLI specifies whether to output the non-interactive version of the command or not
	outputCLI bool
	// CmdFullName records the command's full name
	CmdFullName string
	// whether or not to wait for the service to be ready
	wait bool
	// generic context options common to all commands
	*genericclioptions.Context
	// Context to use when creating service. This will use app and project values from the context
	componentContext string
	// Custom Resrouce to create service from
	CustomResource string
	// Custom Resrouce's Definition fetched from alm-examples
	CustomResourceDefinition map[string]interface{}
	// Group of the GVR
	group string
	// Version of the GVR
	version string
	// Resource of the GVR
	resource string
	// If set to true, DryRun prints the yaml that will create the service
	DryRun bool
	// Location of the file in which yaml specification of CR is stored.
	// TODO: remove this after service create's interactive mode supports creating operator backed services
	fromFile string
}

// NewServiceCreateOptions creates a new ServiceCreateOptions instance
func NewServiceCreateOptions() *ServiceCreateOptions {
	return &ServiceCreateOptions{}
}

// DynamicCRD holds the original CR obtained from the Operator (a CSV), or user
// (when they use --from-file flag), and few other attributes that are likely
// to be used to validate a CRD before creating a service from it
type DynamicCRD struct {
	// contains the CR as obtained from CSV or user
	OriginalCRD map[string]interface{}
}

func NewDynamicCRD() *DynamicCRD {
	return &DynamicCRD{}
}

// Complete completes ServiceCreateOptions after they've been created
func (o *ServiceCreateOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 || !cmd.HasFlags() {
		o.interactive = true
	}

	if experimental.IsExperimentalModeEnabled() {
		o.Context = genericclioptions.NewDevfileContext(cmd)
	} else if o.componentContext != "" {
		o.Context = genericclioptions.NewContext(cmd)
	} else {
		o.Context = genericclioptions.NewContextCreatingAppIfNeeded(cmd)
	}

	client := o.Client

	var class scv1beta1.ClusterServiceClass

	if experimental.IsExperimentalModeEnabled() {
		// we don't support interactive mode for Operator Hub yet
		o.interactive = false

		// if user has just used "odo service create", simply return
		if o.fromFile == "" && len(args) == 0 {
			return
		}

		// if user wants to create service from file and use a name given on CLI
		if o.fromFile != "" {
			if len(args) == 1 {
				o.ServiceName = args[0]
			}
			return
		}
	}

	if o.interactive {
		classesByCategory, err := client.GetServiceClassesByCategory()
		if err != nil {
			return fmt.Errorf("unable to retrieve service classes: %v", err)
		}

		if len(classesByCategory) == 0 {
			return fmt.Errorf("no available service classes")
		}

		class, o.ServiceType = ui.SelectClassInteractively(classesByCategory)

		plans, err := client.GetMatchingPlans(class)
		if err != nil {
			return fmt.Errorf("couldn't retrieve plans for class %s: %v", class.GetExternalName(), err)
		}

		var svcPlan scv1beta1.ClusterServicePlan
		// if there is only one available plan, we select it
		if len(plans) == 1 {
			for k, v := range plans {
				o.Plan = k
				svcPlan = v
			}
			klog.V(4).Infof("Plan %s was automatically selected since it's the only one available for service %s", o.Plan, o.ServiceType)
		} else {
			// otherwise select the plan interactively
			o.Plan = ui.SelectPlanNameInteractively(plans, "Which service plan should we use ")
			svcPlan = plans[o.Plan]
		}

		o.ParametersMap = ui.EnterServicePropertiesInteractively(svcPlan)
		o.ServiceName = ui.EnterServiceNameInteractively(o.ServiceType, "How should we name your service ", o.validateServiceName)
		o.outputCLI = commonui.Proceed("Output the non-interactive version of the selected options")
		o.wait = commonui.Proceed("Wait for the service to be ready")
	} else {
		if experimental.IsExperimentalModeEnabled() {
			// split the name provided on CLI and populate servicetype & customresource
			o.ServiceType, o.CustomResource, err = svc.SplitServiceKindName(args[0])
			if err != nil {
				return fmt.Errorf("invalid service name, use the format <operator-type>/<crd-name>")
			}
		} else {
			o.ServiceType = args[0]
		}
		// if only one arg is given, then it is considered as service name and service type both
		// ONLY if not running in Experimental mode
		if !experimental.IsExperimentalModeEnabled() {
			// This is because an operator with name
			// "etcdoperator.v0.9.4-clusterwide" would lead to creation of a
			// serice with name like
			// "etcdoperator.v0.9.4-clusterwide-c47rf28l56" and that would fail
			// because it's an invalid name for k8s/OCP
			o.ServiceName = o.ServiceType
		}
		// if two args are given, first is service type and second one is service name
		if len(args) == 2 {
			o.ServiceName = args[1]
		}

		// we convert the param list provided in the format of key=value list
		// to a map
		o.ParametersMap = make(map[string]string)
		for _, kv := range o.parameters {
			kvSlice := strings.Split(kv, "=")
			// key value not provided in format of key=value
			if len(kvSlice) != 2 {
				return errors.New("parameters not provided in key=value format")
			}
			o.ParametersMap[kvSlice[0]] = kvSlice[1]
		}
	}

	return
}

// validateServiceName adopts the Validator interface and checks that the name of the service being created is valid
func (o *ServiceCreateOptions) validateServiceName(i interface{}) (err error) {
	s := i.(string)
	err = validation.ValidateName(s)
	if err != nil {
		return err
	}
	exists, err := svc.SvcExists(o.Client, s, o.Application)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s service already exists in the current application", o.ServiceName)
	}
	return
}

// outputNonInteractiveEquivalent outputs the populated options as the equivalent command that would be used in non-interactive mode
func (o *ServiceCreateOptions) outputNonInteractiveEquivalent() string {
	if o.outputCLI {
		var tpl bytes.Buffer
		t := template.Must(template.New("service-create-cli").Parse(equivalentTemplate))
		e := t.Execute(&tpl, o)
		if e != nil {
			panic(e) // shouldn't happen
		}
		return strings.TrimSpace(tpl.String())
	}
	return ""
}

// Validate validates the ServiceCreateOptions based on completed values
func (o *ServiceCreateOptions) Validate() (err error) {
	// if we are in interactive mode, all values are already valid
	if o.interactive {
		return nil
	}

	// we want to find an Operator only if something's passed to the crd flag on CLI
	if experimental.IsExperimentalModeEnabled() {
		d := NewDynamicCRD()
		// if the user wants to create service from a file, we check for
		// existence of file and validate if the requested operator and CR
		// exist on the cluster
		if o.fromFile != "" {
			if _, err := os.Stat(o.fromFile); err != nil {
				return errors.Wrap(err, "unable to find specified file")
			}

			// Parse the file to find Operator and CR info
			fileContents, err := ioutil.ReadFile(o.fromFile)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(fileContents, &d.OriginalCRD)
			if err != nil {
				return err
			}

			// Check if the operator and the CR exist on cluster
			var csv olm.ClusterServiceVersion
			o.CustomResource, csv, err = svc.GetCSV(o.KClient, d.OriginalCRD)
			if err != nil {
				return err
			}

			// all is well, let's populate the fields required for creating operator backed service
			o.group, o.version, o.resource, err = svc.GetGVRFromOperator(csv, o.CustomResource)
			if err != nil {
				return err
			}

			err = d.validateMetadataInCRD()
			if err != nil {
				return err
			}

			if o.ServiceName != "" && !o.DryRun {
				// First check if service with provided name already exists
				svcFullName := strings.Join([]string{o.CustomResource, o.ServiceName}, "/")
				exists, err := svc.OperatorSvcExists(o.KClient, svcFullName)
				if err != nil {
					return err
				}
				if exists {
					return fmt.Errorf("service %q already exists; please provide a different name or delete the existing service first.", svcFullName)
				}

				d.setServiceName(o.ServiceName)
			} else {
				o.ServiceName, err = d.getServiceNameFromCRD()
				if err != nil {
					return err
				}
			}

			// CRD is valid. We can use it further to create a service from it.
			o.CustomResourceDefinition = d.OriginalCRD

			return nil
		} else if o.CustomResource != "" {
			// make sure that CSV of the specified ServiceType exists
			csv, err := o.KClient.GetClusterServiceVersion(o.ServiceType)
			if err != nil {
				// error only occurs when OperatorHub is not installed.
				// k8s does't have it installed by default but OCP does
				return err
			}

			almExample, err := svc.GetAlmExample(csv, o.CustomResource, o.ServiceType)
			if err != nil {
				return err
			}

			d.OriginalCRD = almExample

			o.group, o.version, o.resource, err = svc.GetGVRFromOperator(csv, o.CustomResource)
			if err != nil {
				return err
			}

			if o.ServiceName != "" && !o.DryRun {
				// First check if service with provided name already exists
				svcFullName := strings.Join([]string{o.CustomResource, o.ServiceName}, "/")
				exists, err := svc.OperatorSvcExists(o.KClient, svcFullName)
				if err != nil {
					return err
				}
				if exists {
					return fmt.Errorf("service %q already exists; please provide a different name or delete the existing service first.", svcFullName)
				}

				d.setServiceName(o.ServiceName)
			}

			err = d.validateMetadataInCRD()
			if err != nil {
				return err
			}

			// CRD is valid. We can use it further to create a service from it.
			o.CustomResourceDefinition = d.OriginalCRD

			return nil
		} else {
			// This block is executed only when user has neither provided a
			// file nor a valid `odo service create <operator-name>` to start
			// the service from an Operator. So we raise and error because the
			// correct way is to execute:
			// `odo service create <operator-name>/<crd-name>`

			return fmt.Errorf("please use a valid command to start an Operator backed service; desired format: %q", "odo service create <operator-name>/<crd-name>")
		}
	}
	// make sure the service type exists
	classPtr, err := o.Client.GetClusterServiceClass(o.ServiceType)
	if err != nil {
		return errors.Wrap(err, "unable to create service because Service Catalog is not enabled in your cluster")
	}
	if classPtr == nil {
		return fmt.Errorf("service %v doesn't exist\nRun 'odo catalog list services' to see a list of supported services.\n", o.ServiceType)
	}

	// check plan
	plans, err := o.Client.GetMatchingPlans(*classPtr)
	if err != nil {
		return err
	}
	if len(o.Plan) == 0 {
		// when the plan has not been supplied, if there is only one available plan, we select it
		if len(plans) == 1 {
			for k := range plans {
				o.Plan = k
			}
			klog.V(4).Infof("Plan %s was automatically selected since it's the only one available for service %s", o.Plan, o.ServiceType)
		} else {
			return fmt.Errorf("no plan was supplied for service %v.\nPlease select one of: %v\n", o.ServiceType, strings.Join(ui.GetServicePlanNames(plans), ","))
		}
	} else {
		// when the plan has been supplied, we need to make sure it exists
		if _, ok := plans[o.Plan]; !ok {
			return fmt.Errorf("plan %s is invalid for service %v.\nPlease select one of: %v\n", o.Plan, o.ServiceType, strings.Join(ui.GetServicePlanNames(plans), ","))
		}
	}
	//validate service name
	return o.validateServiceName(o.ServiceName)
}

// Run contains the logic for the odo service create command
func (o *ServiceCreateOptions) Run() (err error) {
	s := &log.Status{}
	if experimental.IsExperimentalModeEnabled() {
		// in case of an opertor backed service, name of the service is
		// provided by the yaml specification in alm-examples. It might also
		// happen that a user spins up Service Catalog based service in
		// experimental mode but we're taking a bet against that for now, so
		// the user won't get to see service name in the log message
		if !o.DryRun {
			log.Infof("Deploying service of type: %s", o.CustomResource)
			s = log.Spinner("Deploying service")
			defer s.End(false)
		}
	} else {
		log.Infof("Deploying service %s of type: %s", o.ServiceName, o.ServiceType)
	}

	if experimental.IsExperimentalModeEnabled() && o.CustomResource != "" {
		// if experimental mode is enabled and o.CustomResource is not empty, we're expected to create an Operator backed service
		if o.DryRun {
			// if it's dry run, only print the alm-example (o.CustomResourceDefinition) and exit
			jsonCR, err := json.MarshalIndent(o.CustomResourceDefinition, "", "  ")
			if err != nil {
				return err
			}

			// convert json to yaml
			yamlCR, err := yaml.JSONToYAML(jsonCR)
			if err != nil {
				return err
			}

			log.Info(string(yamlCR))

			return nil
		} else {
			err = svc.CreateOperatorService(o.KClient, o.group, o.version, o.resource, o.CustomResourceDefinition)
		}
	} else {
		// otherwise just create a ServiceInstance
		err = svc.CreateService(o.Client, o.ServiceName, o.ServiceType, o.Plan, o.ParametersMap, o.Application)
	}
	if err != nil {
		return err
	}
	s.End(true)

	if o.wait {
		s = log.Spinner("Waiting for service to come up")
		_, err = o.Client.WaitAndGetSecret(o.ServiceName, o.Project)
		if err == nil {
			s.End(true)
			log.Successf(`Service '%s' is ready for use`, o.ServiceName)
		}
	} else {
		log.Successf(`Service '%s' was created`, o.ServiceName)
		log.Italic("\nProgress of the provisioning will not be reported and might take a long time\nYou can see the current status by executing 'odo service list'")
	}

	// Information on what to do next
	log.Infof("Optionally, link %s to your component by running: 'odo link <component-name>'", o.ServiceType)

	equivalent := o.outputNonInteractiveEquivalent()
	if len(equivalent) > 0 {
		log.Info("Equivalent command:\n" + ui.StyledOutput(equivalent, "cyan"))
	}
	return
}

// NewCmdServiceCreate implements the odo service create command.
func NewCmdServiceCreate(name, fullName string) *cobra.Command {
	o := NewServiceCreateOptions()
	o.CmdFullName = fullName
	serviceCreateCmd := &cobra.Command{
		Use:     name + " <service_type> --plan <plan_name> [service_name]",
		Short:   createShortDesc,
		Long:    createLongDesc,
		Example: fmt.Sprintf(createExample, fullName),
		Args:    cobra.RangeArgs(0, 2),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	if experimental.IsExperimentalModeEnabled() {
		serviceCreateCmd.Use += fmt.Sprintf(" [flags]\n  %s <operator_type>/<crd_name> [service_name] [flags]", o.CmdFullName)
		serviceCreateCmd.Short = createShortDescExperimental
		serviceCreateCmd.Long = createLongDescExperimental
		serviceCreateCmd.Example += "\n\n" + fmt.Sprintf(createOperatorExample, fullName)
		serviceCreateCmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Print the yaml specificiation that will be used to create the service")
		// remove this feature after enabling service create interactive mode for operator backed services
		serviceCreateCmd.Flags().StringVar(&o.fromFile, "from-file", "", "Path to the file containing yaml specification to use to start operator backed service")
	}

	serviceCreateCmd.Flags().StringVar(&o.Plan, "plan", "", "The name of the plan of the service to be created")
	serviceCreateCmd.Flags().StringArrayVarP(&o.parameters, "parameters", "p", []string{}, "Parameters of the plan where a parameter is expressed as <key>=<value")
	serviceCreateCmd.Flags().BoolVarP(&o.wait, "wait", "w", false, "Wait until the service is ready")
	genericclioptions.AddContextFlag(serviceCreateCmd, &o.componentContext)
	completion.RegisterCommandHandler(serviceCreateCmd, completion.ServiceClassCompletionHandler)
	completion.RegisterCommandFlagHandler(serviceCreateCmd, "plan", completion.ServicePlanCompletionHandler)
	completion.RegisterCommandFlagHandler(serviceCreateCmd, "parameters", completion.ServiceParameterCompletionHandler)
	return serviceCreateCmd
}

// validateMetadataInCRD validates if the CRD has metadata.name field and returns an error
func (d *DynamicCRD) validateMetadataInCRD() error {
	metadata, ok := d.OriginalCRD["metadata"].(map[string]interface{})
	if !ok {
		// this condition is satisfied if there's no metadata at all in the provided CRD
		return fmt.Errorf("couldn't find \"metadata\" in the yaml; need metadata start the service")
	}

	if _, ok := metadata["name"].(string); ok {
		// found the metadata.name; no error
		return nil
	}
	return fmt.Errorf("couldn't find metadata.name in the yaml; provide a name for the service")
}

// setServiceName modifies the CRD to contain user provided name on the CLI
// instead of using the default one in almExample
func (d *DynamicCRD) setServiceName(name string) {
	metaMap := d.OriginalCRD["metadata"].(map[string]interface{})

	for k := range metaMap {
		if k == "name" {
			metaMap[k] = name
			return
		}
		// if metadata doesn't have 'name' field, we set it up
		metaMap["name"] = name
	}
}

// getServiceNameFromCRD fetches the service name from metadata.name field of the CRD
func (d *DynamicCRD) getServiceNameFromCRD() (string, error) {
	metadata, ok := d.OriginalCRD["metadata"].(map[string]interface{})
	if !ok {
		// this condition is satisfied if there's no metadata at all in the provided CRD
		return "", fmt.Errorf("couldn't find \"metadata\" in the yaml; need metadata.name to start the service")
	}

	if name, ok := metadata["name"].(string); ok {
		// found the metadata.name; no error
		return name, nil
	}
	return "", fmt.Errorf("couldn't find metadata.name in the yaml; provide a name for the service")
}
