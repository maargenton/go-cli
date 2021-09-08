package option

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestParseOptsTag(t *testing.T) {
	// assert := asserter.New(t)

	var o = T{}
	var err = o.parseOptsTag(`-d, --db, env: DB, default: admin:admin@tcp(localhist:3306)/test, sep:;`)

	require.That(t, err).IsError(nil)
	require.That(t, o.Short).Eq("d")
	require.That(t, o.Long).Eq("db")
	require.That(t, o.Env).Eq("DB")
	require.That(t, o.Default).Eq("admin:admin@tcp(localhist:3306)/test")
	require.That(t, o.Sep).Eq(";")
}

func TestParseOptsTagPartial(t *testing.T) {
	var o = T{}
	var err = o.parseOptsTag(`--proxy,env:PROXY,sep: \,`)

	require.That(t, err).IsError(nil)
	require.That(t, o.Short).Eq("")
	require.That(t, o.Long).Eq("proxy")
	require.That(t, o.Env).Eq("PROXY")
	require.That(t, o.Default).Eq("")
	require.That(t, o.Sep).Eq(",")
}

func TestParseOptsTagInvalid(t *testing.T) {
	var o = T{}
	var err = o.parseOptsTag(`--proxy,env:PROXY,omitempty`)

	require.That(t, err).IsNotNil()
	require.That(t, err).ToString().Contains("omitempty")
}

func TestSplitSliceValues(t *testing.T) {
	var tcs = []struct {
		input  string
		delim  string
		output []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"a\\,b,c", ",", []string{"a,b", "c"}},
		{"a,b,\\c", ",", []string{"a", "b", "c"}},
		{"a\\\\,b,c", ",", []string{"a\\", "b", "c"}},
		{"a\\\\:b;c", ",:;", []string{"a\\", "b", "c"}},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			var values = splitSliceValues(tc.input, tc.delim)
			require.That(t, values).Eq(tc.output)
		})
	}
}

func TestUnescapeField(t *testing.T) {
	var tcs = []struct {
		in, out string
	}{
		{"ab\\c", "abc"},
		{"abc\\\\", "abc\\"},
		{"abc\\", "abc"},
		{`abc\\def`, `abc\def`},
	}

	for _, tc := range tcs {
		t.Run(tc.in, func(t *testing.T) {
			require.That(t, unescapeField(tc.in)).Eq(tc.out)
		})
	}
}
