package value

import (
	"encoding"
	"flag"
	"fmt"
	"reflect"
)

// RegisterParser registers a function of type 'func(string) (T, error)'
// as value parser for type T
func RegisterParser(parsers ...interface{}) {
	for _, parser := range parsers {
		var f = makeParserFunc(parser)
		var valueType = f.valueType()

		if ff, exist := parseFunctions[valueType]; exist {
			if f != ff {
				panic(fmt.Sprintf(
					"parse function for type '%v' is already registered",
					valueType))
			}
		} else {
			parseFunctions[valueType] = f
		}
	}
}

// CanParseType returns true if the specified type has a registered
// value parser, implements to flag.Value interface or implements
// encoding.TextUnmarshaler interface.
func CanParseType(t reflect.Type) bool {
	if _, exist := parseFunctions[t]; exist {
		return true
	}

	var p = reflect.PtrTo(t)
	if p.Implements(flagValueType) {
		return true
	}

	if p.Implements(textUnmarshalerType) {
		return true
	}

	return false
}

// Parse converts a string into an actual value, using a string
// conversion function specific to the target type. `v` must be a non-nil
// pointer to a variable to parse into, and panics otherwise.
// Types with a registered parser function, including all common primitive types
// are converted using that parser function. Types that conform to the
// flag.Value interface are converted with the Set() method. Types that conform
// to the encoding.TextUnmarshaler interface are converted using the
// UnmarshalText() method.
// The function panics if the target ype is not parsable.
func Parse(v interface{}, s string) error {

	checkPtrToVar(v)
	var vv = reflect.ValueOf(v)
	var valueType = vv.Type().Elem()

	if f, exist := parseFunctions[valueType]; exist {
		r, err := f.call(s)
		if err != nil {
			return fmt.Errorf(
				"invalid value '%v' for type '%v': %w",
				s, valueType, err)
		}
		vv.Elem().Set(reflect.ValueOf(r))
		return nil
	}

	var err error
	if fv, ok := v.(flag.Value); ok {
		err = fv.Set(s)
	} else if tv, ok := v.(encoding.TextUnmarshaler); ok {
		err = tv.UnmarshalText([]byte(s))
	} else {
		err = fmt.Errorf(
			"value.Parse() called with non-parsable value type '%v'",
			valueType)
		panic(err)
	}

	if err != nil {
		err = fmt.Errorf(
			"invalid value '%v' for type '%v': %w",
			s, valueType, err)
	}
	return err
}

func checkPtrToVar(v interface{}) {
	var vv = reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		var msg = fmt.Sprintf(
			"call to value.Parse() requires a non-nil pointer type, got '%v'",
			vv.Type())
		panic(msg)
	}
	if vv.IsNil() {
		var msg = fmt.Sprintf(
			"call to value.Parse() requires a non-nil pointer type, got '(%v)(%v)'",
			vv.Type(), vv.Interface())
		panic(msg)
	}
}

// ---------------------------------------------------------------------------
// reflect.Type for common types
// ---------------------------------------------------------------------------

var (
	stringType          = reflect.TypeOf((*string)(nil)).Elem()
	errorType           = reflect.TypeOf((*error)(nil)).Elem()
	flagValueType       = reflect.TypeOf((*flag.Value)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

// ---------------------------------------------------------------------------
// Global registry for value parser
// ---------------------------------------------------------------------------

type parserFunc reflect.Value

var parseFunctions = map[reflect.Type]parserFunc{}

func (fct parserFunc) valueType() reflect.Type {
	return reflect.Value(fct).Type().Out(0)
}

func (fct parserFunc) call(s string) (interface{}, error) {
	var args = []reflect.Value{reflect.ValueOf(s)}
	var results = reflect.Value(fct).Call(args)
	var result = results[0].Interface()
	var err error
	if r2, ok := results[1].Interface().(error); ok {
		err = r2
	}

	return result, err
}

func makeParserFunc(f interface{}) parserFunc {
	var v = reflect.ValueOf(f)
	var t = v.Type()

	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("value of type '%v' is not a function", t))
	}

	if t.NumIn() != 1 || t.In(0) != stringType {
		panic(fmt.Sprintf("value of type '%v' is not a valid parse function", t))
	}

	if t.NumOut() != 2 || t.Out(1) != errorType {
		panic(fmt.Sprintf("value of type '%v' is not a valid parse function", t))
	}

	return parserFunc(reflect.ValueOf(f))
}
