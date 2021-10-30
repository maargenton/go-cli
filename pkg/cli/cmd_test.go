package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-cli/pkg/cli"
)

type myCmd struct {
	Verbose bool `opts:"-v,--verbose"`
	Arg     int  `opts:"-a, --arg, env:TEST_ARG"`
	didRun  bool
	err     error
	version string
}

func (c *myCmd) Run() error {
	c.didRun = true
	return c.err
}

func (c *myCmd) Version() string {
	return c.version
}

func TestSetProcessEnv(t *testing.T) {
	t.Run("Given an environment as a list of strings", func(t *testing.T) {
		var env = []string{
			"KEY=VALUE",
			"EMPTY_KEY=",
			"_=1",
		}
		t.Run("when calling cmd.SetProcessEnv", func(t *testing.T) {
			var cmd cli.Command
			cmd.SetProcessEnv(env)

			t.Run("then all environment values are recorded", func(t *testing.T) {
				require.That(t, cmd.ProcessEnv).MapKeys().IsEqualSet([]string{
					"KEY", "EMPTY_KEY", "_",
				})
				require.That(t, cmd.ProcessEnv).Field("KEY").Eq("VALUE")
				require.That(t, cmd.ProcessEnv).Field("_").Eq("1")

				v, ok := cmd.ProcessEnv["EMPTY_KEY"]
				require.That(t, v).Eq("")
				require.That(t, ok).IsTrue()

			})
		})
	})
}

func TestCommandRun(t *testing.T) {
	t.Run("Given a well defined command struct", func(t *testing.T) {
		var cmd = &cli.Command{
			Handler:     &myCmd{},
			Description: "command description",
		}
		var c = cmd.Handler.(*myCmd)

		t.Run("when calling run with valid arguments", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--arg", "123"}
			err := cmd.Run()

			t.Run("then the fields are set and the command is run", func(t *testing.T) {
				require.That(t, err).IsNil()
				require.That(t, c.Verbose).IsTrue()
				require.That(t, c.Arg).Eq(123)
				require.That(t, c.didRun).IsTrue()
			})
		})

		t.Run("when the command handler returns an error", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--arg", "123"}
			c.err = fmt.Errorf("myError")
			err := cmd.Run()

			t.Run("then the error is returned", func(t *testing.T) {
				require.That(t, err).IsError(c.err)
			})
		})

		t.Run("when calling run with invalid arguments", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--args", "123"}
			err := cmd.Run()

			t.Run("then an error is returned", func(t *testing.T) {
				require.That(t, err).IsNotNil()
			})
		})
		t.Run("when calling run with invalid environment value", func(t *testing.T) {
			cmd.ProcessArgs = []string{"command-name", "-v", "--arg", "123"}
			cmd.ProcessEnv = map[string]string{
				"TEST_ARG": "argument",
			}
			err := cmd.Run()

			t.Run("then an error is returned", func(t *testing.T) {
				require.That(t, err).IsNotNil()
				require.That(t, err).ToString().Contains("environment variable 'TEST_ARG'")
			})
		})
	})

	t.Run("Given a command with invalid default", func(t *testing.T) {
		type myCmd2 struct {
			myCmd
			Arg int `opts:"-a, --arg, default:bad"`
		}
		var cmd = &cli.Command{
			Handler:     &myCmd2{},
			Description: "command description",
		}
		var c = cmd.Handler.(*myCmd2)

		t.Run("when calling run with valid arguments", func(t *testing.T) {
			err := cmd.Run()

			t.Run("then an error is returned", func(t *testing.T) {
				require.That(t, err).IsNotNil()
				require.That(t, err).ToString().Contains("while applying defaults")
				require.That(t, c.didRun).IsFalse()
			})
		})
	})

	t.Run("Given a command with bad option", func(t *testing.T) {
		type myCmd2 struct {
			myCmd
			Arg int `opts:"-a, --arg, bad:option"`
		}
		var cmd = &cli.Command{
			Handler:     &myCmd2{},
			Description: "command description",
		}
		var c = cmd.Handler.(*myCmd2)

		t.Run("when calling run with valid arguments", func(t *testing.T) {
			err := cmd.Run()

			t.Run("then an error is returned", func(t *testing.T) {
				require.That(t, err).IsNotNil()
				require.That(t, err).ToString().Contains("invalid tag in opts: 'bad'")
				require.That(t, c.didRun).IsFalse()
			})
		})
	})
}

func splitLines(s string) []string {
	var lines = strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func TestCommandUsage(t *testing.T) {
	t.Run("Given a well defined command struct", func(t *testing.T) {
		type myCmd2 struct {
			myCmd
			Port  string   `opts:"arg:1, name:port"     desc:"port to open"`
			Port2 *string  `opts:"arg:2, name:aux-port" desc:"auxiliary port"`
			Ports []string `opts:"args"                 desc:"additional ports to open"`
		}

		var cmd = &cli.Command{
			Handler:     &myCmd2{},
			Description: "command description",
		}

		t.Run("when calling Usage()", func(t *testing.T) {
			cmd.ProcessName = "command-name"
			cmd.ConsoleWidth = 80
			usage := splitLines(cmd.Usage())

			t.Run("then a command description is returned", func(t *testing.T) {
				require.That(t, usage).Length().Eq(7)
				require.That(t, usage[0]).Contains("command-name")
				require.That(t, usage[1]).Contains("command description")
				require.That(t, usage[2]).Contains("<port>")
				require.That(t, usage[3]).Contains("<aux-port>")
				require.That(t, usage[4]).Contains("<value>")
				require.That(t, usage[5]).Contains("--verbose")
				require.That(t, usage[6]).Contains("TEST_ARG")
			})
		})
	})

	t.Run("Given an invalid command struct", func(t *testing.T) {
		type myCmd2 struct {
			myCmd
			Value int `opts:"-v, --value, bad-tag:bad"`
		}

		var cmd = &cli.Command{
			Handler:     &myCmd2{},
			Description: "command description",
		}

		t.Run("when calling Usage()", func(t *testing.T) {
			cmd.ProcessName = "command-name"
			cmd.ConsoleWidth = 80

			usage := cmd.Usage()

			t.Run("then the error is printed in place of the usage", func(t *testing.T) {
				require.That(t, usage).Contains("error initializing the command for Usage")
			})
		})
	})
}

type customUsageCmd struct {
	myCmd
}

func (c *customUsageCmd) Usage(name string, width int) string {
	return "custom usage"
}

func TestCustomUsage(t *testing.T) {
	t.Run("Given a command struct with custom usage", func(t *testing.T) {
		var cmd = &cli.Command{
			Handler:     &customUsageCmd{},
			Description: "command description",
		}

		t.Run("when calling Usage()", func(t *testing.T) {
			cmd.ProcessName = "command-name"
			cmd.ConsoleWidth = 80

			usage := cmd.Usage()

			t.Run("then the error is printed in place of the usage", func(t *testing.T) {
				require.That(t, usage).Contains("custom usage")
			})
		})
	})
}

func TestCommandVersion(t *testing.T) {
	t.Run("Given a command struct with version handler", func(t *testing.T) {
		var cmd = &cli.Command{
			Handler:     &myCmd{},
			Description: "command description",
		}
		var c = cmd.Handler.(*myCmd)

		t.Run("when calling Usage()", func(t *testing.T) {
			c.version = "v1.0.0"
			version := cmd.Version()

			t.Run("then a command description is returned", func(t *testing.T) {
				require.That(t, version).Eq("v1.0.0")
			})
		})
	})
}
