package value_test

import (
	"testing"

	"github.com/maargenton/go-cli/pkg/value"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestParseBool(t *testing.T) {
	var v bool
	var err = value.Parse(&v, "true")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(true)
}

// ---------------------------------------------------------------------------

func TestParseInt(t *testing.T) {
	var v int
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseInt8(t *testing.T) {
	var v int8
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseInt16(t *testing.T) {
	var v int16
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseInt32(t *testing.T) {
	var v int32
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseInt64(t *testing.T) {
	var v int64
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

// ---------------------------------------------------------------------------

func TestParseUInt(t *testing.T) {
	var v uint
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseUInt8(t *testing.T) {
	var v uint8
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseUInt16(t *testing.T) {
	var v uint16
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseUInt32(t *testing.T) {
	var v uint32
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

func TestParseUInt64(t *testing.T) {
	var v uint64
	var err = value.Parse(&v, "123")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123)
}

// ---------------------------------------------------------------------------

func TestParseFloat32(t *testing.T) {
	var v float32
	var err = value.Parse(&v, "123.456")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(float32(123.456))
}

func TestParseFloat64(t *testing.T) {
	var v float64
	var err = value.Parse(&v, "123.456")
	require.That(t, err).IsNil()
	require.That(t, v).Eq(123.456)
}

// ---------------------------------------------------------------------------

func TestParseString(t *testing.T) {
	var v string
	var err = value.Parse(&v, "foobar")
	require.That(t, err).IsNil()
	require.That(t, v).Eq("foobar")
}
