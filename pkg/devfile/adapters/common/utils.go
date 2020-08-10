package common

import (
	"k8s.io/klog"
	"os"
	"strconv"

	devfileParser "github.com/openshift/odo/pkg/devfile/parser"
	"github.com/openshift/odo/pkg/devfile/parser/data"
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
)

// PredefinedDevfileCommands encapsulates constants for predefined devfile commands
type PredefinedDevfileCommands string

// DevfileEventType encapsulates constants for devfile events
type DevfileEventType string

const (
	// DefaultDevfileInitCommand is a predefined devfile command for init
	DefaultDevfileInitCommand PredefinedDevfileCommands = "devinit"

	// DefaultDevfileBuildCommand is a predefined devfile command for build
	DefaultDevfileBuildCommand PredefinedDevfileCommands = "devbuild"

	// DefaultDevfileRunCommand is a predefined devfile command for run
	DefaultDevfileRunCommand PredefinedDevfileCommands = "devrun"

	// DefaultDevfileDebugCommand is a predefined devfile command for debug
	DefaultDevfileDebugCommand PredefinedDevfileCommands = "debugrun"

	// SupervisordInitContainerName The init container name for supervisord
	SupervisordInitContainerName = "copy-supervisord"

	// Default Image that will be used containing the supervisord binary and assembly scripts
	// use GetBootstrapperImage() function instead of this variable
	defaultBootstrapperImage = "registry.access.redhat.com/ocp-tools-4/odo-init-container-rhel8:1.1.5"

	// SupervisordControlCommand sub command which stands for control
	SupervisordControlCommand = "ctl"

	// SupervisordVolumeName Create a custom name and (hope) that users don't use the *exact* same name in their deployment (occlient.go)
	SupervisordVolumeName = "odo-supervisord-shared-data"

	// SupervisordMountPath The supervisord Mount Path for the container mounting the supervisord volume
	SupervisordMountPath = "/opt/odo/"

	// SupervisordBinaryPath The supervisord binary path inside the container volume mount
	SupervisordBinaryPath = "/opt/odo/bin/supervisord"

	// SupervisordConfFile The supervisord configuration file inside the container volume mount
	SupervisordConfFile = "/opt/odo/conf/devfile-supervisor.conf"

	// OdoInitImageContents The path to the odo init image contents
	OdoInitImageContents = "/opt/odo-init/."

	// ENV variable to overwrite image used to bootstrap SupervisorD in S2I and Devfile builder Image
	bootstrapperImageEnvName = "ODO_BOOTSTRAPPER_IMAGE"

	// BinBash The path to sh executable
	BinBash = "/bin/sh"

	// Default volume size for volumes defined in a devfile
	volumeSize = "5Gi"

	// EnvProjectsRoot is the env defined for /projects where component mountSources=true
	EnvProjectsRoot = "PROJECTS_ROOT"

	// EnvOdoCommandRunWorkingDir is the env defined in the runtime component container which holds the work dir for the run command
	EnvOdoCommandRunWorkingDir = "ODO_COMMAND_RUN_WORKING_DIR"

	// EnvOdoCommandRun is the env defined in the runtime component container which holds the run command to be executed
	EnvOdoCommandRun = "ODO_COMMAND_RUN"

	// EnvOdoCommandDebugWorkingDir is the env defined in the runtime component container which holds the work dir for the debug command
	EnvOdoCommandDebugWorkingDir = "ODO_COMMAND_DEBUG_WORKING_DIR"

	// EnvOdoCommandDebug is the env defined in the runtime component container which holds the debug command to be executed
	EnvOdoCommandDebug = "ODO_COMMAND_DEBUG"

	// EnvDebugPort is the env defined in the runtime component container which holds the debug port for remote debugging
	EnvDebugPort = "DEBUG_PORT"

	// ShellExecutable is the shell executable
	ShellExecutable = "/bin/sh"

	// SupervisordCtlSubCommand is the supervisord sub command ctl
	SupervisordCtlSubCommand = "ctl"

	// PreStart is a devfile event
	PreStart DevfileEventType = "preStart"

	// PostStart is a devfile event
	PostStart DevfileEventType = "postStart"

	// PreStop is a devfile event
	PreStop DevfileEventType = "preStop"

	// PostStop is a devfile event
	PostStop DevfileEventType = "postStop"
)

// CommandNames is a struct to store the default and adapter names for devfile commands
type CommandNames struct {
	DefaultName string
	AdapterName string
}

// isContainer checks if the component is a container
func isContainer(component common.DevfileComponent) bool {
	// Currently odo only uses devfile components of type container, since most of the Che registry devfiles use it
	if component.Container != nil {
		klog.V(4).Infof("Found component \"%v\" with name \"%v\"\n", common.ContainerComponentType, component.Container.Name)
		return true
	}
	return false
}

// isVolume checks if the component is a volume
func isVolume(component common.DevfileComponent) bool {
	if component.Volume != nil {
		klog.V(4).Infof("Found component \"%v\" with name \"%v\"\n", common.VolumeComponentType, component.Volume.Name)
		return true
	}
	return false
}

// GetBootstrapperImage returns the odo-init bootstrapper image
func GetBootstrapperImage() string {
	if env, ok := os.LookupEnv(bootstrapperImageEnvName); ok {
		return env
	}
	return defaultBootstrapperImage
}

// GetDevfileContainerComponents iterates through the components in the devfile and returns a list of devfile container components
func GetDevfileContainerComponents(data data.DevfileData) []common.DevfileComponent {
	var components []common.DevfileComponent
	// Only components with aliases are considered because without an alias commands cannot reference them
	for _, comp := range data.GetAliasedComponents() {
		if isContainer(comp) {
			components = append(components, comp)
		}
	}
	return components
}

// GetDevfileVolumeComponents iterates through the components in the devfile and returns a map of devfile volume components
func GetDevfileVolumeComponents(data data.DevfileData) map[string]common.DevfileComponent {
	volumeNameToVolumeComponent := make(map[string]common.DevfileComponent)
	// Only components with aliases are considered because without an alias commands cannot reference them
	for _, comp := range data.GetComponents() {
		if isVolume(comp) {
			volumeNameToVolumeComponent[comp.Volume.Name] = comp
		}
	}
	return volumeNameToVolumeComponent
}

// getCommandsByGroup gets commands by the group kind
func getCommandsByGroup(data data.DevfileData, groupType common.DevfileCommandGroupType) []common.DevfileCommand {
	var commands []common.DevfileCommand

	for _, command := range data.GetCommands() {
		commandGroup := command.GetGroup()
		if commandGroup != nil && commandGroup.Kind == groupType {
			commands = append(commands, command)
		}
	}

	return commands
}

// GetVolumes iterates through the components in the devfile and returns a map of container name to the devfile volumes
func GetVolumes(devfileObj devfileParser.DevfileObj) map[string][]DevfileVolume {
	containerComponents := GetDevfileContainerComponents(devfileObj.Data)
	volumeNameToVolumeComponent := GetDevfileVolumeComponents(devfileObj.Data)

	// containerNameToVolumes is a map of the Devfile container name to the Devfile container Volumes
	containerNameToVolumes := make(map[string][]DevfileVolume)
	for _, containerComp := range containerComponents {
		for _, volumeMount := range containerComp.Container.VolumeMounts {
			size := volumeSize

			// check if there is a volume component name against the container component volume mount name
			if volumeComp, ok := volumeNameToVolumeComponent[volumeMount.Name]; ok {
				// If there is a volume size mentioned in the devfile, use it
				if len(volumeComp.Volume.Size) > 0 {
					size = volumeComp.Volume.Size
				}
			}

			vol := DevfileVolume{
				Name:          volumeMount.Name,
				ContainerPath: volumeMount.Path,
				Size:          size,
			}
			containerNameToVolumes[containerComp.Container.Name] = append(containerNameToVolumes[containerComp.Container.Name], vol)
		}
	}
	return containerNameToVolumes
}

// IsEnvPresent checks if the env variable is present in an array of env variables
func IsEnvPresent(envVars []common.Env, envVarName string) bool {
	for _, envVar := range envVars {
		if envVar.Name == envVarName {
			return true
		}
	}

	return false
}

// IsPortPresent checks if the port is present in the endpoints array
func IsPortPresent(endpoints []common.Endpoint, port int) bool {
	for _, endpoint := range endpoints {
		if endpoint.TargetPort == int32(port) {
			return true
		}
	}

	return false
}

// IsRestartRequired returns if restart is required for devrun command
func IsRestartRequired(command common.DevfileCommand) bool {
	var restart = true
	var err error
	rs, ok := command.Exec.Attributes["restart"]
	if ok {
		restart, err = strconv.ParseBool(rs)
		// Ignoring error here as restart is true for all error and default cases.
		if err != nil {
			klog.V(4).Info("Error converting restart attribute to bool")
		}
	}

	return restart
}

// GetCommandsMap returns a mapping of all of devfile command names to their corresponding DevfileCommand struct
// Allowing us to easily retrieve the DevfileCommand of any command listed in a composite command
func GetCommandsMap(commands []common.DevfileCommand) map[string]common.DevfileCommand {
	commandsMap := make(map[string]common.DevfileCommand)

	for _, command := range commands {
		commandsMap[command.GetID()] = command
	}
	return commandsMap
}

// GetComponentEnvVar returns true if a list of env vars contains the specified env var
// If the env exists, it returns the value of it
func GetComponentEnvVar(env string, envs []common.Env) string {
	for _, envVar := range envs {
		if envVar.Name == env {
			return envVar.Value
		}
	}
	return ""
}
