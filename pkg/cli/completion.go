package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// DefaultCompletion implements a default completion for a given option field,
// and can be used a fallback by command completion handlers. It handles
// specific cases based on the option field type, and simulates default shell
// behavior (filename completion) for string types.
func DefaultCompletion(opt *option.T, w string) []string {
	return DefaultFilenameCompletion(opt, w)
}

// DefaultFilenameCompletion implements a default filename-based completion,
// similar to the default shell completion.
func DefaultFilenameCompletion(opt *option.T, w string) (r []string) {
	files, err := dir.Glob(fmt.Sprintf("%v*", w))
	if err == nil {
		if len(files) == 1 && strings.HasSuffix(files[0], string(fileutils.Separator)) {
			return DefaultFilenameCompletion(opt, files[0])
		}
		for _, f := range files {
			r = append(r, f)
		}
	}
	return
}

// MatchingFileCompletion implements a custom filepath completion scheme,
// matching the provided pattern if possible, using default filename completion
// as a fallback.
func MatchingFilenameCompletion(opt *option.T, pattern string, w string) (r []string) {
	files, err := dir.Glob(pattern)
	if err == nil {
		for _, f := range files {
			if strings.HasPrefix(f, w) {
				r = append(r, f)
			}
		}
		if len(r) > 0 {
			return r
		}
	}
	return DefaultFilenameCompletion(opt, w)
}

// ---

func (cmd *Command) handleCompletionRequest() bool {
	i, w := cmd.getCompletionRequest()
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
	if comp.OptRef != nil {
		comp.OptValues = cmd.complete(comp.OptRef, w)
	}
	if comp.ArgRef != nil {
		comp.ArgValues = cmd.complete(comp.ArgRef, w)
	}

	for _, o := range comp.Options {
		if strings.HasPrefix(o, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}
	for _, o := range comp.OptValues {
		if strings.HasPrefix(o, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}
	for _, o := range comp.ArgValues {
		if strings.HasPrefix(o, w) {
			cmd.Suggestions = append(cmd.Suggestions, o)
		}
	}

	dumpCompletionRequest(comp)

	return true
}

func (cmd *Command) getCompletionRequest() (index int, word string) {
	var env = cmd.ProcessEnv

	word = env["COMP_WORD"]
	if ii, ok := env["COMP_INDEX"]; ok {
		if ii, err := strconv.ParseInt(ii, 0, 0); err == nil {
			index = int(ii)
		}
	}

	var args = cmd.ProcessArgs[:]
	if index == 0 || index >= len(args) || !strings.HasPrefix(args[index], word) {
		word = ""
	}
	return
}

func (cmd *Command) complete(opt *option.T, word string) []string {
	if handler, ok := cmd.Handler.(CompletionHandler); ok {
		return handler.Complete(opt, word)
	}
	return DefaultCompletion(opt, word)
}

func dumpCompletionRequest(comp option.Completion) {
	if s, ok := os.LookupEnv("COMPLETION_DEBUG_OUTPUT"); ok && s != "" {
		type obj = map[string]interface{}
		var v = obj{
			"args": os.Args,
			"env": obj{
				"COMP_INDEX": os.Getenv("COMP_INDEX"),
				"COMP_WORD":  os.Getenv("COMP_WORD"),
			},
			"comp": comp,
		}

		fileutils.WriteFile(s, func(w io.Writer) error {
			var e = json.NewEncoder(w)
			e.SetIndent("", "    ")
			return e.Encode(v)
		})
	}
}
