package option

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/maargenton/go-cli/pkg/value"
)

type OptionType int

const (
	ValueType OptionType = iota
	BoolType
	PtrType
	SliceType
	SpecialType
)

// T represents a single option with a reference back to the OptionSet it
// belongs to.
type T struct {
	Short       string
	Long        string
	Default     string
	Env         string
	Sep         string // optional separator
	Description string
	ValueName   string // optional name for the value
	Position    int    // set to non-zero for fields capturing positional arguments
	Args        bool   // set to true for the field capturing remaining arguments

	FieldName  string
	Index      []int
	FieldType  reflect.Type
	ValueType  reflect.Type
	Type       OptionType
	Optional   bool
	SpecialErr error

	opts *Set
}

// Description captures both an option and its description. This is used
// for both usage and completion.
type Description struct {
	Option      string
	Description string
}

// Name returns a displayable name for the target option, using either the long
// flag, the short flag or the field name.
func (opt *T) Name() string {
	if opt.Args {
		if opt.ValueName != "" {
			return fmt.Sprintf("<%v>...", opt.ValueName)
		}
		return "<args>..."
	}
	if opt.Position != 0 {
		if opt.ValueName != "" {
			return fmt.Sprintf("<%v>", opt.ValueName)
		}
		return fmt.Sprintf("<arg[%d]>", opt.Position)
	}
	if opt.Long == "" {
		return fmt.Sprintf("-%v", opt.Short)
	}
	return fmt.Sprintf("--%v", opt.Long)
}

// GetUsage returns a formated representation of the option used to display
// usage and its description.
func (opt *T) GetUsage() (usage Description) {

	var u strings.Builder
	if opt.Short != "" && opt.Long != "" {
		fmt.Fprintf(&u, "-%v, --%v", opt.Short, opt.Long)
	} else if opt.Short != "" {
		fmt.Fprintf(&u, "-%v", opt.Short)
	} else if opt.Long != "" {
		fmt.Fprintf(&u, "    --%v", opt.Long)
	}

	if v := opt.getValueDescription(); v != "" {
		if u.Len() > 0 {
			u.WriteRune(' ')
		}
		u.WriteString(v)
	}

	return Description{
		Option:      u.String(),
		Description: opt.getDescription(),
	}
}

// GetCompletion returns a value similar to GetUsage(), but using only the long
// flag if defined, the short flag otherwise.
func (opt *T) GetCompletionUsage() (usage Description) {
	var u strings.Builder
	if opt.Long != "" {
		fmt.Fprintf(&u, "--%v", opt.Long)
	} else {
		fmt.Fprintf(&u, "-%v", opt.Short)
	}

	if v := opt.getValueDescription(); v != "" {
		if u.Len() > 0 {
			u.WriteRune(' ')
		}
		u.WriteString(v)
	}

	return Description{
		Option:      u.String(),
		Description: opt.getDescription(),
	}
}

func (opt *T) getValueDescription() string {
	if opt.Type == BoolType || opt.Type == SpecialType {
		return ""
	}
	if opt.ValueName != "" {
		return fmt.Sprintf("<%v>", opt.ValueName)
	}
	return "<value>"
}

func (opt *T) getDescription() string {
	var d strings.Builder
	d.WriteString(opt.Description)
	if opt.Default != "" {
		if d.Len() > 0 {
			fmt.Fprintf(&d, ", ")
		}
		fmt.Fprintf(&d, "default: %v", opt.Default)
	}

	if opt.Env != "" {
		if d.Len() > 0 {
			fmt.Fprintf(&d, ", ")
		}
		fmt.Fprintf(&d, "env: %v", opt.Env)
	}

	return d.String()
}

// SetBool is a special setter usable only on boolean flags to set them to true.
func (opt *T) SetBool() {
	if opt.Type != BoolType {
		panic("cannot call Option.SetBool() on non-bool fields")
	}

	var fv = opt.opts.target.FieldByIndex(opt.Index)
	if fv.Kind() == reflect.Ptr {
		var v = reflect.New(fv.Type().Elem())
		v.Elem().SetBool(true)
		fv.Set(v)
	} else {
		fv.SetBool(true)
	}
}

// SetValue convert the strign `s` into a value of the desired type and assigns
// it to the struct field backing the current options. It supports value type,
// pointer type and slice type. For pointer type and slice type, an empty value
// reverts the field to a null pointer or an empty slice. For slice types
// defining a delimiter, the value is split accordingly and the delimited values
// are added to the slice.
func (opt *T) SetValue(s string) error {
	var fv = opt.opts.target.FieldByIndex(opt.Index)
	var err error
	if fv.Kind() == reflect.Ptr {
		err = opt.setPtrValue(fv, s)
	} else if fv.Kind() == reflect.Slice {
		err = opt.setSliceValue(fv, s)
	} else {
		err = value.Parse(fv.Addr().Interface(), s)
	}
	if err != nil {
		err = fmt.Errorf("failed to set value for '%v': %w", opt.Name(), err)
	}
	return err
}

func (opt *T) setPtrValue(fv reflect.Value, s string) error {
	var v = reflect.New(fv.Type().Elem())
	if err := value.Parse(v.Interface(), s); err != nil {
		return err
	}
	fv.Set(v)

	return nil
}

func (opt *T) setSliceValue(fv reflect.Value, s string) error {
	if len(s) == 0 {
		fv.Set(reflect.Zero(fv.Type()))
	} else if opt.Sep != "" {
		var updatedSlice = fv
		for _, vs := range splitSliceValues(s, opt.Sep) {
			var v = reflect.New(fv.Type().Elem())
			if err := value.Parse(v.Interface(), vs); err != nil {
				return err
			}
			updatedSlice = reflect.Append(updatedSlice, v.Elem())
		}
		fv.Set(updatedSlice)
	} else {
		var v = reflect.New(fv.Type().Elem())
		if err := value.Parse(v.Interface(), s); err != nil {
			return err
		}
		fv.Set(reflect.Append(fv, v.Elem()))
	}
	return nil
}

func splitSliceValues(s string, delim string) (r []string) {
	var escape = false
	var b strings.Builder

	for _, c := range s {
		if escape {
			escape = false
			b.WriteRune(c)
		} else if strings.ContainsRune(delim, c) {
			r = append(r, b.String())
			b.Reset()
		} else if c == '\\' {
			escape = true
		} else {
			b.WriteRune(c)
		}
	}
	if b.Len() > 0 {
		r = append(r, b.String())
	}

	return r
}

func (opt *T) parseOptsTag(tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("invalid empty 'opts' tag")
	}

	var s = tag
	for len(s) > 0 {
		var k, v string
		k, v, s = scanTagFields(s)

		if strings.HasPrefix(k, "--") && len(k) > 3 {
			opt.Long = k[2:]
		} else if k != "--" && strings.HasPrefix(k, "-") && len(k) == 2 {
			opt.Short = k[1:]
		} else if k == "args" && v == "" {
			opt.Args = true
		} else if k == "arg" {
			n, err := strconv.ParseInt(v, 0, 0)
			if err != nil {
				return fmt.Errorf("invalid index '%v' for arg: tag, %w", v, err)
			}
			if n == 0 {
				return fmt.Errorf("invalid index '0' for arg: tag")
			}
			opt.Position = int(n)
		} else if k == "default" {
			opt.Default = v
		} else if k == "env" {
			opt.Env = v
		} else if k == "sep" {
			opt.Sep = v
		} else if k == "name" {
			opt.ValueName = v
		} else {
			return fmt.Errorf("invalid tag in opts: '%v'", k)
		}
	}

	return nil
}
