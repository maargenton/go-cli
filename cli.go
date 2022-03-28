package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/maargenton/go-fileutils"
	"golang.org/x/term"

	"github.com/maargenton/go-cli/pkg/cli"
	"github.com/maargenton/go-cli/pkg/option"
)

// Command is the main public type used to define all the details of a comand to
// be handled.
type Command = cli.Command

// Suggestion is an alias of `option.Description` used to describe one
// suggestions in the context of a completion request.
type Suggestion = option.Description

// SimpleSuggestionList generate a list of Suggestion objects with no
// description from a list of strings
func SimpleSuggestionList(options ...string) []Suggestion {
	var r = make([]option.Description, 0, len(options))
	for _, o := range options {
		r = append(r, option.Description{Option: o})
	}
	return r
}

// DefaultCompletion acts like the shell default completion and suggests file
// and folder names under the current directory. It is used by default when the
// command does not implement a specific completion handler, and should be used
// from the command completion handler when no other completion logic is
// suitable.
var DefaultCompletion = cli.DefaultCompletion

// FilepathCompletion implements a custom filepath completion scheme, matching
// the provided pattern if possible.
var FilepathCompletion = cli.FilepathCompletion

// Run takes the command line arguments, parses them and execute the
// command or sub-command with the corresponding options.
func Run(cmd *Command) {
	if cmd.ProcessName == "" {
		cmd.ProcessName = fileutils.Base(os.Args[0])
	}
	if cmd.ProcessArgs == nil {
		cmd.ProcessArgs = os.Args
	}
	cmd.ConsoleWidth = consoleWidth()
	cmd.SetProcessEnv(os.Environ())

	var err = cmd.Run()
	if errors.Is(err, cli.ErrCompletionScriptRequested) {
		fmt.Print(cli.BashCompletionScript(cmd.ProcessName))

	} else if errors.Is(err, cli.ErrHelpRequested) {
		fmt.Print(cmd.Usage())

	} else if errors.Is(err, cli.ErrVersionRequested) {
		var version = cmd.Version()
		if version != "" {
			fmt.Printf("%v\n", version)
		}
	} else if errors.Is(err, cli.ErrCompletionRequested) {
		fmt.Printf("%v",
			cli.FormatCompletionSuggestions(cmd.ConsoleWidth, cmd.Suggestions),
		)
		// cli.DumpCompletionSuggestions("competion.txt", cmd.Suggestions)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func consoleWidth() int {
	var width = 80
	if ww, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = ww
	}
	return width
}
