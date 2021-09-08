# go-cli

A library to define command line interfaces in Go.

[![GoDoc](
    https://godoc.org/github.com/maargenton/go-cli?status.svg)](
    https://godoc.org/github.com/maargenton/go-cli)
[![Build Status](
    https://travis-ci.org/maargenton/go-cli.svg?branch=master)](
    https://travis-ci.org/maargenton/go-cli)
[![codecov](
    https://codecov.io/gh/maargenton/go-cli/branch/master/graph/badge.svg)](
    https://codecov.io/gh/maargenton/go-cli)
[![Go Report Card](
    https://goreportcard.com/badge/github.com/maargenton/go-cli)](
    https://goreportcard.com/report/github.com/maargenton/go-cli)


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
A collection of various example commands using go-cli can be found at

```bash
git clone git@github.com:maargenton/go-cli-examples
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
    cli.Run(os.Args, &cli.Command{
        Handler:     &cmd{},
        Description: "...",
    })
}
```

The `cmd` type defines a collection of options for command being defined. The
main function parses command-line arguments and invokes `cmd.Run()` with all
the fields initialized inside the `options`.

The `opts` tag defines, in a comma separated list:
- `-b,--baudrate`: either or both of a short and long flag name for the option
- `arg:<n>` : captures a positional argument
- `args` : captures all remaining arguments
- `default` : a default value for the field if not specified on the command-line
- `env` : the name of an environment variable that can override the default
- `sep`: a separator for fields that can accept multiple values. By
  default, multiple values must be specified by repeating the option flag multiple times.
- `name`: the display name for the value, used when printing out description of
  the field.

A separate `desc` tag contains the description for the option.

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
