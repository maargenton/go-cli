package option

import (
	"strings"
)

// Completion records a set of completion suggestions, including usable flags,
// values for a specific flag and / or values for next remaining argument.
type Completion struct {
	Options   []string
	OptValues []string
	ArgValues []string
	OptRef    *T
	ArgRef    *T
}

// GetCompletion evaluate the list of command line arguments `args` in the
// context of the receiver, and determines a list of completion suggestions for
// the `partial` argument given. The result is a partially filled `Completion`
// object with either a list of `Options` or one of `Opt` or `Arg` set the the
// `option.T` whose value needs to be completed.
func (opts *Set) GetCompletion(args []string, partial string) (r Completion) {

	// Evaluate commandline arguments, discarding values
	var opt *T
	var remainingArgs []string
	var usedOptions = make(map[*T]struct{})
	for _, arg := range args {
		if opt != nil {
			opt = nil // swallow value
		} else if strings.HasPrefix(arg, "--") {
			opt = opts.GetOption(arg[2:])
			if opt != nil {
				usedOptions[opt] = struct{}{}
				if opt.Type == Bool || opt.Type == Special {
					opt = nil // no value expected
				}
			}
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
			for i, c := range arg {
				opt = opts.GetOption(string(c))
				if opt != nil {
					usedOptions[opt] = struct{}{}
					if opt.Type == Bool || opt.Type == Special {
						opt = nil // no value expected
					} else {
						value := arg[i+1:]
						if len(value) > 0 {
							opt = nil
						}
						break
					}
				}
			}
		} else {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	if opt != nil {
		r.OptRef = opt
		return r
	}

	var nonExclusiveUsed = len(remainingArgs) > 0
	for o := range usedOptions {
		if o.Type == Special {
			// Exclusive flag has been used, nothing more to suggest
			return r
		}
		nonExclusiveUsed = true
	}
	for _, o := range opts.Options {
		if _, used := usedOptions[o]; !used || o.Type == Slice {
			if o.Type == Special && nonExclusiveUsed {
				// Non-exclusive flag has been used, skip special flags
				continue
			}
			r.Options = append(r.Options, o.Name())
		}
	}

	if len(remainingArgs) < len(opts.Positional) {
		r.ArgRef = opts.Positional[len(remainingArgs)]
	} else if opts.Args != nil {
		r.ArgRef = opts.Args
	}

	return r
}
