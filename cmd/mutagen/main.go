package main

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/fatih/color"

	"github.com/mutagen-io/mutagen/cmd"
	"github.com/mutagen-io/mutagen/cmd/mutagen/compose"
	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"
	"github.com/mutagen-io/mutagen/cmd/mutagen/forward"
	"github.com/mutagen-io/mutagen/cmd/mutagen/project"
	"github.com/mutagen-io/mutagen/cmd/mutagen/root"
	"github.com/mutagen-io/mutagen/cmd/mutagen/sync"

	"github.com/mutagen-io/mutagen/pkg/prompting"
)

func init() {
	// Disable alphabetical sorting of commands in help output. This is a global
	// setting that affects all Cobra command instances.
	cobra.EnableCommandSorting = false

	// Disable Cobra's use of mousetrap. This breaks daemon registration on
	// Windows because it tries to enforce that the CLI only be launched from
	// a console, which it's not when running automatically.
	cobra.MousetrapHelpText = ""

	// Register commands.
	root.RootCommand.AddCommand(
		sync.SyncCommand,
		forward.ForwardCommand,
		project.ProjectCommand,
		compose.RootCommand,
		daemon.DaemonCommand,
		versionCommand,
		legalCommand,
		generateCommand,
	)

	// HACK If we're on Windows, enable color support for command usage and
	// error output by recursively replacing the output streams for Cobra
	// commands.
	if runtime.GOOS == "windows" {
		enableColorForCommand(root.RootCommand)
	}
}

// enableColorForCommand recursively enables colorized usage and error output
// for a command and all of its child commands.
func enableColorForCommand(command *cobra.Command) {
	// Enable color support for the command itself.
	command.SetOut(color.Output)
	command.SetErr(color.Error)

	// Recursively enable color support for child commands.
	for _, c := range command.Commands() {
		enableColorForCommand(c)
	}
}

func main() {
	// Check if a prompting environment is set. If so, treat this as a prompt
	// request. Prompting is sort of a special pseudo-command that's indicated
	// by the presence of an environment variable, and hence it has to be
	// handled in a bit of a special manner.
	if _, ok := os.LookupEnv(prompting.PrompterEnvironmentVariable); ok {
		if err := promptMain(os.Args[1:]); err != nil {
			cmd.Fatal(err)
		}
		return
	}

	// Handle terminal compatibility issues. If this call returns, it means that
	// we should proceed normally.
	cmd.HandleTerminalCompatibility()

	// HACK: If we're performing command line completion, then remove the
	// adapter command that we use to keep the Docker Compose command hierarchy
	// separate and replace it with the actual Docker Compose command hierarchy
	// so that completions work properly.
	if cmd.PerformingShellCompletion {
		root.RootCommand.RemoveCommand(compose.RootCommand)
		root.RootCommand.AddCommand(compose.ComposeCommand)
	}

	// Execute the root command.
	if err := root.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
