package strcase_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-cli/pkg/strcase"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestFormat(t *testing.T) {
	var tcs = []struct {
		name, group string
		format      strcase.Format
		output      string
	}{
		{"AccessReadOnly", "AccessLevel", strcase.CamelCase, "AccessReadOnly"},
		{"AccessReadOnly", "AccessLevel", strcase.LowerCamelCase, "accessReadOnly"},
		{"AccessReadOnly", "AccessLevel", strcase.SnakeCase, "access_read_only"},
		{"AccessReadOnly", "AccessLevel", strcase.HyphenCase, "access-read-only"},
		{"AccessReadOnly", "AccessLevel", strcase.FilteredCamelCase, "ReadOnly"},
		{"AccessReadOnly", "AccessLevel", strcase.FilteredLowerCamelCase, "readOnly"},
		{"AccessReadOnly", "AccessLevel", strcase.FilteredSnakeCase, "read_only"},
		{"AccessReadOnly", "AccessLevel", strcase.FilteredHyphenCase, "read-only"},

		{"AccessREADOnly", "AccessLevel", strcase.CamelCase, "AccessREADOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.LowerCamelCase, "accessREADOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.NormalizedCamelCase, "AccessReadOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.NormalizedLowerCamelCase, "accessReadOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.FilteredCamelCase, "READOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.FilteredLowerCamelCase, "readOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.NormalizedFilteredCamelCase, "ReadOnly"},
		{"AccessREADOnly", "AccessLevel", strcase.NormalizedFilteredLowerCamelCase, "readOnly"},
	}

	for _, tc := range tcs {
		bdd.Given(t, fmt.Sprintf("a format target %v", tc.format), func(t *bdd.T) {
			t.When("applying format", func(t *bdd.T) {
				var r = tc.format.Apply(tc.name, tc.group)
				t.Then("something happens", func(t *bdd.T) {
					require.That(t, r).Eq(tc.output)
				})
			})
		})
	}
}
