package common

import (
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
	"github.com/openshift/odo/pkg/machineoutput"
)

// supervisorCommand encapsulates a supervisor-specific command
type supervisorCommand struct {
	adapter commandExecutor
	cmd     []string
	info    ComponentInfo
}

// newSupervisorInitCommand creates a command that initializes the supervisor for the specified devfile if needed
// nil is returned if no devfile-specified container needing supervisor initialization is found
func newSupervisorInitCommand(command common.DevfileCommand, adapter commandExecutor) (command, error) {
	cmd := []string{SupervisordBinaryPath, "-c", SupervisordConfFile, "-d"}
	info, err := adapter.SupervisorComponentInfo(command)
	if err != nil {
		adapter.Logger().ReportError(err, machineoutput.TimestampNow())
		return nil, err
	}
	if !info.IsEmpty() {
		return supervisorCommand{
			adapter: adapter,
			cmd:     cmd,
			info:    info,
		}, nil
	}
	return nil, nil
}

// newSupervisorStopCommand creates a command implementation that stops the specified command via the supervisor
func newSupervisorStopCommand(command common.DevfileCommand, executor commandExecutor) (command, error) {
	cmd := []string{SupervisordBinaryPath, SupervisordCtlSubCommand, "stop", "all"}
	if stop, err := newOverriddenSimpleCommand(command, executor, cmd); err == nil {
		// use empty spinner message to avoid showing it altogether
		stop.msg = ""
		return stop, err
	} else {
		return nil, err
	}
}

// newSupervisorStartCommand creates a command implementation that starts the specified command via the supervisor
func newSupervisorStartCommand(command common.DevfileCommand, cmd string, adapter commandExecutor) (command, error) {
	cmdLine := []string{SupervisordBinaryPath, SupervisordCtlSubCommand, "start", cmd}
	return newOverriddenSimpleCommand(command, adapter, cmdLine)
}

func (s supervisorCommand) Execute(show bool) error {
	err := ExecuteCommand(s.adapter, s.info, s.cmd, true, nil, nil)
	if err != nil {
		s.adapter.Logger().ReportError(err, machineoutput.TimestampNow())
		return err
	}
	return nil
}
