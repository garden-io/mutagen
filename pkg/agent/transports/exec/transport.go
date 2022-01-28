package exec

import (
	"errors"
	"os"
	osExec "os/exec"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/agent"
	"github.com/mutagen-io/mutagen/pkg/agent/transport"
	"github.com/mutagen-io/mutagen/pkg/process"
)

// transport implements the agent.Transport interface using SSH.
type execTransport struct {
	// host is the target host.
	command string
	// prompter is the prompter identifier to use for prompting.
	prompter string
}

// NewTransport creates a new exec transport using the specified parameters.
func NewTransport(command string, prompter string) (agent.Transport, error) {
	return &execTransport{
		command:   command,
		prompter:  prompter,
	}, nil
}

// Copy implements the Copy method of agent.Transport.
func (t *execTransport) Copy(localPath, remoteName string) error {
	// This is a no-op. We instead expect the Command handler to start with
	return nil
}

// Command implements the Command method of agent.Transport.
func (t *execTransport) Command(command string) (*osExec.Cmd, error) {
	if (strings.Contains(command, agent.BaseName)) {
		// This is a little hacky, but we need to override the given command with the one specified in the exec:<command>
		// URL when calling the mutagen-agent
		command = t.command
	}

	// Compute the command and args. Note that we ignore the provided command argument!
	split := strings.Split(command, " ")
	name := split[0]

	var args []string
	if len(split) > 1 {
		args = split[1:]
	} else {
		args = []string{}
	}

	// Create the process.
	execCommand := osExec.Command(name, args...)

	// Force it to run detached.
	execCommand.SysProcAttr = transport.ProcessAttributes()

	// Create a copy of the current environment.
	environment := os.Environ()

	// Add locale environment variables.
	// QUESTION: Is this necessary in this context?
	// environment = addLocaleVariables(environment)

	// Set the environment.
	execCommand.Env = environment

	// Done.
	return execCommand, nil
}

// ClassifyError implements the ClassifyError method of agent.Transport.
func (t *execTransport) ClassifyError(processState *os.ProcessState, errorOutput string) (bool, bool, error) {
	// Note: We always return false for the tryInstall value, since installation should never be attempted.

	if process.IsPOSIXShellInvalidCommand(processState) {
		return false, false, nil
	} else if process.IsPOSIXShellCommandNotFound(processState) {
		return false, false, nil
	} else if process.OutputIsWindowsInvalidCommand(errorOutput) {
		// A Windows invalid command error doesn't necessarily indicate that
		// the agent isn't installed, but instead usually indicates that we were
		// trying to invoke the agent using the POSIX shell syntax in a Windows
		// cmd.exe environment. Thus we return false here for re-installation,
		// but we still indicate that this is a Windows platform to potentially
		// change the dialer's platform hypothesis and force it to reconnect
		// under the Windows hypothesis.
		// HACK: We're relying on the fact that the agent dialing logic will
		// attempt a reconnect under the cmd.exe hypothesis, which it will, but
		// this is potentially a bit fragile. We've sort of codified this
		// behavior in the transport interface definition, but it's hard to make
		// super explicit.
		return false, true, nil
	} else if process.OutputIsWindowsCommandNotFound(errorOutput) {
		return false, true, nil
	}

	// Just bail if we weren't able to determine the nature of the error.
	return false, false, errors.New("unknown error condition encountered")
}
