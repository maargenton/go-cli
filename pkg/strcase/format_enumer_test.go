// GENERATED CODE -- DO NOT EDIT

package strcase_test

import (
	"math"
	"testing"

	"github.com/maargenton/go-cli/pkg/strcase"
)

// ---------------------------------------------------------------------------
// Format

func TestFormatEnummer(t *testing.T) {
	var l = []strcase.Format{
		strcase.CamelCase,
		strcase.LowerCamelCase,
		strcase.NormalizedCamelCase,
		strcase.NormalizedLowerCamelCase,
		strcase.SnakeCase,
		strcase.HyphenCase,
		strcase.FilteredCamelCase,
		strcase.FilteredLowerCamelCase,
		strcase.NormalizedFilteredCamelCase,
		strcase.NormalizedFilteredLowerCamelCase,
		strcase.FilteredSnakeCase,
		strcase.FilteredHyphenCase,
	}

	for _, v := range l {
		var vv strcase.Format
		var err = vv.Set(v.String())
		if err != nil {
			t.Errorf("failed to parse %v", v.String())
		}
		if v != vv {
			t.Errorf("%v != %v", v, vv)
		}
	}

	var v = strcase.Format(math.MaxUint8)
	_ = v.String()
	if len(v.EnumValues()) == 0 {
		t.Errorf("unexpected empty EnumValues()")
	}
	if err := v.Set("--**--some-string-that-should-never-match-anything--??--"); err == nil {
		t.Errorf("Set() with invalid values should generate an error")
	}

	var b, _ = v.MarshalText()
	_ = v.UnmarshalText(b)
}

// Format
// ---------------------------------------------------------------------------
