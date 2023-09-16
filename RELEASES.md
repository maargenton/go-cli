# v0.5.0

## Key Features

- Simplify completion option and associated interface to a simple list of
  strings
- Handle special case in zsh completions, discard `COMP_WORD=--` when not
  matching start of argument at `COMP_INDEX`
- Change available completion helper functions to:
    - `DefaultCompletion()` to return the default completion behavior from a
      custom completion handler
    - `DefaultFilenameCompletion()` matching default shell completion
    - `MatchingFilenameCompletion()` for pattern based filename matching
- Surrounding spaces are automatically trimmed from values unless `keep-spaces`
  option is specified on the argument
- List values automatically drop empty values unless `keep-empty` option is
  specified on the argument. Note that even with `keep-empty`, if the last
  character is a separator, the last empty values is always drop, as was the
  case in the prior versions.
- Update `golang.org/x/tools` package dependency to fix issues in `enumer`
  command loading some of the standard packages

## Code changes

- Simplify completion handling and suggestions ([#9](https://github.com/maargenton/go-cli/pull/9))
- Update build scripts, fix windows test issues ([#11](https://github.com/maargenton/go-cli/pull/11))
- Change handling of spaces and empty list values ([#10](https://github.com/maargenton/go-cli/pull/10))
- Fix `enumer` command package loading issues ([#13](https://github.com/maargenton/go-cli/pull/13))


# v0.4.0

## Key Features

- Add `enumer` command to generate flag.Value interface methods for enumerated
  types and make them usable as direct recipients for command-line arguments
- Add support for custom value parser defined on pointer types (like
  `url.Parse`)

## Improvements

- Clarify README `opts` struct tag documentation
- Add note in README for setting up bash completion compatibility in zsh

# v0.3.1

## Major Features

- Bump minimum Go version requirement to v1.17 due to dependencies
- Preserve preset `cmd.ProcessName` and `cmd.ProcessArgs` if set before invoking
  `cli.Run()`
- Add option to disable completion machinery
- Improve usage display for positional arguments

## Improvements

- Fix typos in test names
- Update dependencies
- Add vscode debug configuration


# v0.3.0

## Major Features

- Add support for `--long=<value>` format
- Add support for empty value in `--long=<value>` format
- Add support for `--` end of options delimiter

## Improvements

- Improve reported errors from `value.Parse()`
- Update documentation with details about short flags, long flags, non-option
  arguments and limitations

# v0.2.0

## Major Features

- Resolve windows compatibility issues, update go-filetuils to v0.6.0, which
  addresses the same issue.

# v0.1.0

## Major Features

- Support for simple commands, no sub-command
- Parses flags, flag values and arguments
- Full support for bash completion
