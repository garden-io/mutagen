package exec

import (
	"errors"
	"os"
	osExec "os/exec"
	"runtime"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/agent"
	"github.com/mutagen-io/mutagen/pkg/agent/transport"
	"github.com/mutagen-io/mutagen/pkg/process"
)

// execTransport implements the agent.Transport interface using a shell command.
type execTransport struct {
	// host is the target host.
	command string
	// prompter is the prompter identifier to use for prompting.
	prompter string
}

// NewTransport creates a new exec transport using the specified parameters.
func NewTransport(command string, prompter string) (agent.Transport, error) {
	return &execTransport{
		command:  command,
		prompter: prompter,
	}, nil
}

// Copy implements the Copy method of agent.Transport.
func (t *execTransport) Copy(localPath, remoteName string) error {
	// This is a no-op. In this custom implementation, we expect the agent to be
	// pre-installed on the remote.
	return nil
}

// Command implements the Command method of agent.Transport.
func (t *execTransport) Command(command string) (*osExec.Cmd, error) {
	if strings.Contains(command, agent.BaseName) {
		// This is a little hacky, but we need to override the given command with the one specified in the exec:<command>
		// URL when calling the mutagen-agent
		command = t.command
	} else {
		// Since this transport only works on the primary agent dialing path, we
		// disallow all other commands.
		return nil, errors.New("transport does not support this command")
	}

	// Compute the command and arguments. Note that we're ignoring the original
	// command in this case for the one passed in via the exec URL. Unlike a
	// standard Mutagen-issued command, we can't rely on the ability to lex the
	// command by string splitting, because the prefix for kubectl can contain
	// spaces (e.g. in a home directory name). Thus, we split off the kubectl
	// target manually with some assumptions about how that command will be
	// structured.
	var split []string
	var firstSplitIndicator string
	if runtime.GOOS == "windows" {
		firstSplitIndicator = "\\kubectl.exe "
	} else {
		firstSplitIndicator = "/kubectl "
	}
	if index := strings.Index(command, firstSplitIndicator); index >= 0 {
		split = append(split, command[:index+len(firstSplitIndicator)-1])
		remaining := command[index+len(firstSplitIndicator):]
		split = append(split, strings.Split(remaining, " ")...)
	} else {
		return nil, errors.New("unable to identify kubectl invocation")
	}
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
