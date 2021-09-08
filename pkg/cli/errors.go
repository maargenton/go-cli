package cli

import (
	"github.com/maargenton/go-errors"
)

// ErrHelpRequested is a sentinel error indicating that help was requested with
// no explicit flag to handle the request
const ErrHelpRequested = errors.Sentinel("ErrHelpRequested")

// ErrVersionRequested is a sentinel error indicating that version was requested
// with no explicit flag to handle the request
const ErrVersionRequested = errors.Sentinel("ErrVersionRequested")

// ErrCompletionScriptRequested is a sentinel error indicating the bash
// completion script was requested and should be printed to stdout for
// evaluation
const ErrCompletionScriptRequested = errors.Sentinel("ErrCompletionScriptRequested")

// ErrCompletionRequested is a sentinel error indicating that the command was
// invoked in completion mode and that the completion options should be printed
// out instead of running the command
const ErrCompletionRequested = errors.Sentinel("ErrCompletionRequested")
