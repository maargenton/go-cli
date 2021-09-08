package value

import (
	"strconv"
	"time"
)

func init() {

	RegisterParser(parseBool)
	RegisterParser(parseInt, parseInt8, parseInt16, parseInt32, parseInt64)
	RegisterParser(parseUint, parseUint8, parseUint16, parseUint32, parseUint64)
	RegisterParser(parseFloat32, parseFloat64)

	RegisterParser(parseString)
	RegisterParser(time.ParseDuration)
}

// ---------------------------------------------------------------------------

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ---------------------------------------------------------------------------

func parseInt(s string) (int, error) {
	v, err := strconv.ParseInt(s, 0, 0)
	return int(v), err
}

func parseInt8(s string) (int8, error) {
	v, err := strconv.ParseInt(s, 0, 8)
	return int8(v), err
}

func parseInt16(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 0, 16)
	return int16(v), err
}

func parseInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 0, 32)
	return int32(v), err
}

func parseInt64(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	return int64(v), err
}

// ---------------------------------------------------------------------------

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 0, 0)
	return uint(v), err
}

func parseUint8(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 0, 8)
	return uint8(v), err
}

func parseUint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 0, 16)
	return uint16(v), err
}

func parseUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 0, 32)
	return uint32(v), err
}

func parseUint64(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 0, 64)
	return uint64(v), err
}

// ---------------------------------------------------------------------------

func parseFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

func parseFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	return float64(v), err
}

// ---------------------------------------------------------------------------

func parseString(s string) (string, error) {
	return s, nil
}
