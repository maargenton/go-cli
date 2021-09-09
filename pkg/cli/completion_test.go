package cli_test

import (
	"path/filepath"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-cli/pkg/cli"
	"github.com/maargenton/go-cli/pkg/option"
)

func TestBashCompletionScript(t *testing.T) {
	var name = "command-name"
	var script = cli.BashCompletionScript(name)

	require.That(t, script).Contains("_command-name_completion()")
	require.That(t, script).Contains("complete -F _command-name_completion command-name")
}

func TestFormatCompletionSuggestions(t *testing.T) {

	t.Run("Given a list of suggestions", func(t *testing.T) {
		suggestions := []option.Description{
			{Option: "foo <value>", Description: "foo description"},
			{Option: "bar <value>", Description: "bar description"},
		}
		t.Run("when calling FormatCompletionSuggestions()", func(t *testing.T) {
			width := 40
			lines := splitLines(cli.FormatCompletionSuggestions(width, suggestions))

			t.Run("then all lines are at most the width", func(t *testing.T) {
				verify.That(t, lines).Length().Eq(2)
				verify.That(t, lines).All(
					subexpr.Value().Length().Le(width))
			})
		})
		t.Run("when calling FormatCompletionSuggestions() with a single option", func(t *testing.T) {
			width := 40
			lines := splitLines(cli.FormatCompletionSuggestions(width, suggestions[:1]))

			t.Run("then only the first word of the option is printed", func(t *testing.T) {
				require.That(t, lines).Length().Eq(1)
				verify.That(t, lines[0]).Eq("foo")
			})
		})
	})
}

func TestDefaultCompletion(t *testing.T) {
	t.Run("Given the current directory structure", func(t *testing.T) {
		t.Run("when Calling DefaultCompletion() with an empty string", func(t *testing.T) {
			suggestions := cli.DefaultCompletion("")
			t.Run("then suggestions include the local files", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})
		t.Run("when Calling DefaultCompletion() with partial filename", func(t *testing.T) {
			suggestions := cli.DefaultCompletion("comp")
			t.Run("then suggestions include only the matching files", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsEqualSet(
					[]string{"completion.go", "completion_test.go"})
			})
		})
		t.Run("when Calling DefaultCompletion() with partial unique folder name", func(t *testing.T) {
			suggestions := cli.DefaultCompletion(filepath.FromSlash("../cl"))
			t.Run("then suggestions include the files in that folder", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsSupersetOf([]string{
					filepath.FromSlash("../cli/completion.go"),
					filepath.FromSlash("../cli/completion_test.go"),
				})
			})
		})
	})
}

func TestFilepathCompletion(t *testing.T) {
	t.Run("Given a call to FilepathCompletion()", func(t *testing.T) {
		t.Run("when passing a pattern and an empty string", func(t *testing.T) {
			suggestions := cli.FilepathCompletion("*_test.go", "")

			t.Run("then suggestions include all filenames matching the pattern", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsEqualSet(
					[]string{"cmd_test.go", "completion_test.go"})
			})
		})
		t.Run("when passing a pattern and a partial name", func(t *testing.T) {
			suggestions := cli.FilepathCompletion("*_test.go", "co")

			t.Run("then suggestions include only filenames matching both", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsEqualSet(
					[]string{"completion_test.go"})
			})
		})
		t.Run("when passing a pattern and a non-matching partial name", func(t *testing.T) {
			suggestions := cli.FilepathCompletion("*_test.go", "er")

			t.Run("then the pattern is ignored and all matching files are returned", func(t *testing.T) {
				require.That(t, suggestions).Field("Option").IsEqualSet(
					[]string{"errors.go"})
			})
		})
	})
}

// ---------------------------------------------------------------------------

type compCmd struct {
	Verbose bool     `opts:"-v,--verbose"`
	Option  string   `opts:"-o, --option"`
	Args    []string `opts:"args"`

	didRun bool
}

func (c *compCmd) Run() error {
	c.didRun = true
	return nil
}

type compCmd2 struct {
	compCmd
}

func (c *compCmd2) Complete(opt *option.T, partial string) []option.Description {
	if opt.Long == "option" {
		return []option.Description{
			{Option: "aaa"}, {Option: "bbb"}, {Option: "ccc"},
		}
	}
	if opt.Args {
		return []option.Description{
			{Option: "ddd"}, {Option: "eee"}, {Option: "fff"},
		}
	}
	return nil
}

func TestCommandRunCompletion(t *testing.T) {
	t.Run("Given a well defined command struct", func(t *testing.T) {
		var cmd = &cli.Command{
			Handler:     &compCmd{},
			Description: "command description",
		}
		var c = cmd.Handler.(*compCmd)

		t.Run("when calling Run() with completion request and partial option flag", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--o"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "--o",
				"COMP_INDEX": "2",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Run("then the command is not run", func(t *testing.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Run("then the completion request error is returned", func(t *testing.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Run("then the suggestions contain the matching flag", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Length().Eq(1)
				require.That(t, cmd.Suggestions[0].Option).StartsWith("--option")
			})
		})

		t.Run("when calling Run() with a partial option argument", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--option", "co"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "co",
				"COMP_INDEX": "3",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Run("then the command is not run", func(t *testing.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Run("then the completion request error is returned", func(t *testing.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Run("then the suggestions contain matching local filenames", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})

		t.Run("when calling Run() with nothing", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "1",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Run("then the command is not run", func(t *testing.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Run("then the completion request error is returned", func(t *testing.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Run("then the suggestions include option flags", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsSupersetOf(
					[]string{"--verbose", "--option <value>"})
			})
			t.Run("then the suggestions include argument options", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})
	})

	t.Run("Given a command with custom completion handler", func(t *testing.T) {
		var cmd = &cli.Command{
			Handler:     &compCmd2{},
			Description: "command description",
		}
		var c = cmd.Handler.(*compCmd2)
		_ = c

		t.Run("when calling Run() with nothing", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "1",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Run("then the command is not run", func(t *testing.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Run("then the completion request error is returned", func(t *testing.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Run("then the suggestions include option flags", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsSupersetOf(
					[]string{"--verbose", "--option <value>"})
			})
			t.Run("then the suggestions include argument options", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsSupersetOf(
					[]string{"ddd", "eee", "fff"})
			})
		})

		t.Run("when calling Run() with missing option argument", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--option"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "3",
			}
			cmd.Suggestions = nil
			cmd.Run()

			t.Run("then the suggestions contain matching local filenames", func(t *testing.T) {
				require.That(t, cmd.Suggestions).Field("Option").IsEqualSet(
					[]string{"aaa", "bbb", "ccc"})
			})
		})
	})

}
