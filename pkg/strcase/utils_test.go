package strcase_test

import (
	"testing"

	"github.com/maargenton/go-cli/pkg/strcase"

	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
)

func TestFilterParts(t *testing.T) {
	bdd.Given(t, "split sample enum type and value", func(t *bdd.T) {
		var v = strcase.Split("AccessReadOnly")
		var g = strcase.Split("AccessLevel")
		t.When("calling FilterParts()", func(t *bdd.T) {
			var r = strcase.FilterParts(v, g)
			t.Then("common prefix word s dropped", func(t *bdd.T) {
				require.That(t, r).Eq([]string{"Read", "Only"})

			})
		})
	})
}
func TestUniqueStrings(t *testing.T) {
	bdd.Given(t, "a list of strings with duplicates", func(t *bdd.T) {
		var v = []string{"aaa", "bbb", "ccc", "bbb", "aaa"}
		t.When("calling UniqueStrings()", func(t *bdd.T) {
			var r = strcase.UniqueStrings(v)
			t.Then("the order is preserved with duplicates removed", func(t *bdd.T) {
				require.That(t, r).Eq([]string{"aaa", "bbb", "ccc"})
			})
		})
	})
}
