package option

import (
	"fmt"
)

// ErrInvalidFlag is a custom error type, raised while parsing command-line
// arguments, indicating that a specific flag was not recognized
type ErrInvalidFlag struct {
	Flag string
}

func (err *ErrInvalidFlag) Error() string {
	return fmt.Sprintf("invalid flag '%v'", err.Flag)
}
