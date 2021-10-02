package option

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/maargenton/go-cli/pkg/value"
)

// Set represents the set of options exposed by a struct through `opts`
// tags. It is attached to a specific instacne of that struct and allows field
// values to be set.
type Set struct {
	target reflect.Value

	Options    []*T
	Positional []*T
	Args       *T
}

// NewOptionSet creates a new Set that reflects the field in type `t`
// that can be set through commandline arguments
func NewOptionSet(v interface{}) (*Set, error) {
	var pv = reflect.ValueOf(v)
	if pv.Kind() != reflect.Ptr || pv.IsNil() || pv.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid argument of type '%v', non-null pointer to struct expected", pv.Type())
	}

	var opts = &Set{
		target: pv.Elem(),
	}
	if err := opts.parseStruct(); err != nil {
		return nil, err
	}

	return opts, nil
}

// GetOption return the option matching the specified name, which can be the
// short or the long name of the flag, without any leading dash. Returns nil is
// no matching flag is found.
func (opts *Set) GetOption(name string) (opt *T) {
	if name != "" {
		for _, opt := range opts.Options {
			if opt.Short == name || opt.Long == name {
				return opt
			}
		}
	}
	return nil
}

// AddSpecialFlag appends a special flag to the option set, that sends an
// sentinel error when found on the command-line. Used for `--version` and
// `--help`. The short flag is up-cased or dropped if conflicting with existing
// flags. The whole special flag is dropped if the long flag is conflicting.
func (opts *Set) AddSpecialFlag(short, long, desc string, err error) {
	if opts.GetOption(long) != nil {
		return
	}

	if opts.GetOption(short) != nil {
		short = strings.ToUpper(short)
		if opts.GetOption(short) != nil {
			short = ""
		}
	}

	var opt = &T{
		Short:       short,
		Long:        long,
		Description: desc,
		Type:        Special,
		SpecialErr:  err,

		opts: opts,
	}

	opts.Options = append(opts.Options, opt)
}

// ApplyDefaults scans through a parsed option set and applies the values
// defined in the environment to the fields backed by a matching environment
// variable.
func (opts *Set) ApplyDefaults() error {
	for _, opt := range opts.Options {
		if opt.Default != "" {
			if err := opt.SetValue(opt.Default); err != nil {
				return fmt.Errorf("while applying defaults, %w", err)
			}
		}
	}
	return nil
}

// ApplyEnv scans through a parsed option set and applies the corresponding
// default values to the fields of the target struct value.
func (opts *Set) ApplyEnv(env map[string]string) error {
	for _, opt := range opts.Options {
		if opt.Env != "" {
			if v, ok := env[opt.Env]; ok {
				if err := opt.SetValue(v); err != nil {
					return fmt.Errorf(
						"while applying value from environment variable '%v', %w",
						opt.Env, err)
				}
			}
		}
	}
	return nil
}

// ApplyArgs scans through a parsed option set and applies the corresponding
// command-line arguments to the fields of the target struct value.
func (opts *Set) ApplyArgs(args []string) error {
	opt, remainingArgs, err := opts.applyArgsToOptions(args)
	if err != nil {
		return err
	}
	if opt != nil {
		return fmt.Errorf("missing argument for '%v'", opt.Name())
	}

	for _, opt := range opts.Positional {
		if len(remainingArgs) == 0 {
			if opt.Optional {
				break
			}
			return fmt.Errorf("missing argument for '%v'", opt.Name())
		}
		if err := opt.SetValue(remainingArgs[0]); err != nil {
			return err
		}
		remainingArgs = remainingArgs[1:]
	}

	if len(remainingArgs) != 0 {
		if opts.Args != nil {
			for _, arg := range remainingArgs {
				if err := opts.Args.SetValue(arg); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf(
				"unsupported extra arguments: %v",
				strings.Join(remainingArgs, " "))
		}
	}

	return nil
}

func (opts *Set) applyArgsToOptions(
	args []string) (
	opt *T, remainingArgs []string, err error) {

	for _, arg := range args {
		if opt != nil {
			if err := opt.SetValue(arg); err != nil {
				return nil, nil, err
			}
			opt = nil

		} else if strings.HasPrefix(arg, "--") {
			optName := arg[2:]
			value := ""
			if i := strings.IndexByte(optName, '='); i >= 0 {
				value = optName[i+1:]
				optName = optName[:i]
			}
			opt = opts.GetOption(optName)
			if opt == nil {
				return nil, nil, &ErrInvalidFlag{arg}
			}
			if opt.Type == Special {
				return nil, nil, opt.SpecialErr
			}
			if opt.Type == Bool && value == "" {
				opt.SetBool()
				opt = nil
			}
			if value != "" {
				if err := opt.SetValue(value); err != nil {
					return nil, nil, err
				}
				opt = nil
			}
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
			for i, c := range arg {
				opt = opts.GetOption(string(c))
				if opt == nil {
					return nil, nil, &ErrInvalidFlag{"-" + string(c)}
				}
				if opt.Type == Special {
					return nil, nil, opt.SpecialErr
				}
				if opt.Type == Bool {
					opt.SetBool()
					opt = nil
				} else {
					value := arg[i+1:]
					if len(value) > 0 {
						if err := opt.SetValue(value); err != nil {
							return nil, nil, err
						}
						opt = nil
					}
					break
				}
			}
		} else {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return
}

// ---------------------------------------------------------------------------
// Private support functions for NewOptionSet()
// ---------------------------------------------------------------------------

func (opts *Set) parseStruct() error {

	var t = opts.target.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		err := opts.parseField(field, []int{})
		if err != nil {
			return fmt.Errorf(
				"error while generating Set for '%v': %v",
				t, err)
		}
	}

	var positional = make(map[int]*T)
	for i, opt := range opts.Options {
		if opt.Args {
			if opts.Args == nil {
				opts.Args = opt
				opts.Options[i] = nil
			} else {
				return fmt.Errorf(
					"multiple fields capturing remaning args: '%v' amd '%v'",
					opts.Args.FieldName, opt.FieldName,
				)
			}
		} else if p := opt.Position; p != 0 {
			if o, exists := positional[p]; exists {
				return fmt.Errorf(
					"positional argument '%v' defined by multiple fields: '%v' amd '%v'",
					p, o.FieldName, opt.FieldName,
				)
			}
			positional[p] = opt
			opts.Options[i] = nil
		}
	}

	for i := 1; i <= len(positional); i++ {
		opt, exists := positional[i]
		if !exists {
			return fmt.Errorf(
				"%v positional arguments defined, but 'arg:%v' is missing",
				len(positional), i,
			)
		}
		opts.Positional = append(opts.Positional, opt)
	}

	var compact []*T
	for _, opt := range opts.Options {
		if opt != nil {
			compact = append(compact, opt)
		}
	}
	opts.Options = compact

	if opts.Args != nil && opts.Args.Type != Slice {
		return fmt.Errorf(
			"field '%v' of type '%v' must be a slice to receive additional arguments",
			opts.Args.Name(), opts.Args.FieldType)
	}

	for _, arg := range opts.Positional {
		if arg.Type == Slice {
			return fmt.Errorf(
				"field '%v' of type '%v' cannot be a slice to receive positional argument",
				arg.Name(), arg.FieldType)

		}
	}

	// Traverse positional arguments backward and mark all trailing pointer type
	// positional as optional.
	for i := len(opts.Positional) - 1; i >= 0; i-- {
		arg := opts.Positional[i]
		if arg.FieldType.Kind() == reflect.Ptr {
			arg.Optional = true
		} else {
			break
		}
	}

	return nil
}

func (opts *Set) parseField(f reflect.StructField, index []int) error {

	if f.Anonymous && f.Type.Kind() == reflect.Struct {
		var t = f.Type
		var fieldIndex = mergeIndexes(index, f.Index)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			err := opts.parseField(field, fieldIndex)
			if err != nil {
				return err
			}
		}
	} else {
		tag, ok := f.Tag.Lookup("opts")
		if !ok {
			return nil
		}

		var fieldType = f.Type
		var valueType = fieldType
		if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice {
			valueType = fieldType.Elem()
		}

		if !value.CanParseType(valueType) {
			return fmt.Errorf(
				"type '%v' of field '%v' is not parsable",
				fieldType, f.Name)
		}

		var index = mergeIndexes(index, f.Index)
		var vv = opts.target.FieldByIndex(index)
		if !vv.CanSet() {
			return fmt.Errorf(
				"field '%v' of '%v' is not settable",
				f.Name, opts.target.Type())

		}

		var optionType = Value
		if fieldType.Kind() == reflect.Slice {
			optionType = Slice
		} else if valueType.Kind() == reflect.Bool {
			optionType = Bool
		} else if fieldType.Kind() == reflect.Ptr {
			optionType = Ptr
		}

		var opt = &T{
			FieldName: f.Name,
			FieldType: fieldType,
			ValueType: valueType,
			Type:      optionType,
			Index:     index,
			opts:      opts,
		}
		err := opt.parseOptsTag(tag)
		if err != nil {
			return err
		}

		if desc, ok := f.Tag.Lookup("desc"); ok {
			opt.Description = strings.TrimSpace(desc)
		}

		opts.Options = append(opts.Options, opt)
	}

	return nil
}

func mergeIndexes(indexes ...[]int) []int {
	var r []int
	for _, ii := range indexes {
		r = append(r, ii...)
	}
	return r
}
