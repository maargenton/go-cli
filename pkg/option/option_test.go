package option_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/require"

	"github.com/maargenton/go-cli/pkg/option"
)

// ---------------------------------------------------------------------------
// Option.Name()
// ---------------------------------------------------------------------------

func TestOptionName(t *testing.T) {
	t.Run("Given a struct with opts tags", func(t *testing.T) {
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

		t.Run("when getting the name of a field with long flag", func(t *testing.T) {
			var option = optionSet.GetOption("duration")
			require.That(t, option).IsNotNil()

			t.Run("the name includes the long flag", func(t *testing.T) {
				require.That(t, option.Name()).Eq("--duration")
			})
		})

		t.Run("when getting the name of a field with only short flag", func(t *testing.T) {
			var option = optionSet.GetOption("c")
			require.That(t, option).IsNotNil()

			t.Run("the name includes the short flag", func(t *testing.T) {
				require.That(t, option.Name()).Eq("-c")
			})
		})

		t.Run("when getting the name of a positional argument field", func(t *testing.T) {
			var option = optionSet.Positional[0]
			require.That(t, option).IsNotNil()

			t.Run("the name includes the value name", func(t *testing.T) {
				require.That(t, option.Name()).Eq("<name>")
			})
		})
		t.Run("when getting the name of a positional argument field", func(t *testing.T) {
			var option = optionSet.Positional[1]
			require.That(t, option).IsNotNil()

			t.Run("the name includes the argument position", func(t *testing.T) {
				require.That(t, option.Name()).Eq("<arg[2]>")
			})
		})
		t.Run("when getting the name of an extra arguments field", func(t *testing.T) {
			var option = optionSet.Args
			require.That(t, option).IsNotNil()

			t.Run("the name includes the value name", func(t *testing.T) {
				require.That(t, option.Name()).Eq("<inputs>...")
			})
		})
	})

	t.Run("Given another struct with opts tags", func(t *testing.T) {
		type options struct {
			Inputs []string `opts:"args"`
		}

		var o = &options{}
		optionSet, err := option.NewOptionSet(o)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		t.Run("when getting the name of an extra arguments field", func(t *testing.T) {
			var option = optionSet.Args
			require.That(t, option).IsNotNil()

			t.Run("the name includes default name and ellipsis", func(t *testing.T) {
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
			name: "an Option{} with short and long",
			opt: option.T{
				Short: "p",
				Long:  "port",
			},
			usage: "-p, --port <value>",
		},
		{
			name: "an Option{} with short only",
			opt: option.T{
				Short: "p",
			},
			usage: "-p <value>",
		},
		{
			name: "an Option{} with long only",
			opt: option.T{
				Long: "port",
			},
			usage: "    --port <value>",
		},
		{
			name: "an Option{} with named value",
			opt: option.T{
				Short:     "p",
				Long:      "port",
				ValueName: "port",
			},
			usage: "-p, --port <port>",
		},
		{
			name: "an positional Option{} with nameonly",
			opt: option.T{
				ValueName: "port",
			},
			usage: "<port>",
		},
		{
			name: "an positional Option{} with nameonly",
			opt: option.T{
				ValueName: "port",
			},
			usage: "<port>",
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("Given %v", tc.name), func(t *testing.T) {
			t.Run("when calling Usage()", func(t *testing.T) {
				t.Run("then it returns formatted usage", func(t *testing.T) {
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
		t.Run(fmt.Sprintf("Given %v", tc.name), func(t *testing.T) {
			t.Run("when calling Usage()", func(t *testing.T) {
				t.Run("then it returns formatted description", func(t *testing.T) {
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
	t.Run("Given a struct with bool options", func(t *testing.T) {
		args := struct {
			A   bool  `opts:"-a,--aaa"`
			Ptr *bool `opts:"-p,--ppp"`
			N   int   `opts:"-n,--number"`
		}{}
		optionSet, err := option.NewOptionSet(&args)
		require.That(t, err).IsNil()

		t.Run("when calling Set() on a bool field", func(t *testing.T) {
			optionSet.GetOption("a").SetBool()

			t.Run("then the field is set to true", func(t *testing.T) {
				require.That(t, args.A).IsTrue()
			})
		})

		t.Run("when calling Set() on a bool pointer field", func(t *testing.T) {
			optionSet.GetOption("p").SetBool()

			t.Run("then the field is set to point to a true value", func(t *testing.T) {
				require.That(t, args.Ptr).IsNotNil()
				require.That(t, *args.Ptr).IsTrue()
			})
		})

		t.Run("when calling Set() on a non-bool field", func(t *testing.T) {
			t.Run("then it panics", func(t *testing.T) {

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
	t.Run("Given a struct with `opts` tags", func(t *testing.T) {

		args := struct {
			Period    time.Duration   `opts:"-p,--period"`
			Duration  *time.Duration  `opts:"-d,--duration"`
			Intervals []time.Duration `opts:"-i,--intervals,sep:\\,"`
		}{}
		optionSet, err := option.NewOptionSet(&args)
		require.That(t, err).IsNil()
		require.That(t, optionSet).IsNotNil()

		t.Run("when setting the value of a scalar field", func(t *testing.T) {
			option := optionSet.GetOption("period")
			require.That(t, option).IsNotNil()

			t.Run("with a valid value", func(t *testing.T) {
				args.Period = 0 * time.Second
				err := option.SetValue("5m30s")
				expected := 5*time.Minute + 30*time.Second

				t.Run("then the field is set accordingly", func(t *testing.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Period).Eq(expected)
				})
			})

			t.Run("with an invalid value", func(t *testing.T) {
				args.Period = 0 * time.Second
				err := option.SetValue("5m30p")

				t.Run("then an error is returned and the value is not changed", func(t *testing.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Period).Eq(0 * time.Second)
				})
			})
		})

		t.Run("when setting the value of a pointer field", func(t *testing.T) {
			option := optionSet.GetOption("duration")
			require.That(t, option).IsNotNil()

			t.Run("with a valid value", func(t *testing.T) {
				args.Duration = nil
				err := option.SetValue("5m30s")
				expected := 5*time.Minute + 30*time.Second

				t.Run("then the field is set accordingly", func(t *testing.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Duration).IsNotNil()
					require.That(t, *args.Duration).Eq(expected)
				})
			})

			t.Run("with an invalid value", func(t *testing.T) {
				args.Duration = nil
				err := option.SetValue("5m30p")

				t.Run("then an error is returned and the value is not changed", func(t *testing.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Duration).IsNil()
				})
			})
		})

		t.Run("when setting the value of a slice field", func(t *testing.T) {
			option := optionSet.GetOption("intervals")
			require.That(t, option).IsNotNil()

			t.Run("with a single valid value", func(t *testing.T) {
				args.Intervals = nil
				err := option.SetValue("5m30s")
				expected := []time.Duration{
					5*time.Minute + 30*time.Second,
				}

				t.Run("then the field is set accordingly", func(t *testing.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Intervals).IsNotNil()
					require.That(t, args.Intervals).Length().Eq(1)
					require.That(t, args.Intervals).Eq(expected)
				})
			})

			t.Run("with multiple delimited values", func(t *testing.T) {
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

				t.Run("then all values are recorded", func(t *testing.T) {
					require.That(t, err1).IsNil()
					require.That(t, err2).IsNil()
					require.That(t, args.Intervals).Eq(expected)
				})
			})

			t.Run("with empty value", func(t *testing.T) {
				args.Intervals = []time.Duration{
					1 * time.Minute,
					2 * time.Minute,
					3 * time.Minute,
					4 * time.Minute,
					5 * time.Minute,
				}
				err := option.SetValue("")

				t.Run("then all previous values are deleted", func(t *testing.T) {
					require.That(t, err).IsNil()
					require.That(t, args.Intervals).IsEmpty()
				})
			})

			t.Run("with invalid value", func(t *testing.T) {
				args.Intervals = nil
				err := option.SetValue("1m,2p,3m")

				t.Run("then all values are recorded", func(t *testing.T) {
					require.That(t, err).IsNotNil()
					require.That(t, args.Intervals).IsEmpty()
				})
			})
		})
	})
}
