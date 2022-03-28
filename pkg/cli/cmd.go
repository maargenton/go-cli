package cli

import (
	"fmt"
	"strings"

	"github.com/maargenton/go-cli/pkg/option"
)

// ---------------------------------------------------------------------------
// Command type definition
// ---------------------------------------------------------------------------

// Command is the representation of a runnable command, with reference to a
// runnable command attached to an options struct
type Command struct {
	Handler     Handler
	Description string

	ProcessName       string
	ProcessArgs       []string
	ProcessEnv        map[string]string
	ConsoleWidth      int
	DisableCompletion bool

	Suggestions []option.Description

	opts *option.Set
}

// Handler defines the interface necessary to run a command once the command
// lien arguments have been parsed
type Handler interface {
	Run() error
}

// VersionHandler defines an optional interface for the command handler to
// return a relevant version string
type VersionHandler interface {
	Version() string
}

// UsageHandler defines an optional interface for the command handler to
// print a custom usage message
type UsageHandler interface {
	Usage(name string, width int) string
}

// CompletionHandler is an optional interface for the command handler to provide
// meaningful values for a specific option or argument.
type CompletionHandler interface {
	Complete(opt *option.T, partial string) []option.Description
}

// ---------------------------------------------------------------------------
// Command type public interface
// ---------------------------------------------------------------------------

// SetProcessEnv sets the command `ProcessEnv` from a list of environment
// strings as returned by os.Environ().
func (cmd *Command) SetProcessEnv(env []string) {
	var ee = make(map[string]string, len(env))
	for _, e := range env {
		var i = strings.IndexByte(e, '=')
		if i >= 0 {
			var key = e[:i]
			var value = e[i+1:]
			ee[key] = value
		}
	}
	cmd.ProcessEnv = ee
}

// Run is the main invocation point for a command. The command must be seeded
// with all the necessary runtime arguments from the process context
// (`ProcessName`, `ProcessArgs`, `ProcessEnv` and `ConsoleWidth`). It sets up
// the command option struct, applies the defaults abd environment variable,
// decodes the command-line and run the command.
func (cmd *Command) Run() error {

	if err := cmd.initialize(); err != nil {
		return err
	}

	if _, ok := cmd.Handler.(VersionHandler); ok {
		cmd.opts.AddSpecialFlag(
			"v", "version", "display version information",
			ErrVersionRequested)
	}

	cmd.opts.AddSpecialFlag(
		"h", "help", "display usage information",
		ErrHelpRequested)

	if !cmd.DisableCompletion {
		cmd.opts.AddSpecialFlag(
			"", "bash-completion-script",
			"generate a bash script that sets up completion for this command; "+
				"to use, run the following line or add it to your .bash_profile:\n"+
				"eval $("+cmd.ProcessName+" --bash-completion-script)",
			ErrCompletionScriptRequested)

		if cmd.handleCompletionRequest() {
			return ErrCompletionRequested
		}
	}
	if err := cmd.opts.ApplyDefaults(); err != nil {
		return err
	}
	if err := cmd.opts.ApplyEnv(cmd.ProcessEnv); err != nil {
		return err
	}
	if err := cmd.opts.ApplyArgs(cmd.ProcessArgs[1:]); err != nil {
		return err
	}

	if err := cmd.Handler.Run(); err != nil {
		return err
	}

	return nil
}

// Usage returns a string containign the usage for the command. The display name
// for the command is expected as first argument.
func (cmd *Command) Usage() string {

	if err := cmd.initialize(); err != nil {
		return fmt.Sprintf("error initializing the command for Usage:  %v", err)
	}

	if uh, ok := cmd.Handler.(UsageHandler); ok {
		return uh.Usage(cmd.ProcessName, cmd.ConsoleWidth)
	}

	var args []string
	for _, opt := range cmd.opts.Positional {
		var name = opt.Name()
		if opt.Optional {
			args = append(args, fmt.Sprintf("[%v]", name))
		} else {
			args = append(args, fmt.Sprintf("%v", name))
		}
	}
	if cmd.opts.Args != nil {
		args = append(args, cmd.opts.Args.Name())
	}

	var usage strings.Builder
	fmt.Fprintf(&usage,
		"Usage: %v [options] %v\n",
		cmd.ProcessName, strings.Join(args, " "))
	fmt.Fprintf(&usage, "%v\n", cmd.Description)

	var options []option.Description
	for _, arg := range cmd.opts.Positional {
		var usage = arg.GetUsage()
		if usage.Description != "" {
			options = append(options, usage)
		}
	}
	if arg := cmd.opts.Args; arg != nil {
		var usage = arg.GetUsage()
		if usage.Description != "" {
			options = append(options, usage)
		}
	}
	for _, opt := range cmd.opts.Options {
		options = append(options, opt.GetUsage())
	}
	fmt.Fprint(&usage, option.FormatOptionDescription("  ", cmd.ConsoleWidth, options))

	return usage.String()
}

// Version returns a version string for the command
func (cmd *Command) Version() (version string) {
	if vh, ok := cmd.Handler.(VersionHandler); ok {
		version = vh.Version()
	}
	return
}

// ---------------------------------------------------------------------------
// Command type private implementation
// ---------------------------------------------------------------------------

// initialize parses the tags of the handler struct and records all the
// available options. The function is safe to call more than once.
func (cmd *Command) initialize() error {
	if cmd.opts == nil {
		var opts, err = option.NewOptionSet(cmd.Handler)
		if err != nil {
			return err
		}
		cmd.opts = opts
	}
	return nil
}
