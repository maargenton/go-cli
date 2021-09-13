# go-cli

A library to define command line interfaces in Go.

[![Latest](
  https://img.shields.io/github/v/tag/maargenton/go-cli?color=blue&label=latest&logo=go&logoColor=white&sort=semver)](
  https://pkg.go.dev/github.com/maargenton/go-cli)
[![Build](
  https://img.shields.io/github/workflow/status/maargenton/go-cli/build?label=build&logo=github&logoColor=aaaaaa)](
  https://github.com/maargenton/go-cli/actions?query=branch%3Amaster)
[![Codecov](
  https://img.shields.io/codecov/c/github/maargenton/go-cli?label=codecov&logo=codecov&logoColor=aaaaaa&token=fVZ3ZMAgfo)](
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


## Standards

The implementation of go-clo follows the conventions outlined in [The Open Group
    Base Specifications Issue 7, 2018 edition, Chapter 12. Utility
    Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)


## Installation

```bash
go get github.com/maargenton/go-cli
```

## Usage

```go
type cmd struct {
    Port     string  `opts:"arg:1, name:port" desc:"..."`
    Baudrate *uint32 `opts:"-b, --baudrate"   desc:"..."`
    Format   *string `opts:"-f, --format"     desc:"..."`

    Timestamp bool `opts:"-t, --timestamp" desc:"..."`
    Verbose   bool `opts:"-v, --verbose"   desc:"..."`
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
struct tags, containing a comma separated list of the following:
- `-b,--baudrate`: either or both of a short and long flag name for the option
- `arg:<n>` : captures a positional argument
- `args` : captures all remaining arguments
- `default` : a default value for the field if not specified on the command-line
- `env` : the name of an environment variable that can override the default
- `sep`: a separator for fields that can accept multiple values. By default,
  multiple values must be specified by repeating the option flag multiple times.
- `name`: the display name for the value, used when printing out description of
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

- `Version() string`, if defined, adds a `-v, --verions` option that print the
  command version returned by thise function
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
