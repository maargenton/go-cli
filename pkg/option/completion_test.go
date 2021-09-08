package option_test

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-cli/pkg/option"
)

func TestGetCompletion(t *testing.T) {

	type command struct {
		Port     *string `opts:"arg:1, name:port"           desc:"name of the port to open"`
		Baudrate *uint32 `opts:"-b, --baudrate, name:speed" desc:"baudrate to use for communication"`
		Format   *string `opts:"-f, --format"               desc:"communcation format, e.g. 8N1"`

		Timestamp bool     `opts:"-t,--timestamp" desc:"prefix every line with elapsed time"`
		Verbose   bool     `opts:"-v"             desc:"display additional information on startup"`
		Extra     []string `opts:"args"`
	}

	t.Run("Given an configured OptionSet", func(t *testing.T) {
		var cmd command
		optionSet, err := option.NewOptionSet(&cmd)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		optionSet.AddSpecialFlag(
			"v", "version", "display version information",
			ErrSpecialFlag,
		)

		t.Run("when calling GetCompletion() with no arguments", func(t *testing.T) {
			var args = []string{}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then no specific option is being completed", func(t *testing.T) {
				require.That(t, completion.Opt).IsNil()
				require.That(t, completion.OptValues).IsEmpty()
			})
			t.Run("then the first argument is being completed", func(t *testing.T) {
				require.That(t, completion.Arg).Eq(optionSet.Positional[0])
			})
			t.Run("then available options are listed", func(t *testing.T) {
				require.That(t, completion.Options).Length().Eq(5)
			})
		})

		t.Run("when calling GetCompletion() with flag expecting a value", func(t *testing.T) {
			var args = []string{"--baudrate"}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then only yhe option is returned", func(t *testing.T) {
				var expected = optionSet.GetOption("baudrate")
				require.That(t, completion.Opt).Eq(expected)
				require.That(t, completion.Arg).IsNil()
				require.That(t, completion.Options).IsEmpty()
			})
		})

		t.Run("when calling GetCompletion() with flag and value", func(t *testing.T) {
			var args = []string{"--baudrate", "115200", "--timestamp"}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then remaining options are listed", func(t *testing.T) {
				require.That(t, completion.Opt).IsNil()
				require.That(t, completion.Options).Field("Option").IsEqualSet(
					[]string{"--format <value>", "-v"})
			})
			t.Run("then the first argument is being completed", func(t *testing.T) {
				require.That(t, completion.Arg).Eq(optionSet.Positional[0])
			})
		})

		t.Run("when calling GetCompletion() with short flag and value", func(t *testing.T) {
			var args = []string{"-b115200", "-t"}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then remaining options are listed", func(t *testing.T) {
				require.That(t, completion.Opt).IsNil()
				require.That(t, completion.Options).Field("Option").IsEqualSet(
					[]string{"--format <value>", "-v"})
			})
			t.Run("then the first argument is being completed", func(t *testing.T) {
				require.That(t, completion.Arg).Eq(optionSet.Positional[0])
			})
		})

		t.Run("when calling GetCompletion() with exclusive flag", func(t *testing.T) {
			var args = []string{"--version"}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then remaining options are listed", func(t *testing.T) {
				require.That(t, completion.Opt).IsNil()
				require.That(t, completion.Arg).IsNil()
				require.That(t, completion.Options).IsEmpty()
			})
		})

		t.Run("when calling GetCompletion() with argument and flags", func(t *testing.T) {
			var args = []string{"/dev/tty.usb-1210", "-b", "115200"}
			var partial = ""
			completion := optionSet.GetCompletion(args, partial)

			t.Run("then remaining options are listed", func(t *testing.T) {
				require.That(t, completion.Opt).IsNil()
				require.That(t, completion.Options).Field("Option").IsEqualSet(
					[]string{"--format <value>", "-v", "--timestamp"})
			})
			t.Run("then next argument is being completed", func(t *testing.T) {
				require.That(t, completion.Arg).Eq(optionSet.Args)
			})
		})
	})
}
