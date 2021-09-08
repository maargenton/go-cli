package option

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestScanTagFields(t *testing.T) {

	var tcs = []struct {
		name   string
		input  string
		output []string
	}{
		{
			name:  "short and long",
			input: `-v,--verbose`,
			output: []string{
				"-v", "",
				"--verbose", "",
			},
		},
		{
			name:  "short and long with spaces",
			input: `-v, --verbose`,
			output: []string{
				"-v", "",
				"--verbose", "",
			},
		},
		{
			name:  "with env default",
			input: `-d,--db,env:DB,default:admin:admin@tcp(localhist:3306)/test`,
			output: []string{
				"-d", "",
				"--db", "",
				"env", "DB",
				"default", "admin:admin@tcp(localhist:3306)/test",
			},
		},
		{
			name:  "with env default and spaces",
			input: `-d, --db, env: DB, default: admin:admin@tcp(localhist:3306)/test`,
			output: []string{
				"-d", "",
				"--db", "",
				"env", "DB",
				"default", "admin:admin@tcp(localhist:3306)/test",
			},
		},
		{
			name:  "with escape",
			input: `-p, --proxy, env: PROXY, delim: \,, default: localhost:8080\,localhost:3000\`,
			output: []string{
				"-p", "",
				"--proxy", "",
				"env", "PROXY",
				"delim", ",",
				"default", "localhost:8080,localhost:3000",
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			var s = tc.input
			var result []string
			for len(s) > 0 {
				var k, v string
				k, v, s = scanTagFields(s)
				result = append(result, k, v)
			}

			require.That(t, result).Eq(tc.output)
		})
	}
}
