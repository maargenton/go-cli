package cli_test

import (
	"testing"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-cli/pkg/cli"
	"github.com/maargenton/go-cli/pkg/option"
)

func TestBashCompletionScript(t *testing.T) {
	var name = "command-name"
	var script = cli.BashCompletionScript(name)

	require.That(t, script).Contains("_command-name_completion()")
	require.That(t, script).Contains("complete -F _command-name_completion command-name")
}

func TestDefaultCompletion(t *testing.T) {
	bdd.Given(t, "the current directory structure", func(t *bdd.T) {
		t.When("calling DefaultCompletion() with an empty string", func(t *bdd.T) {
			suggestions := cli.DefaultCompletion(nil, "")
			t.Then("suggestions include the local files", func(t *bdd.T) {
				require.That(t, suggestions).IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})
		t.When("calling DefaultCompletion() with partial filename", func(t *bdd.T) {
			suggestions := cli.DefaultCompletion(nil, "comp")
			t.Then("suggestions include only the matching files", func(t *bdd.T) {
				require.That(t, suggestions).IsEqualSet(
					[]string{"completion.go", "completion_test.go"})
			})
		})
		t.When("calling DefaultCompletion() with partial unique folder name", func(t *bdd.T) {
			suggestions := cli.DefaultCompletion(nil, "../cl")
			t.Then("suggestions include the files in that folder", func(t *bdd.T) {
				require.That(t, suggestions).IsSupersetOf([]string{
					"../cli/completion.go",
					"../cli/completion_test.go",
				})
			})
		})
	})
}

func TestMatchingFilenameCompletion(t *testing.T) {
	bdd.Given(t, "a call to MatchingFilenameCompletion()", func(t *bdd.T) {
		t.When("passing a pattern and an empty string", func(t *bdd.T) {
			suggestions := cli.MatchingFilenameCompletion(nil, "*_test.go", "")

			t.Then("suggestions include all filenames matching the pattern", func(t *bdd.T) {
				require.That(t, suggestions).IsEqualSet(
					[]string{"cmd_test.go", "completion_test.go"})
			})
		})
		t.When("passing a pattern and a partial name", func(t *bdd.T) {
			suggestions := cli.MatchingFilenameCompletion(nil, "*_test.go", "co")

			t.Then("suggestions include only filenames matching both", func(t *bdd.T) {
				require.That(t, suggestions).IsEqualSet(
					[]string{"completion_test.go"})
			})
		})
		t.When("passing a pattern and a non-matching partial name", func(t *bdd.T) {
			suggestions := cli.MatchingFilenameCompletion(nil, "*_test.go", "er")

			t.Then("the pattern is ignored and all matching files are returned", func(t *bdd.T) {
				require.That(t, suggestions).IsEqualSet(
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

func (c *compCmd2) Complete(opt *option.T, partial string) []string {
	if opt.Long == "option" {
		return []string{"aaa", "bbb", "ccc"}
	}
	if opt.Args {
		return []string{"ddd", "eee", "fff"}
	}
	return nil
}

func TestCommandRunCompletion(t *testing.T) {
	bdd.Given(t, "a well defined command struct", func(t *bdd.T) {
		var cmd = &cli.Command{
			Handler:     &compCmd{},
			Description: "command description",
		}
		var c = cmd.Handler.(*compCmd)

		t.When("calling Run() with completion request and partial option flag", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"-v", "--o"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "--o",
				"COMP_INDEX": "2",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Then("the command is not run", func(t *bdd.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Then("the completion request error is returned", func(t *bdd.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Then("the suggestions contain the matching flag", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).Length().Eq(1)
				require.That(t, cmd.Suggestions[0]).StartsWith("--option")
			})
		})

		t.When("calling Run() with a partial option argument", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"-v", "--option", "co"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "co",
				"COMP_INDEX": "3",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Then("the command is not run", func(t *bdd.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Then("the completion request error is returned", func(t *bdd.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Then("the suggestions contain matching local filenames", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})

		t.When("calling Run() with nothing", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"command-name"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "1",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Then("the command is not run", func(t *bdd.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Then("the completion request error is returned", func(t *bdd.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Then("the suggestions include option flags", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsSupersetOf(
					[]string{"--verbose", "--option"})
			})
			t.Then("the suggestions include argument options", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsSupersetOf(
					[]string{"completion.go", "completion_test.go"})
			})
		})
	})

	bdd.Given(t, "a command with custom completion handler", func(t *bdd.T) {
		var cmd = &cli.Command{
			Handler:     &compCmd2{},
			Description: "command description",
		}
		var c = cmd.Handler.(*compCmd2)
		_ = c

		t.When("calling Run() with nothing", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"command-name"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "1",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Then("the command is not run", func(t *bdd.T) {
				require.That(t, c.didRun).IsFalse()
			})
			t.Then("the completion request error is returned", func(t *bdd.T) {
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
			t.Then("the suggestions include option flags", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsSupersetOf(
					[]string{"--verbose", "--option"})
			})
			t.Then("the suggestions include argument options", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsSupersetOf(
					[]string{"ddd", "eee", "fff"})
			})
		})

		t.When("calling Run() with missing option argument", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"-v", "--option"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "",
				"COMP_INDEX": "3",
			}
			cmd.Suggestions = nil
			cmd.Run()

			t.Then("the suggestions contain matching local filenames", func(t *bdd.T) {
				require.That(t, cmd.Suggestions).IsEqualSet(
					[]string{"aaa", "bbb", "ccc"})
			})
		})
	})
}

func TestCommandRunCompletionDebuf(t *testing.T) {

	bdd.Given(t, "an environment with COMPLETION_DEBUG_OUTPUT set", func(t *bdd.T) {
		var cmd = &cli.Command{Handler: &compCmd{}}
		var tmp = t.TempDir()
		var filename = fileutils.Join(tmp, "comp.json")
		t.Setenv("COMPLETION_DEBUG_OUTPUT", filename)

		t.When("calling Run() with completion request", func(t *bdd.T) {
			cmd.ProcessArgs = []string{"-v", "--o"}
			cmd.ProcessEnv = map[string]string{
				"COMP_WORD":  "--o",
				"COMP_INDEX": "2",
			}
			cmd.Suggestions = nil
			err := cmd.Run()

			t.Then("the specified file is created", func(t *bdd.T) {
				require.That(t, fileutils.Exists(filename)).IsTrue()
				require.That(t, err).IsError(cli.ErrCompletionRequested)
			})
		})
	})
}
