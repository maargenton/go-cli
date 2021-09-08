package value_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/maargenton/go-cli/pkg/value"

	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestRegisterParserPanicsWithInvalidArgument(t *testing.T) {

	require.That(t, func() {
		value.RegisterParser(time.ParseDuration("1s"))
	}).Panics()

	require.That(t, func() {
		value.RegisterParser(func(int) {})
	}).Panics()

	require.That(t, func() {
		value.RegisterParser(func(string) error { return nil })
	}).Panics()

	require.That(t, func() {
		value.RegisterParser(func(string) (int, int) { return 0, 0 })
	}).Panics()

}

func TestRegisterParserPanicsWithMultipleParserForTheSameType(t *testing.T) {
	value.RegisterParser(ParseCustomParserType)
	require.That(t, func() {
		value.RegisterParser(ParseCustomParserType2)
	}).Panics()
}

type customParserType int

func ParseCustomParserType(s string) (customParserType, error) {
	return customParserType(0), nil
}

func ParseCustomParserType2(s string) (customParserType, error) {
	return customParserType(0), nil
}

// ---------------------------------------------------------------------------
// Tests with flag.Value conforming type
// ---------------------------------------------------------------------------
func TestFlagValue(t *testing.T) {
	var customValueType = reflect.TypeOf((*customValue)(nil)).Elem()
	require.That(t, value.CanParseType(customValueType)).Eq(true)

	var v customValue
	var err = value.Parse(&v, "foobar")
	require.That(t, err).IsNil()
	require.That(t, v).Eq("foobar")
}

func TestFlagValueError(t *testing.T) {

	var customValueType = reflect.TypeOf((*customValue)(nil)).Elem()
	require.That(t, value.CanParseType(customValueType)).Eq(true)

	var v customValue
	var err = value.Parse(&v, "invalid-value")
	require.That(t, err).IsNotNil()
	require.That(t, v).Eq("")
}

type customValue string

func (v customValue) String() string {
	return string(v)
}

func (v *customValue) Set(s string) error {
	if s == "invalid-value" {
		return fmt.Errorf("cannot convert '%v' into customValue", s)
	}
	*v = customValue(s)
	return nil
}

// ---------------------------------------------------------------------------
// Tests with encoding.TextUnmarshaller conforming type
// ---------------------------------------------------------------------------
func TestTextUnmarshaller(t *testing.T) {
	var customTextUnmarshallerType = reflect.TypeOf((*customTextUnmarshaller)(nil)).Elem()
	require.That(t, value.CanParseType(customTextUnmarshallerType)).Eq(true)

	var v customTextUnmarshaller
	var err = value.Parse(&v, "foobar")
	require.That(t, err).IsNil()
	require.That(t, v).Eq("foobar")
}

type customTextUnmarshaller string

func (v *customTextUnmarshaller) UnmarshalText(text []byte) error {
	*v = customTextUnmarshaller(text)
	return nil
}

// ---------------------------------------------------------------------------
// Tests with non-parsable types
// ---------------------------------------------------------------------------
func TestUnparsable(t *testing.T) {
	var unparsableType = reflect.TypeOf((*unparsable)(nil)).Elem()
	require.That(t, value.CanParseType(unparsableType)).Eq(false)

	var v unparsable
	require.That(t, func() {
		value.Parse(&v, "foobar")
	}).Panics()
}

type unparsable string

// ---------------------------------------------------------------------------
// Tests with built-in and standrad types
// ---------------------------------------------------------------------------
func TestTimeDuration(t *testing.T) {
	var vt = reflect.TypeOf((*time.Duration)(nil)).Elem()
	require.That(t, value.CanParseType(vt)).Eq(true)

	var v time.Duration
	var err = value.Parse(&v, "1s")
	require.That(t, err).IsNil()

	require.That(t, value.Parse(&v, "foobar")).IsNotNil()
}

// ---------------------------------------------------------------------------
// Tests with bad arguments to Parse
// ---------------------------------------------------------------------------

func TestParseArgument(t *testing.T) {
	var v time.Duration
	require.That(t, func() {
		value.Parse(v, "1s")
	}).Panics()

	var pv *time.Duration
	require.That(t, func() {
		value.Parse(pv, "1s")
	}).Panics()

	require.That(t, func() {
		value.Parse(&pv, "1s")
	}).Panics()
}
