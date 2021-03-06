// ctlcli holds the interfaces and helpers for the current CLI library
// (codeganster/cli).
//
// Abstracting away any cli library implementation from the commands themselves,
// keeping our commands testable and generic.
package ctlcli

import (
	"io"
	"os"

	"github.com/koding/logging"

	"github.com/codegangsta/cli"
)

// Command is the basic interface for all klientctl commands.
type Command interface {
	// Run implements the main CLI function for running a command. Error is returned
	// in combination with the exit status for easy API usage.
	Run() (int, error)

	// Help, presents help to the user. However, many CLI libraries take care of the
	// Help presentation, such as formatting and args. This method usually calls back to
	// the parent CLI library via an internal reference to the ctlcli.Helper type. See
	// ctlcli.Helper docs for further explanation.
	Help()
}

// AutocompleteCommand is an interface for a Command that wants to provide
// Autocomplete functionality.
type AutocompleteCommand interface {
	// Autocomplete prints to Stdout what is to be autocompleted, one item per line.
	// Handling of the autocompletion is done by the shell (bash/fish/etc), this
	// method simply prints to Stdout.
	Autocomplete(args ...string) error
}

// Helper implements an abstraction between Commands and codegansta/cli. How the
// helpers are implemented varies depending on the actual Helper function we're
// wrapping, but typically they write the given io.Writer to the cli.Context
// before calling the cli.ShowHelp-like commands.
type Helper func(io.Writer)

// ExitingCommand is a function that returns an exit code
type ExitingCommand func(*cli.Context, logging.Logger, string) int

// ExitingErrCommand is a function that returns an exit code and an error. Behavior
// is the same as ExitingCommand, but it also supports an error return.
type ExitingErrCommand func(*cli.Context, logging.Logger, string) (int, error)

// CommandFactory returns a struct implementing the Command interface.
type CommandFactory func(*cli.Context, logging.Logger, string) Command

// CommandHelper maps the codegansta/cli Help to our generic Helper type.
// It does so by calling cli.ShowCommandHelper after setting the proper writer. For
// reference, see:
//
// cli.ShowCommandHelp https://github.com/codegangsta/cli/blob/master/help.go#L104
//
// The context and command for this are typically provided by the command factory.
func CommandHelper(ctx *cli.Context, cmd string) Helper {
	return func(w io.Writer) {
		ctx.App.Writer = w
		cli.ShowCommandHelp(ctx, cmd)
	}
}

// ExitAction implements a cli.Command's Action field for an ExitingCommand type.
func ExitAction(f ExitingCommand, log logging.Logger, cmdName string) func(*cli.Context) {
	return func(c *cli.Context) {
		os.Exit(f(c, log, cmdName))
	}
}

// ExitErrAction implements a cli.Command's Action field for an ExitingErrCommand
func ExitErrAction(f ExitingErrCommand, log logging.Logger, cmdName string) func(*cli.Context) {
	return func(c *cli.Context) {
		exit, err := f(c, log, cmdName)
		log.Error("ExitErrAction encountered error. err:%s", err)
		os.Exit(exit)
	}
}

// FactoryAction implements a cli.Command's Action field.
func FactoryAction(factory CommandFactory, log logging.Logger, cmdName string) func(*cli.Context) {
	return func(c *cli.Context) {
		cmd := factory(c, log, cmdName)
		exit, err := cmd.Run()

		// For API reasons, we may return an error but a zero exit code. So we want
		// to check and log both.
		if exit != 0 || err != nil {
			log.Error(
				"Command encountered error. command:%s, exit:%d, err:%s",
				cmdName, exit, err,
			)
		}

		os.Exit(exit)
	}
}

// FactoryCompletion implements codeganstas cli.Command's bash completion field
func FactoryCompletion(factory CommandFactory, log logging.Logger, cmdName string) func(*cli.Context) {
	return func(c *cli.Context) {
		cmd := factory(c, log, cmdName)

		// If the command implements AutocompleteCommand, run the autocomplete.
		if aCmd, ok := cmd.(AutocompleteCommand); ok {
			if err := aCmd.Autocomplete(c.Args()...); err != nil {
				log.Error(
					"Autocompletion of a command encountered error. command:%s, err:%s",
					cmdName, err,
				)
			}
		}
	}
}
