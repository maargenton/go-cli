package option

import (
	"strings"
)

// Completion records a set of completion suggestions, including usable flags,
// values for a specific flag and / or values for next remaining argument.
type Completion struct {
	Options   []Description
	Opt       *T
	OptValues []Description
	Arg       *T
	ArgValues []Description
}

func (opts *Set) GetCompletion(args []string, partial string) Completion {

	var suggestions Completion

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
				if opt.Type == BoolType || opt.Type == SpecialType {
					opt = nil // no value expected
				}
			}
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
			for i, c := range arg {
				opt = opts.GetOption(string(c))
				if opt != nil {
					usedOptions[opt] = struct{}{}
					if opt.Type == BoolType || opt.Type == SpecialType {
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
		suggestions.Opt = opt
		return suggestions
	}

	var nonExclusiveUsed = len(remainingArgs) > 0
	for o := range usedOptions {
		if o.Type == SpecialType {
			// Exclusive flag has been used nothin more to suggest
			return suggestions
		} else {
			nonExclusiveUsed = true
		}
	}
	for _, o := range opts.Options {
		if _, used := usedOptions[o]; !used || o.Type == SliceType {
			if o.Type == SpecialType && nonExclusiveUsed {
				// Non-exclusive flag has been used, skip special flags
				continue
			}
			suggestions.Options = append(suggestions.Options, o.GetCompletionUsage())
		}
	}

	if len(remainingArgs) < len(opts.Positional) {
		suggestions.Arg = opts.Positional[len(remainingArgs)]
	} else if opts.Args != nil {
		suggestions.Arg = opts.Args
	}

	return suggestions
}
