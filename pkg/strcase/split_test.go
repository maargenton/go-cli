package strcase_test

import (
	"fmt"
	"testing"

	"github.com/maargenton/go-cli/pkg/strcase"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestSplit(t *testing.T) {
	var tcs = []struct {
		input  string
		output []string
	}{
		{"camel_case", []string{"camel", "case"}},
		{"camel-case", []string{"camel", "case"}},
		{"camelCase", []string{"camel", "Case"}},
		{"CamelCase", []string{"Camel", "Case"}},
		{"hyphen_case", []string{"hyphen", "case"}},
		{"hyphen-case", []string{"hyphen", "case"}},
		{"hyphenCase", []string{"hyphen", "Case"}},
		{"HyphenCase", []string{"Hyphen", "Case"}},
		{"kebab_case", []string{"kebab", "case"}},
		{"kebab-case", []string{"kebab", "case"}},
		{"kebabCase", []string{"kebab", "Case"}},
		{"KebabCase", []string{"Kebab", "Case"}},
		{"lower_camel_case", []string{"lower", "camel", "case"}},
		{"lower-camel-case", []string{"lower", "camel", "case"}},
		{"lowerCamelCase", []string{"lower", "Camel", "Case"}},
		{"LowerCamelCase", []string{"Lower", "Camel", "Case"}},
		{"snake_case", []string{"snake", "case"}},
		{"snake-case", []string{"snake", "case"}},
		{"snakeCase", []string{"snake", "Case"}},
		{"SnakeCase", []string{"Snake", "Case"}},

		{"SQLStmt", []string{"SQL", "Stmt"}},
		{"CACert", []string{"CA", "Cert"}},

		{"foo bar baz", []string{"foo", "bar", "baz"}},
	}

	bdd.Given(t, "a string", func(t *bdd.T) {
		for _, tc := range tcs {
			t.When(fmt.Sprintf("calling Split(%#v)", tc.input), func(t *bdd.T) {
				var actual = strcase.Split(tc.input)
				t.Then("a properly split string is returned", func(t *bdd.T) {
					require.That(t, actual).Eq(tc.output)
				})
			})
		}
	})
}
