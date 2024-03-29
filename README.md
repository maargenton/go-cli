# go-cli

A library to define command line interfaces in Go.

[![Latest](
  https://img.shields.io/github/v/tag/maargenton/go-cli?color=blue&label=latest&logo=go&logoColor=white&sort=semver)](
  https://pkg.go.dev/github.com/maargenton/go-cli)
[![Build](
  https://img.shields.io/github/actions/workflow/status/maargenton/go-cli/build.yaml?branch=master&label=build&logo=github&logoColor=aaaaaa)](
  https://github.com/maargenton/go-cli/actions?query=branch%3Amaster)
[![Codecov](
  https://img.shields.io/codecov/c/github/maargenton/go-cli?label=codecov&logo=codecov&logoColor=aaaaaa&token=1JFLIu042X)](
  https://codecov.io/gh/maargenton/go-cli)
[![Go Report Card](
  https://goreportcard.com/badge/github.com/maargenton/go-cli)](
  https://goreportcard.com/report/github.com/maargenton/go-cli)


---------------------------

Package `go-cli` provides a declarative way to define full featured command-line
interfaces. It follows the POSIX/GNU-style guidelines and supports custom bash
completion. The current version support only simple commands with flags. Future
versions will support composite commands with sub-commands, like `go` and `git`.

The options supported by a command are defined by the fields of a `struct` with
specific struct-tags, and the command itself is defined by a `Run()` method
attached to that struct.


## Features

The implementation of go-clo follows the conventions outlined in [The Open Group
    Base Specifications Issue 7, 2018 edition, Chapter 12. Utility
    Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)

### Short flags support

If a command defines `-c`, `-v` and `-z` as boolean flags and `-f` as a string
option, then:
- `-cvz` is equivalent to `-c -v -z`
- `-cvzf foo` is equivalent to `-c -v -z -f foo`
- `-cvzffoo` is equivalent to `-c -v -z -f foo`
- `-cffoovz` is equivalent to `-c -f foovz`

### Long flags support

Long flags always start with a `--` and can accept value either inline after an
`=` sign, or as the next argument.

- `--filename foo`
- `--filename=foo`
- `--filename=` specifies an empty filename

Boolean long flags do not accept a value unless attached with an `=` sign;
`--bool-flag=false` is the only way to set a boolean flag with a default value
of true back to false.

### Positional and addition arguments

All non-option command-line arguments must be captured by a field in the command
options struct, otherwise an error is generated. Non-option arguments can be
captures as either positional arguments or additional arguments; additional
arguments must be be backed by a slice type.

A special delimiter `--` marks the end of option flags and capture the remaining
arguments as non-option, assigned to positional and additional arguments fields.
Note that after parsing, it is not possible to determine if arguments were
specified before or after a `--` delimiter.

### Limitation

- Each option flag can appear only once unless it is backed by an slice type.
- The order in which arguments are specified on the command line cannot be
  retrieved after parsing (except for arguments backed by a slice type in which
  values are stored in order).
- All option and non-option arguments are handled independently of their
  position on the command-line, so you cannot have different flag values
  attached to different arguments. For example, the following hypothetical
  compiler command cannot be handled by go-cli: `cc -O2 foo.cpp -O0 bar.ccp`.
- `-vvvv` to e.g. increase verbosity to level 4 is ***not supported***.


## Installation

```bash
go get github.com/maargenton/go-cli
```

## Usage

```go
type cmd struct {
    Port      *string  `opts:"arg:1, name:port" desc:"..."`
    Format    *string  `opts:"-f, --format"     desc:"..."`
    Verbose   bool     `opts:"-v, --verbose"    desc:"..."`
    List      []string `opts:"-l, --list, sep:\\, , env:LIST"`
}

func (options *cmd) Run() error {
    // This function is invoked once the options struct has been
    // initialized with default values, values from the environment
    // and values from the command-line.
    return nil;
}

func main() {
    cli.Run(&cli.Command{
        Handler:     &cmd{},
        Description: "...",
    })
}
```

The `cmd` type is a struct that will capture all the options and arguments
passed on the command line. The main function defines a command referring to the
`cmd` type, and through `cli.Run()` parses the command-line arguments and
environment variable to configure the command before calling `cmd.Run()`.

### Struct tags

The details of the command-line interface are defined on the struct with `opts`
struct tags. Each tag contains a comma-separated list of a short name and/or
long name and additional options in the form `key:value`. Within the value,
reserved characters (comma and colon) can be escaped with a double-backslash
`\\`. For example, to accept a comma-separated list of values, a option should
include `sep:\\,`.

The `opts` tag can consist of:
- `-b,--baudrate`: either or both of a short and long flag name for the option
- `arg:<n>` : captures a positional argument
- `args` : captures all remaining arguments
- `default:` : a default value for the field if not specified on the
  command-line
- `env:` : the name of an environment variable that can override the default
- `sep:`: a list separator characters for fields that can accept multiple
  values. To use a comma or colon as separator, those characters must be escaped
  with a double-backslash (`\\,` or `\\:`). To use spaces and newlines as
  separator, the special value `\\s` can be used. Unless a separator is
  specified with this option, additional values must be specified by repeating
  the option flag multiple times.
- `keep-spaces` : when specified with no value, this option preserves the spaces
  around the argument values, that would otherwise be trimmed by default.
- `keep-empty` : for fields thar accept multiple values using a separator, this
  option preserve empty values after splitting and trimming, that would
  otherwise be dropped by default.
- `name:`: the display name for the value, used when printing out description of
  the field.

A separate `desc` struct tag contains the description for the option.

### Supported field types

Fields in the command options struct can be:
- any parsable value type
- a pointer to a parsable type
- a slice of parsable value type.

A parsable types is a type whose value that can be set from a string using some
standard interface, and includes:
- `bool` type
- All built-in integer and float types,
- `string` type
- `time.Duration`
- Any type that conforms to `flag.Value`.
- Any type that conforms to `encoding.TextUnmarshaler`
- Any other type that registers a custom parser function through
  `value.RegisterParser()`

Unless otherwise initialized, all pointer fields are initialized to `nil`, all
slice fields are initialized to an empty slice, and all scalar fields are
initialized to their built-in zero value. If the command needs to differentiate
between the built-in zero value of a scalar field and a specified zero value, a
pointer type should be used and checked against `nil`.

By default, slice field values must be provided by repeating the option flag
multiple times, once of each value. If a separator is defined (`sep:`), multiple
or all values can be provided with one command-line argument. To provide multiple values through an environment variable, a separated must be defined.

> New in v0.5.0: A breaking change has been introduced to better handle lists of
> values and spaces around values. Prior behavior can be restores with the
> `keep-spaces` option for all fields and `keep-empty` for lists. With this new
> behavior, spaces around the values are automatically trimmed, and empty values
> after trimming are dropped from value lists. In effect, this allows for
> trailing separators to be ignored, and for list values to be specified with
> additional spaces around the separators, including newlines.  Note that, even
> with `keep-empty`, if the last character is a separator, the last empty values
> is always drop, as was the case in prior versions.

In a command-line interface, all option flags are optional. Required options
should use positional arguments. Positional arguments are required unless
defined with a pointer type, and only if all subsequent positional arguments are
also optional.

### Optional command behavior

Every command struct must define a `Run() error` function to comply with the
`cli.Handler` interface. The struct can also define additional methods to
support specific behaviors:

- `Version() string`, if defined, adds a `-v, --version` option that print the
  command version returned by this function
- `Usage(name string, width int) string`, if defined, let the command completely
  redefine the usage printout triggered by `-h, --help` option
- `Complete(opt *option.T, partial string) []option.Description`, if defined,
  let the command override the list of suggestions offered during completion of
  an option or an argument. By default, the completion mechanism emulates the
  default behavior of bash completion and suggests matching local files.

### Completion support

Completion integration with bash is supported by running:
```
eval $(<command> --bash-completion-script)
```
Once setup, bash will invoke the command to get completion suggestions, with two
special environment variables set, `COMP_WORD` and `COMP_INDEX`.

When using zsh, you can still leverage bash completion scripts by adding the
following to you `~/.zshrc`, before th eval command:

```zsh
autoload -U +X compinit && compinit
autoload -U +X bashcompinit && bashcompinit
```

If for some reason a command should not support the built-in completion, the
completion machinery can be disabled by setting `cmd.DisableCompletion` on the
root command.


## Enum support

To help with command-line handling of enumerated types, the `go-cli` package
includes a code generator command for use with `go:generate` that adds the
necessary `flag.Value` interface methods to enumerated types.

The command is primarily intended to be run in the context of go:generate, from
a file that contains definitions for one or more enumerated types. In that
context, it is executed from the folder of that file and processes the current
package. It generates one file with `_enumer.go` suffix for every file in the
package that contains an enumerated type.

For example, the file pkg/strcase/format.go contains the following line:

```go
//go:generate go run github.com/maargenton/go-cli/cmd/enumer format.go
```

which generate `format_enumer.go` for the `Format` type. The command is
restricted to process only the definitions found in `format.go`.

The string representation of enumerated values is flexible and configurable, and
defaults to `filtered-hyphen-case`. The generated `Parse...()` function and
`Set()` method accept all supported string representations; the `String()`
method prints the value in the specified representation. The representation
format can be specified on the command-line with `-f` or `--format` option and
applies to all the enum types generated by one invocation of the command.

For additional convenience, the generated code also defines methods to support
both `encoding.TextMarshaler` and `encoding.TextUnmarshaler`, making the enum
values represented as a strings in both json and yaml serialization.
