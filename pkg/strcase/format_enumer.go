// GENERATED CODE -- DO NOT EDIT

package strcase

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maargenton/go-cli/pkg/enumer/enum"
)

// ---------------------------------------------------------------------------
// Format

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CamelCase-(0)]
	_ = x[LowerCamelCase-(1)]
	_ = x[NormalizedCamelCase-(2)]
	_ = x[NormalizedLowerCamelCase-(3)]
	_ = x[SnakeCase-(4)]
	_ = x[HyphenCase-(5)]
	_ = x[FilteredCamelCase-(6)]
	_ = x[FilteredLowerCamelCase-(7)]
	_ = x[NormalizedFilteredCamelCase-(8)]
	_ = x[NormalizedFilteredLowerCamelCase-(9)]
	_ = x[FilteredSnakeCase-(10)]
	_ = x[FilteredHyphenCase-(11)]
}

var _ enum.Type = (*Format)(nil)

var FormatValues = []enum.Value{
	{
		Name:     "camel-case",
		AltNames: []string{"camel-case", "CamelCase", "camelCase", "camel_case"},
		Value:    CamelCase,
	},
	{
		Name:     "lower-camel-case",
		AltNames: []string{"lower-camel-case", "LowerCamelCase", "lowerCamelCase", "lower_camel_case"},
		Value:    LowerCamelCase,
	},
	{
		Name:     "normalized-camel-case",
		AltNames: []string{"normalized-camel-case", "NormalizedCamelCase", "normalizedCamelCase", "normalized_camel_case"},
		Value:    NormalizedCamelCase,
	},
	{
		Name:     "normalized-lower-camel-case",
		AltNames: []string{"normalized-lower-camel-case", "NormalizedLowerCamelCase", "normalizedLowerCamelCase", "normalized_lower_camel_case"},
		Value:    NormalizedLowerCamelCase,
	},
	{
		Name:     "snake-case",
		AltNames: []string{"snake-case", "SnakeCase", "snakeCase", "snake_case"},
		Value:    SnakeCase,
	},
	{
		Name:     "hyphen-case",
		AltNames: []string{"hyphen-case", "HyphenCase", "hyphenCase", "hyphen_case"},
		Value:    HyphenCase,
	},
	{
		Name:     "filtered-camel-case",
		AltNames: []string{"filtered-camel-case", "FilteredCamelCase", "filteredCamelCase", "filtered_camel_case"},
		Value:    FilteredCamelCase,
	},
	{
		Name:     "filtered-lower-camel-case",
		AltNames: []string{"filtered-lower-camel-case", "FilteredLowerCamelCase", "filteredLowerCamelCase", "filtered_lower_camel_case"},
		Value:    FilteredLowerCamelCase,
	},
	{
		Name:     "normalized-filtered-camel-case",
		AltNames: []string{"normalized-filtered-camel-case", "NormalizedFilteredCamelCase", "normalizedFilteredCamelCase", "normalized_filtered_camel_case"},
		Value:    NormalizedFilteredCamelCase,
	},
	{
		Name:     "normalized-filtered-lower-camel-case",
		AltNames: []string{"normalized-filtered-lower-camel-case", "NormalizedFilteredLowerCamelCase", "normalizedFilteredLowerCamelCase", "normalized_filtered_lower_camel_case"},
		Value:    NormalizedFilteredLowerCamelCase,
	},
	{
		Name:     "filtered-snake-case",
		AltNames: []string{"filtered-snake-case", "FilteredSnakeCase", "filteredSnakeCase", "filtered_snake_case"},
		Value:    FilteredSnakeCase,
	},
	{
		Name:     "filtered-hyphen-case",
		AltNames: []string{"filtered-hyphen-case", "FilteredHyphenCase", "filteredHyphenCase", "filtered_hyphen_case"},
		Value:    FilteredHyphenCase,
	},
}

func (v Format) EnumValues() []enum.Value {
	return FormatValues
}

func (v Format) String() string {
	switch v {
	case CamelCase:
		return "camel-case"
	case LowerCamelCase:
		return "lower-camel-case"
	case NormalizedCamelCase:
		return "normalized-camel-case"
	case NormalizedLowerCamelCase:
		return "normalized-lower-camel-case"
	case SnakeCase:
		return "snake-case"
	case HyphenCase:
		return "hyphen-case"
	case FilteredCamelCase:
		return "filtered-camel-case"
	case FilteredLowerCamelCase:
		return "filtered-lower-camel-case"
	case NormalizedFilteredCamelCase:
		return "normalized-filtered-camel-case"
	case NormalizedFilteredLowerCamelCase:
		return "normalized-filtered-lower-camel-case"
	case FilteredSnakeCase:
		return "filtered-snake-case"
	case FilteredHyphenCase:
		return "filtered-hyphen-case"
	}
	return "Format(" + strconv.FormatInt(int64(v), 10) + ")"
}

func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "camel-case", "camelcase", "camel_case":
		return CamelCase, nil
	case "lower-camel-case", "lowercamelcase", "lower_camel_case":
		return LowerCamelCase, nil
	case "normalized-camel-case", "normalizedcamelcase", "normalized_camel_case":
		return NormalizedCamelCase, nil
	case "normalized-lower-camel-case", "normalizedlowercamelcase", "normalized_lower_camel_case":
		return NormalizedLowerCamelCase, nil
	case "snake-case", "snakecase", "snake_case":
		return SnakeCase, nil
	case "hyphen-case", "hyphencase", "hyphen_case":
		return HyphenCase, nil
	case "filtered-camel-case", "filteredcamelcase", "filtered_camel_case":
		return FilteredCamelCase, nil
	case "filtered-lower-camel-case", "filteredlowercamelcase", "filtered_lower_camel_case":
		return FilteredLowerCamelCase, nil
	case "normalized-filtered-camel-case", "normalizedfilteredcamelcase", "normalized_filtered_camel_case":
		return NormalizedFilteredCamelCase, nil
	case "normalized-filtered-lower-camel-case", "normalizedfilteredlowercamelcase", "normalized_filtered_lower_camel_case":
		return NormalizedFilteredLowerCamelCase, nil
	case "filtered-snake-case", "filteredsnakecase", "filtered_snake_case":
		return FilteredSnakeCase, nil
	case "filtered-hyphen-case", "filteredhyphencase", "filtered_hyphen_case":
		return FilteredHyphenCase, nil
	}
	return 0, fmt.Errorf("invalid Format value '%v'", s)
}

func (v *Format) Set(s string) error {
	vv, err := ParseFormat(s)
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

// Format
// ---------------------------------------------------------------------------
