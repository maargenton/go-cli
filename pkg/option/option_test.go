package option_test

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/maargenton/go-errors"
	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"

	"github.com/maargenton/go-cli/pkg/option"
	"github.com/maargenton/go-cli/pkg/value"
)

// ---------------------------------------------------------------------------
// Option.Name()
// ---------------------------------------------------------------------------

func TestOptionName(t *testing.T) {
	bdd.Given(t, "a struct with opts tags", func(t *bdd.T) {
		type options struct {
			Period time.Duration `opts:"-d,--duration"`
			Count  int           `opts:"-c"`
			Name   string        `opts:"arg:1, name:name"`
			Name2  string        `opts:"arg:2"`
			Inputs []string      `opts:"args, name:inputs"`
		}

		var o = &options{}
		optionSet, err := option.NewOptionSet(o)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		t.When("getting the name of a field with long flag", func(t *bdd.T) {
			var option = optionSet.GetOption("duration")
			require.That(t, option).IsNotNil()

			t.Then("the name includes the long flag", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("--duration")
			})
		})

		t.When("getting the name of a field with only short flag", func(t *bdd.T) {
			var option = optionSet.GetOption("c")
			require.That(t, option).IsNotNil()

			t.Then("the name includes the short flag", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("-c")
			})
		})

		t.When("getting the name of a positional argument field", func(t *bdd.T) {
			var option = optionSet.Positional[0]
			require.That(t, option).IsNotNil()

			t.Then("the name includes the value name", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("<name>")
			})
		})
		t.When("getting the name of a positional argument field", func(t *bdd.T) {
			var option = optionSet.Positional[1]
			require.That(t, option).IsNotNil()

			t.Then("the name includes the argument position", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("<arg2>")
			})
		})
		t.When("getting the name of an extra arguments field", func(t *bdd.T) {
			var option = optionSet.Args
			require.That(t, option).IsNotNil()

			t.Then("the name includes the value name", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("<inputs>...")
			})
		})
	})

	bdd.Given(t, "another struct with opts tags", func(t *bdd.T) {
		type options struct {
			Inputs []string `opts:"args"`
		}

		var o = &options{}
		optionSet, err := option.NewOptionSet(o)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		t.When("getting the name of an extra arguments field", func(t *bdd.T) {
			var option = optionSet.Args
			require.That(t, option).IsNotNil()

			t.Then("the name includes default name and ellipsis", func(t *bdd.T) {
				require.That(t, option.Name()).Eq("<args>...")
			})
		})

	})
}

// ---------------------------------------------------------------------------
// Option.Usage()
// ---------------------------------------------------------------------------

func TestOptionDescriptionUsage(t *testing.T) {
	var tcs = []struct {
		name  string
		opt   option.T
		usage string
	}{
		{
			name: "an Option with short and long",
			opt: option.T{
				Short: "p",
				Long:  "port",
			},
			usage: "-p, --port <value>",
		},
		{
			name: "an Option with short only",
			opt: option.T{
				Short: "p",
			},
			usage: "-p <value>",
		},
		{
			name: "an Option with long only",
			opt: option.T{
				Long: "port",
			},
			usage: "    --port <value>",
		},
		{
			name: "an Option with named value",
			opt: option.T{
				Short:     "p",
				Long:      "port",
				ValueName: "port",
			},
			usage: "-p, --port <port>",
		},
		{
			name: "a positional Option",
			opt: option.T{
				ValueName: "port",
				Position:  1,
			},
			usage: "<port>",
		},
		{
			name: "a positional Option with no name",
			opt: option.T{
				Position: 2,
			},
			usage: "<arg2>",
		},
		{
			name: "a remaining arguments Option",
			opt: option.T{
				ValueName: "ports",
				Args:      true,
			},
			usage: "<ports>...",
		},
		{
			name: "a remaining arguments Option with no name",
			opt: option.T{
				Args: true,
			},
			usage: "<args>...",
		},
		{
			name: "a boolean Option",
			opt: option.T{
				Short: "d",
				Long:  "debug",
				Type:  option.Bool,
			},
			usage: "-d, --debug",
		},
		{
			name: "a special Option",
			opt: option.T{
				Short:      "v",
				Long:       "version",
				Type:       option.Special,
				SpecialErr: errors.Sentinel("ErrDisplayVersion"),
			},
			usage: "-v, --version",
		},
	}

	for _, tc := range tcs {
		bdd.Given(t, tc.name, func(t *bdd.T) {
			t.When("calling Usage()", func(t *bdd.T) {
				t.Then("it returns formatted usage", func(t *bdd.T) {
					require.That(t, tc.opt.GetUsage()).Field("Option").Eq(tc.usage)
				})
			})
		})
	}
}

func TestOptionDescriptionDescription(t *testing.T) {
	var tcs = []struct {
		name string
		opt  option.T
		desc string
	}{
		{
			name: "an Option{} with no description",
			opt:  option.T{},
			desc: "",
		},
		{
			name: "an Option{} with description",
			opt: option.T{
				Description: "description",
			},
			desc: "description",
		},
		{
			name: "an Option{} with default and env",
			opt: option.T{
				Description: "description",
				Default:     "default",
				Env:         "ENV",
			},
			desc: "description, default: default, env: ENV",
		},
		{
			name: "an Option{} with only default and env",
			opt: option.T{
				Default: "default",
				Env:     "ENV",
			},
			desc: "default: default, env: ENV",
		},
	}

	for _, tc := range tcs {
		bdd.Given(t, tc.name, func(t *bdd.T) {
			t.When("calling Usage()", func(t *bdd.T) {
				t.Then("it returns formatted description", func(t *bdd.T) {
					require.That(t, tc.opt.GetUsage()).Field("Description").Eq(tc.desc)
				})
			})
		})
	}
}

// ---------------------------------------------------------------------------
// Option.SetBool()
// ---------------------------------------------------------------------------

func TestOption_SetBool(t *testing.T) {
	bdd.Given(t, "a struct with bool options", func(t *bdd.T) {
		args := struct {
			A   bool  `opts:"-a,--aaa"`
			Ptr *bool `opts:"-p,--ppp"`
			N   int   `opts:"-n,--number"`
		}{}
		optionSet, err := option.NewOptionSet(&args)
		require.That(t, err).IsNil()

		t.When("calling Set() on a bool field", func(t *bdd.T) {
			optionSet.GetOption("a").SetBool()

			t.Then("the field is set to true", func(t *bdd.T) {
				require.That(t, args.A).IsTrue()
			})
		})

		t.When("calling Set() on a bool pointer field", func(t *bdd.T) {
			optionSet.GetOption("p").SetBool()

			t.Then("the field is set to point to a true value", func(t *bdd.T) {
				require.That(t, args.Ptr).IsNotNil()
				require.That(t, *args.Ptr).IsTrue()
			})
		})

		t.When("calling Set() on a non-bool field", func(t *bdd.T) {
			t.Then("it panics", func(t *bdd.T) {

				require.That(t, func() {
					optionSet.GetOption("n").SetBool()
				}).PanicsAndRecoveredValue().Eq(
					"cannot call Option.SetBool() on non-bool fields",
				)
			})
		})
	})
}

// ---------------------------------------------------------------------------
// Option.SetValue()
// ---------------------------------------------------------------------------

func TestOption_SetValue(t *testing.T) {
	bdd.Given(t, "a struct with `opts` tags", func(t *bdd.T) {

		args := struct {
			Period    time.Duration   `opts:"-p,--period"`
			Duration  *time.Duration  `opts:"-d,--duration"`
			Intervals []time.Duration `opts:"-i,--intervals,sep:\\,"`
		}{}
		optionSet, err := option.NewOptionSet(&args)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		t.When("setting the value of a scalar field", func(t *bdd.T) {
			option := optionSet.GetOption("period")
			require.That(t, option).IsNotNil()

			t.Run("with a valid value", func(t *bdd.T) {
				args.Period = 0 * time.Second
				err := option.SetValue("5m30s")
				expected := 5*time.Minute + 30*time.Second

				t.Then("the field is set accordingly", func(t *bdd.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Period).Eq(expected)
				})
			})

			t.Run("with an invalid value", func(t *bdd.T) {
				args.Period = 0 * time.Second
				err := option.SetValue("5m30p")

				t.Then("an error is returned and the value is not changed", func(t *bdd.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Period).Eq(0 * time.Second)
				})
			})
		})

		t.When("setting the value of a pointer field", func(t *bdd.T) {
			option := optionSet.GetOption("duration")
			require.That(t, option).IsNotNil()

			t.Run("with a valid value", func(t *bdd.T) {
				args.Duration = nil
				err := option.SetValue("5m30s")
				expected := 5*time.Minute + 30*time.Second

				t.Then("the field is set accordingly", func(t *bdd.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Duration).IsNotNil()
					require.That(t, *args.Duration).Eq(expected)
				})
			})

			t.Run("with an invalid value", func(t *bdd.T) {
				args.Duration = nil
				err := option.SetValue("5m30p")

				t.Then("an error is returned and the value is not changed", func(t *bdd.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Duration).IsNil()
				})
			})
		})

		t.When("setting the value of a slice field", func(t *bdd.T) {
			option := optionSet.GetOption("intervals")
			require.That(t, option).IsNotNil()

			t.Run("with a single valid value", func(t *bdd.T) {
				args.Intervals = nil
				err := option.SetValue("5m30s")
				expected := []time.Duration{
					5*time.Minute + 30*time.Second,
				}

				t.Then("the field is set accordingly", func(t *bdd.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Intervals).IsNotNil()
					require.That(t, args.Intervals).Length().Eq(1)
					require.That(t, args.Intervals).Eq(expected)
				})
			})

			t.Run("with multiple delimited values", func(t *bdd.T) {
				args.Intervals = nil
				err1 := option.SetValue("1m,2m,3m")
				err2 := option.SetValue("4m,5m")
				expected := []time.Duration{
					1 * time.Minute,
					2 * time.Minute,
					3 * time.Minute,
					4 * time.Minute,
					5 * time.Minute,
				}

				t.Then("all values are recorded", func(t *bdd.T) {
					require.That(t, err1).IsNil()
					require.That(t, err2).IsNil()
					require.That(t, args.Intervals).Eq(expected)
				})
			})

			t.Run("with empty value", func(t *bdd.T) {
				args.Intervals = []time.Duration{
					1 * time.Minute,
					2 * time.Minute,
					3 * time.Minute,
					4 * time.Minute,
					5 * time.Minute,
				}
				err := option.SetValue("")

				t.Then("all previous values are deleted", func(t *bdd.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Intervals).IsEmpty()
				})
			})

			t.Run("with invalid value", func(t *bdd.T) {
				args.Intervals = nil
				err := option.SetValue("1m,2p,3m")

				t.Then("all values are recorded", func(t *bdd.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Intervals).IsEmpty()
				})
			})
		})
	})
}

// ---------------------------------------------------------------------------
// Option.SetValue() -- parsable pointer type
// ---------------------------------------------------------------------------

func TestOption_SetValue_Ptr(t *testing.T) {
	value.RegisterParser(url.Parse)

	args := struct {
		URL *url.URL `opts:"--url"`
	}{}
	optionSet, err := option.NewOptionSet(&args)
	require.That(t, err).IsNil()
	require.That(t, optionSet).IsNotNil()

	urlType := reflect.TypeOf(args.URL)
	opt := optionSet.Options[0]

	// Because *url.URL is directly parsable because of the registered parser,
	// both FieldType and ValueType should be *url.URL, and the field should be
	// handled as a regular value.
	verify.That(t, opt.FieldType).Eq(urlType)
	verify.That(t, opt.ValueType).Eq(urlType)
	verify.That(t, opt.Type).Eq(option.Value)
}

// ---------------------------------------------------------------------------
// Option.SetValue() -- slice type
// ---------------------------------------------------------------------------

func TestOption_SetValue_Slices(t *testing.T) {
	bdd.Given(t, "an option set with a string-slice field", func(t *bdd.T) {
		args := struct {
			Default    []string `opts:"-v, --values, sep:\\,, env:DEFAULT"`
			KeepSpaces []string `opts:"-v, --values, sep:\\,, keep-spaces, env:KEEP_SPACES"`
			KeepEmpty  []string `opts:"-v, --values, sep:\\,, keep-empty, env:KEEP_EMPTY"`
			KeepBoth   []string `opts:"-v, --values, sep:\\,, keep-spaces,keep-empty, env:KEEP_BOTH"`
		}{}
		opts, err := option.NewOptionSet(&args)
		require.That(t, err).IsNil()
		require.That(t, opts).IsNotNil()

		t.When("using field with default options (trim spaces, discard empty)", func(t *bdd.T) {
			opt := opts.Options[0]
			opt.SetValue("aaa, bbb, ,")

			t.Then("spaces are trimmed from the individual values", func(t *bdd.T) {
				verify.That(t, args.Default).Length().Eq(2)
				verify.That(t, args.Default).Eq([]string{"aaa", "bbb"})
			})
		})
		t.When("using field with keep-spaces", func(t *bdd.T) {
			opt := opts.Options[1]
			opt.SetValue("aaa, bbb, ")

			t.Then("spaces are preserved for the individual values", func(t *bdd.T) {
				verify.That(t, args.KeepSpaces).Length().Eq(3)
				verify.That(t, args.KeepSpaces).Eq([]string{"aaa", " bbb", " "})
			})
		})
		t.When("using field with keep-empty", func(t *bdd.T) {
			opt := opts.Options[2]
			opt.SetValue("aaa, bbb,,")

			t.Then("spaces are trimmed but empty values are preserved", func(t *bdd.T) {
				verify.That(t, args.KeepEmpty).Length().Eq(3)
				verify.That(t, args.KeepEmpty).Eq([]string{"aaa", "bbb", ""})
			})
		})
		t.When("using field with keep-spaces and keep-empty", func(t *bdd.T) {
			opt := opts.Options[3]
			opt.SetValue("aaa, bbb,, ")

			t.Then("all values are preserved with their spaces", func(t *bdd.T) {
				verify.That(t, args.KeepBoth).Length().Eq(4)
				verify.That(t, args.KeepBoth).Eq([]string{"aaa", " bbb", "", " "})
			})
		})
	})
}
