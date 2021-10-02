package option_test

import (
	"strings"
	"testing"
	"time"

	"github.com/maargenton/go-errors"
	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-cli/pkg/option"
)

// ---------------------------------------------------------------------------
// option.NewOptionSet()
// ---------------------------------------------------------------------------

type common struct {
	NonFlag int
	Verbose bool `opts:"-v,--verbose" desc:"verbose description"`
	Debug   bool `opts:"--debug"      desc:"debug description"`
}

type options struct {
	common
	I  int   `opts:"-i, --int, name:i, default:-1"`
	PI *int  `opts:"-p, --pint"`
	IA []int `opts:"-p, --pint, sep:\\,"`
}

func TestOptionSet(t *testing.T) {
	var v options
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).IsNil()
	require.That(t, opts).IsNotNil()
	require.That(t, opts.Options).Length().Eq(5)
}

// -----------
// Error cases

func TestNewOptionSet_NonStruct(t *testing.T) {
	var v *int
	var opts, err = option.NewOptionSet(v)

	require.That(t, err).ToString().Contains("invalid argument of type '*int'")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_InvalidFieldType(t *testing.T) {
	type invalidFieldType struct {
		Ppi **int `opts:"-i"`
	}
	var v invalidFieldType
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("is not parsable")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_InvalidTag(t *testing.T) {
	type invalidTag struct {
		Pi *int `opts:"-i, foobar:none"`
	}
	var v invalidTag
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("invalid tag")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_NestedInvalidTag(t *testing.T) {
	type invalidTag struct {
		Pi *int `opts:"-i, foobar:none"`
	}
	type nestedInvalidTag struct {
		invalidTag
	}
	var v nestedInvalidTag
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("invalid tag")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_UnexportedField(t *testing.T) {
	type invalidFieldType struct {
		pi *int `opts:"-i"`
	}
	var v invalidFieldType
	var opts, err = option.NewOptionSet(&v)
	_ = v.pi

	require.That(t, err).ToString().Contains("is not settable")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_WithEmptyOptsTag(t *testing.T) {
	type invalidFieldType struct {
		Value int `opts:""`
	}
	var v invalidFieldType
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("invalid empty 'opts'")
	require.That(t, opts).IsNil()
}

// -----------------------------
// Test for positional arguments

func TestNewOptionSet_ArgN(t *testing.T) {
	type argN struct {
		Arg1 *string `opts:"arg:1"`
		Arg2 string  `opts:"arg:2"`
		Arg3 *string `opts:"arg:3"`
	}
	var v argN
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).IsNil()
	require.That(t, opts.Options).Length().Eq(0)
	require.That(t, opts.Positional).Length().Eq(3)

	// Only Arg3 should be marked optional
	require.That(t, opts.Positional).Field("Optional").Eq([]bool{
		false, false, true,
	})

	err = opts.ApplyArgs([]string{"aaa", "bbb"})
	require.That(t, err).IsNil()
}

func TestNewOptionSet_ArgN_BadType(t *testing.T) {
	type argN struct {
		Arg1 string   `opts:"arg:1"`
		Arg2 []string `opts:"arg:2"`
		Arg3 string   `opts:"arg:3"`
	}
	var v argN
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).IsNotNil()
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_Arg0(t *testing.T) {
	type arg0 struct {
		Arg1 string `opts:"arg:1"`
		Arg2 string `opts:"arg:2"`
		Arg3 string `opts:"arg:0"`
	}
	var v arg0
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("invalid index '0' for arg:")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_ArgNonInt(t *testing.T) {
	type argNonInt struct {
		Arg1 string `opts:"arg:1"`
		Arg2 string `opts:"arg:nonInt"`
		Arg3 string `opts:"arg:0"`
	}
	var v argNonInt
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("invalid index 'nonInt' for arg:")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_ArgGaps(t *testing.T) {
	type argGaps struct {
		Arg1 string `opts:"arg:1"`
		Arg2 string `opts:"arg:3"`
		Arg3 string `opts:"arg:5"`
	}
	var v argGaps
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("'arg:2' is missing")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_ArgDup(t *testing.T) {
	type argDup struct {
		Arg1 string `opts:"arg:1"`
		Arg2 string `opts:"arg:2"`
		Arg3 string `opts:"arg:1"`
	}
	var v argDup
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("positional argument '1' defined by multiple fields")
	require.That(t, opts).IsNil()
}

// ----------------------
// Test for argument list

func TestNewOptionSet_Args(t *testing.T) {
	type args struct {
		Args []string `opts:"args"`
	}
	var v args
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).IsNil()
	require.That(t, opts.Options).Length().Eq(0)
	require.That(t, opts.Args).IsNotNil()
}

func TestNewOptionSet_ArgsDup(t *testing.T) {
	type argsDup struct {
		Args  []string `opts:"args"`
		Args2 []string `opts:"args"`
	}
	var v argsDup
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).ToString().Contains("multiple fields capturing remaning args")
	require.That(t, opts).IsNil()
}

func TestNewOptionSet_ArgsNotSlice(t *testing.T) {
	type args struct {
		Args string `opts:"args"`
	}
	var v args
	var opts, err = option.NewOptionSet(&v)

	require.That(t, err).IsNotNil()
	require.That(t, opts).IsNil()
}

// ---------------------------------------------------------------------------
// OptionSet.GetOption()
// ---------------------------------------------------------------------------

func TestOption(t *testing.T) {
	args := struct {
		Period time.Duration `opts:"-d,--duration"`
		Name   string        `opts:"-n,--name"`
		Port   string        `opts:"--port"`
	}{}

	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()

	var tcs = []struct {
		optionName string
		fieldName  string
	}{
		{"d", "Period"},
		{"duration", "Period"},
		{"n", "Name"},
		{"name", "Name"},

		{"q", ""},
		{"", ""},
	}

	for _, tc := range tcs {
		t.Run(tc.optionName, func(t *testing.T) {
			opt := optionSet.GetOption(tc.optionName)
			if tc.fieldName != "" {
				require.That(t, opt).IsNotNil()
				require.That(t, opt.FieldName).Eq(tc.fieldName)
			} else {
				require.That(t, opt).IsNil()
			}
		})
	}
}

// ---------------------------------------------------------------------------
// OptionSet.AddSpecialFlag()
// ---------------------------------------------------------------------------
const ErrSpecialFlag = errors.Sentinel("ErrSpecialFlag")
const ErrSpecialFlag2 = errors.Sentinel("ErrSpecialFlag2")

func TestAddSpecialFlag(t *testing.T) {
	type command struct {
		A       bool          `opts:"-a, --aaa"`
		D       time.Duration `opts:"-d, --duration"`
		Hello   bool          `opts:"-h, --hello"`
		Verbose bool          `opts:"-v, --verbose"`
		Version bool          `opts:"-V, --version"`
	}

	t.Run("Given an configured OptionSet", func(t *testing.T) {
		var cmd command
		optionSet, err := option.NewOptionSet(&cmd)
		require.That(t, err).IsNil()

		t.Run("when adding special flag with conflicting short", func(t *testing.T) {
			optionSet.AddSpecialFlag("h", "help", "", ErrSpecialFlag)
			t.Run("then short flag is up-cased", func(t *testing.T) {
				var opt = optionSet.GetOption("help")
				require.That(t, opt).IsNotNil()
				require.That(t, opt.Short).Eq("H")
				require.That(t, opt.SpecialErr).Eq(ErrSpecialFlag)
			})
		})
		t.Run("when adding special flag with conflicting short and up-cased short", func(t *testing.T) {
			optionSet.AddSpecialFlag("v", "version-info", "", ErrSpecialFlag)
			t.Run("then short flag is dropped", func(t *testing.T) {
				require.That(t, optionSet.GetOption("version-info")).IsNotNil()
				require.That(t, optionSet.GetOption("version-info")).Field("Short").Eq("")
			})
		})
		t.Run("when adding special flag with conflicting long", func(t *testing.T) {
			optionSet.AddSpecialFlag("v", "version", "", ErrSpecialFlag)
			t.Run("then the existing flag is preserved", func(t *testing.T) {
				require.That(t, optionSet.GetOption("version")).Field("SpecialErr").IsNil()
				require.That(t, optionSet.GetOption("v")).Field("SpecialErr").IsNil()
				require.That(t, optionSet.GetOption("V")).Field("SpecialErr").IsNil()
			})
		})
	})
}

// ---------------------------------------------------------------------------
// OptionSet.ApplyDefaults()
// ---------------------------------------------------------------------------

func TestApplyDefaults(t *testing.T) {

	args := struct {
		Period time.Duration `opts:"-d,--duration, default:5m"`
	}{}
	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()

	err = optionSet.ApplyDefaults()
	require.That(t, err).IsNil()
	require.That(t, args.Period).Eq(5 * time.Minute)
}

func TestApplyDefaultsError(t *testing.T) {

	args := struct {
		Period time.Duration `opts:"-d,--duration, default:5p"`
	}{}
	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()

	err = optionSet.ApplyDefaults()
	require.That(t, err).IsNotNil()
	require.That(t, err).ToString().Contains("defaults")
	require.That(t, args.Period).Eq(0 * time.Minute)
}

// ---------------------------------------------------------------------------
// OptionSet.ApplyEnv()
// ---------------------------------------------------------------------------

func TestApplyEnv(t *testing.T) {

	args := struct {
		Period time.Duration `opts:"-d,--duration, env:PERIOD"`
	}{}
	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()

	var env = map[string]string{
		"PERIOD": "5m",
	}
	err = optionSet.ApplyEnv(env)
	require.That(t, err).IsNil()
	require.That(t, args.Period).Eq(5 * time.Minute)
}

func TestApplyEnvError(t *testing.T) {

	args := struct {
		Period time.Duration `opts:"-d,--duration, env:PERIOD"`
	}{}
	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()

	var env = map[string]string{
		"PERIOD": "5p",
	}
	err = optionSet.ApplyEnv(env)
	require.That(t, err).IsNotNil()
	require.That(t, err).ToString().Contains("environment variable")
	require.That(t, err).ToString().Contains("PERIOD")
	require.That(t, args.Period).Eq(0 * time.Minute)
}

// ---------------------------------------------------------------------------
// OptionSet.ApplyArgs()
// ---------------------------------------------------------------------------

func TestApplyArgsCombiningBoolFlags(t *testing.T) {
	type command struct {
		A bool   `opts:"-a"`
		B bool   `opts:"-b"`
		C bool   `opts:"-c"`
		F string `opts:"-f"`
	}
	var tcs = []struct {
		args []string
		cmd  command
	}{
		{[]string{"-ac"}, command{A: true, C: true}},
		{[]string{"-a", "-c"}, command{A: true, C: true}},
		{[]string{"-fabc"}, command{F: "abc"}},
		{[]string{"-abfabc"}, command{A: true, B: true, F: "abc"}},
		{[]string{"-afabc", "-b"}, command{A: true, B: true, F: "abc"}},
		{[]string{"-af", "abc", "-b"}, command{A: true, B: true, F: "abc"}},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {

			var cmd command
			optionSet, err := option.NewOptionSet(&cmd)
			require.That(t, err).IsNil()

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, cmd).Eq(tc.cmd)
			require.That(t, err).IsNil()
		})
	}
}

func TestApplyArgsLongFlags(t *testing.T) {
	type command struct {
		A bool   `opts:"-a, --aaa"`
		B bool   `opts:"-b, --bbb"`
		C bool   `opts:"-c, --ccc"`
		F string `opts:"-f, --file"`
	}
	var tcs = []struct {
		args []string
		cmd  command
	}{
		{[]string{"--aaa", "--ccc", "--file", "abc"}, command{A: true, C: true, F: "abc"}},
		{[]string{"--aaa=false", "--ccc", "--file=abc"}, command{A: false, C: true, F: "abc"}},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {
			var cmd command
			optionSet, err := option.NewOptionSet(&cmd)
			require.That(t, err).IsNil()

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, err).IsNil()
			require.That(t, cmd).Eq(tc.cmd)
		})
	}
}

func TestApplyArgs_Delimiter(t *testing.T) {
	type command struct {
		A bool     `opts:"-a, --aaa"`
		D string   `opts:"-d, --duration"`
		L []string `opts:"args"`
	}

	bdd.Given(t, "a struct with args field", func(t *bdd.T) {
		var cmd command
		optionSet, err := option.NewOptionSet(&cmd)
		require.That(t, err).IsNil()

		t.When("arguments contain `--`", func(t *bdd.T) {
			var args = []string{"--duration", "1m", "--", "--aaa"}

			t.Then("remaining flags are not parsed", func(t *bdd.T) {
				err = optionSet.ApplyArgs(args)
				require.That(t, err).IsNil()
				require.That(t, cmd.A).IsFalse()
				require.That(t, cmd.L).Eq([]string{"--aaa"})
			})
		})

		t.When("option value is `--`", func(t *bdd.T) {
			var args = []string{"--duration", "--", "1m", "--aaa"}

			t.Then("remaining flags are parsed", func(t *bdd.T) {
				err = optionSet.ApplyArgs(args)
				require.That(t, err).IsNil()
				require.That(t, cmd.A).IsTrue()
				require.That(t, cmd.D).Eq("--")
				require.That(t, cmd.L).Eq([]string{"1m"})
			})
		})
	})
}

func TestApplyArgs_Errors(t *testing.T) {
	type command struct {
		A bool          `opts:"-a, --aaa"`
		D time.Duration `opts:"-d, --duration"`
	}
	var tcs = []struct {
		args []string
		err  string
	}{
		{[]string{"--aaa", "--ddd"}, "invalid flag"},
		{[]string{"-abgc"}, "invalid flag"},
		{[]string{"-d4p"}, "invalid value"},
		{[]string{"-d", "4p"}, "invalid value"},
		{[]string{"--duration", "4p"}, "invalid value"},
		{[]string{"--duration=4p"}, "invalid value"},
		{[]string{"a", "--duration"}, "missing argument"},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {
			var cmd command
			optionSet, err := option.NewOptionSet(&cmd)
			require.That(t, err).IsNil()

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, err).IsNotNil()
			require.That(t, err).ToString().Contains(tc.err)
		})
	}
}

func TestApplyArgs_WithSpecialFlags(t *testing.T) {
	type command struct {
		A bool          `opts:"-a, --aaa"`
		D time.Duration `opts:"-d, --duration"`
	}
	var tcs = []struct {
		args []string
		err  error
	}{
		{[]string{"--aaa", "--help"}, ErrSpecialFlag},
		{[]string{"--aaa", "-h"}, ErrSpecialFlag},
		{[]string{"--aaa", "--version"}, ErrSpecialFlag2},
		{[]string{"--aaa", "-v"}, ErrSpecialFlag2},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {
			var cmd command
			optionSet, err := option.NewOptionSet(&cmd)
			require.That(t, err).IsNil()

			optionSet.AddSpecialFlag("h", "help", "", ErrSpecialFlag)
			optionSet.AddSpecialFlag("v", "version", "", ErrSpecialFlag2)

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, err).IsError(tc.err)
		})
	}
}

func TestApplyArgs_NonFlagArguments(t *testing.T) {
	type command struct {
		A bool          `opts:"-a, --aaa"`
		D time.Duration `opts:"-d, --duration"`
		U string        `opts:"arg:1"`
		V string        `opts:"arg:2"`
		W []string      `opts:"args"`
	}
	var tcs = []struct {
		args []string
		cmd  command
	}{
		{
			[]string{"--aaa", "--duration", "5m", "aaa", "bbb", "ccc", "ddd"},
			command{
				A: true,
				D: 5 * time.Minute,
				U: "aaa",
				V: "bbb",
				W: []string{"ccc", "ddd"},
			},
		},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {
			var cmd command
			optionSet, err := option.NewOptionSet(&cmd)
			require.That(t, err).IsNil()

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, err).IsNil()
			require.That(t, cmd).Eq(tc.cmd)
		})
	}
}

func TestApplyArgs_NonFlagArguments_Errors(t *testing.T) {
	var tcs = []struct {
		args []string
		cmd  interface{}
		err  string
	}{
		{
			[]string{"--aaa", "aaa"},
			&struct {
				A bool          `opts:"-a, --aaa"`
				D time.Duration `opts:"-d, --duration"`
				U string        `opts:"arg:1"`
				V string        `opts:"arg:2"`
				W []string      `opts:"args"`
			}{},
			"missing argument",
		},
		{
			[]string{"--aaa", "aaa", "bbb", "ccc"},
			&struct {
				A bool          `opts:"-a, --aaa"`
				D time.Duration `opts:"-d, --duration"`
				U string        `opts:"arg:1"`
				V string        `opts:"arg:2"`
				// W []string      `opts:"args"`
			}{},
			"unsupported extra arguments",
		},
		{
			[]string{"--aaa", "aaa", "bbb", "ccc"},
			&struct {
				A bool          `opts:"-a, --aaa"`
				D time.Duration `opts:"-d, --duration"`
				U int           `opts:"arg:1"`
				V int           `opts:"arg:2"`
				W []int         `opts:"args"`
			}{},
			"failed to set value",
		},
		{
			[]string{"--aaa", "1", "2", "ccc"},
			&struct {
				A bool          `opts:"-a, --aaa"`
				D time.Duration `opts:"-d, --duration"`
				U int           `opts:"arg:1"`
				V int           `opts:"arg:2"`
				W []int         `opts:"args"`
			}{},
			"failed to set value",
		},
	}

	for _, tc := range tcs {
		var name = strings.Join(tc.args, " ")
		t.Run(name, func(t *testing.T) {

			optionSet, err := option.NewOptionSet(tc.cmd)
			require.That(t, err).IsNil()

			err = optionSet.ApplyArgs(tc.args)
			require.That(t, err).ToString().Contains(tc.err)
		})
	}
}
