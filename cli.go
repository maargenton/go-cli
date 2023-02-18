package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/maargenton/go-fileutils"
	"golang.org/x/term"

	"github.com/maargenton/go-cli/pkg/cli"
)

// Command is the main public type used to define all the details of a command
// to be handled.
type Command = cli.Command

// DefaultCompletion acts like the shell default completion and suggests file
// and folder names under the current directory. It is used by default when the
// command does not implement a specific completion handler, and should be used
// from the command completion handler when no other completion logic is
// suitable.
var DefaultCompletion = cli.DefaultCompletion

// DefaultFilenameCompletion implements a default filename-based completion,
// similar to the default shell completion.
var DefaultFilenameCompletion = cli.DefaultFilenameCompletion

// MatchingFileCompletion implements a custom filepath completion scheme,
// matching the provided pattern if possible, using default filename completion
// as a fallback.
var MatchingFilenameCompletion = cli.MatchingFilenameCompletion

// Run takes the command line arguments, parses them and execute the
// command or sub-command with the corresponding options.
func Run(cmd *Command) {
	if cmd.ProcessName == "" {
		cmd.ProcessName = fileutils.Base(os.Args[0])
	}
	if cmd.ProcessArgs == nil {
		cmd.ProcessArgs = os.Args[1:]
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
		for _, v := range cmd.Suggestions {
			fmt.Println(v)
		}
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
