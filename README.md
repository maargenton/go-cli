# go-cli

A library to define command line interfaces in Go.

[![Latest](
  https://img.shields.io/github/v/tag/maargenton/go-cli?color=blue&label=latest&logo=go&logoColor=white&sort=semver)](
  https://pkg.go.dev/github.com/maargenton/go-cli)
[![Build](
  https://img.shields.io/github/workflow/status/maargenton/go-cli/build?label=build&logo=github&logoColor=aaaaaa)](
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
- `sep:`: a separator for fields that can accept multiple values. By default,
  multiple values must be specified by repeating the option flag multiple times.
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
initialized to their built-in zero value. If the built-in zero value is a valid
value and if the command needs to determine if an option has been provided, a
pointer type should be used and checked against `nil`.

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

If for some reason a command should not support the built-in completion, the
completion machinery can be disabled by setting `cmd.DisableCompletion` on the
root command.
