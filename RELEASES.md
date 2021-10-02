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
