package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-fileutils/pkg/dir"

	"github.com/maargenton/go-cli/pkg/option"
)

// BashCompletionScript returns a string containing the bash script necessary to
// setup bash completion for the command.
func BashCompletionScript(command string) string {
	var bashCompletionTemplate string = "" +
		"_%[1]v_completion() {\n" +
		"    local IFS=$'\\n' ;\n" +
		"    COMPREPLY=($(COMP_INDEX=$COMP_CWORD COMP_WORD=$2 ${COMP_WORDS[@]})) ;\n" +
		"    return 0 ;\n" +
		"} ;\n" +
		"complete -F _%[1]v_completion %[1]v ;\n"

	return fmt.Sprintf(bashCompletionTemplate, command)
}

// FormatCompletionSuggestions takes pre-filtered completion suggestions and
// returns a formatted string ready to pass back to the shell during a
// completion request. The description and addition suffixes are dropped when
// only one option is available.
func FormatCompletionSuggestions(width int, suggestions []option.Description) string {
	var s strings.Builder
	if len(suggestions) == 1 {
		v := suggestions[0].Option
		v = strings.Split(v, " ")[0]
		fmt.Fprintf(&s, "%v\n", v)
	} else {
		fmt.Fprintf(&s, "%v", option.FormatCompletion(width, suggestions))
	}

	return s.String()
}

// DefaultCompletion acts like the shell default completion and suggests file
// and folder names under the current directory. It is used by default when the
// command does not implement a specific completion handler, and should be used
// from the command completion handler when no other completion logic is
// suitable.
func DefaultCompletion(w string) []option.Description {
	var r []option.Description
	files, err := dir.Glob(fmt.Sprintf("%v*", w))
	if err == nil {
		if len(files) == 1 && strings.HasSuffix(files[0], string(fileutils.Separator)) {
			return DefaultCompletion(files[0])
		}
		for _, f := range files {
			r = append(r, option.Description{Option: f})
		}
	}
	return r
}

// FilepathCompletion implements a custom filepath completion scheme, matching
// the provided pattern if possible. If the pattern does not yield any result,
// or if none of the suggestions match the partial argument being completed,
// then if folds back to the default filename completion.
func FilepathCompletion(pattern string, w string) []option.Description {
	files, err := dir.Glob(pattern)
	if err == nil {
		var r []option.Description
		for _, f := range files {
			if strings.HasPrefix(f, w) {
				r = append(r, option.Description{Option: f})
			}
		}
		if len(r) > 0 {
			return r
		}
	}
	return DefaultCompletion(w)
}

// ---

func (cmd *Command) handleCompletionRequest() bool {
	i, w := checkCompletionRequest(cmd.ProcessEnv)
	if i == 0 {
		return false
	}

	// Trim process arguments past the completion point
	var args = cmd.ProcessArgs[:]
	if i > 0 && i < len(args) {
		args = args[:i]
	}
	if len(args) > 0 {
		args = args[1:]
	}

	var comp = cmd.opts.GetCompletion(args, w)
	if comp.Opt != nil {
		comp.OptValues = cmd.complete(comp.Opt, w)
	}
	if comp.Arg != nil {
		comp.ArgValues = cmd.complete(comp.Arg, w)
	}

	for _, o := range comp.Options {
		if strings.HasPrefix(o.Option, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}
	for _, o := range comp.OptValues {
		if strings.HasPrefix(o.Option, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}
	for _, o := range comp.ArgValues {
		if strings.HasPrefix(o.Option, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}

	return true
}

func checkCompletionRequest(env map[string]string) (index int, word string) {
	i, ok := env["COMP_INDEX"]
	w, ok2 := env["COMP_WORD"]
	if ok && ok2 {
		if ii, err := strconv.ParseInt(i, 0, 0); err == nil {
			index = int(ii)
			word = w
		}
	}
	return
}

func (cmd *Command) complete(opt *option.T, word string) []option.Description {
	if handler, ok := cmd.Handler.(CompletionHandler); ok {
		return handler.Complete(opt, word)
	}
	return DefaultCompletion(word)
}
