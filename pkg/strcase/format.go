package strcase

import (
	"strings"
)

//go:generate go run github.com/maargenton/go-cli/cmd/enumer format.go

type Format uint8

const (
	CamelCase Format = iota
	LowerCamelCase
	NormalizedCamelCase
	NormalizedLowerCamelCase
	SnakeCase
	HyphenCase

	FilteredCamelCase
	FilteredLowerCamelCase
	NormalizedFilteredCamelCase
	NormalizedFilteredLowerCamelCase
	FilteredSnakeCase
	FilteredHyphenCase
)

var SupportedFormats = []Format{
	CamelCase, LowerCamelCase,
	NormalizedCamelCase, NormalizedLowerCamelCase,
	SnakeCase, HyphenCase,
}

var AllFormats = []Format{
	CamelCase, LowerCamelCase,
	NormalizedCamelCase, NormalizedLowerCamelCase,
	SnakeCase, HyphenCase,
	FilteredCamelCase, FilteredLowerCamelCase,
	NormalizedFilteredCamelCase, NormalizedFilteredLowerCamelCase,
	FilteredSnakeCase, FilteredHyphenCase,
}

func (f Format) Apply(name, groupName string) string {
	var parts = Split(name)
	var groupParts = Split(groupName)
	var altParts = FilterParts(parts, groupParts)
	return f.ApplySlice(parts, altParts)
}

func (f Format) ApplySlice(parts, filteredParts []string) string {
	var filter = filters[f]
	return filter.apply(parts, filteredParts)
}

type nameFilter struct {
	Transform      func(string) string
	FirstTransform func(string) string
	Separator      string
	AtlInput       bool
}

func (f *nameFilter) apply(input, alt []string) string {
	var parts = input
	if f.AtlInput && len(alt) > 0 {
		parts = alt
	}

	var p = make([]string, 0, len(parts))
	for i, v := range parts {
		if i == 0 && f.FirstTransform != nil {
			p = append(p, f.FirstTransform(v))
		} else if f.Transform != nil {
			p = append(p, f.Transform(v))
			// } else {
			// 	p = append(p, v)
		}
	}
	return strings.Join(p, f.Separator)
}

func normalizedTitle(s string) string {
	return strings.Title(strings.ToLower(s))
}

var filters = map[Format]nameFilter{
	CamelCase: {
		Transform: strings.Title,
	},
	LowerCamelCase: {
		Transform:      strings.Title,
		FirstTransform: strings.ToLower,
	},
	NormalizedCamelCase: {
		Transform: normalizedTitle,
	},
	NormalizedLowerCamelCase: {
		Transform:      normalizedTitle,
		FirstTransform: strings.ToLower,
	},
	SnakeCase: {
		Transform: strings.ToLower,
		Separator: "_",
	},
	HyphenCase: {
		Transform: strings.ToLower,
		Separator: "-",
	},
	FilteredCamelCase: {
		Transform: strings.Title,
		AtlInput:  true,
	},
	FilteredLowerCamelCase: {
		Transform:      strings.Title,
		FirstTransform: strings.ToLower,
		AtlInput:       true,
	},
	NormalizedFilteredCamelCase: {
		Transform: normalizedTitle,
		AtlInput:  true,
	},
	NormalizedFilteredLowerCamelCase: {
		Transform:      normalizedTitle,
		FirstTransform: strings.ToLower,
		AtlInput:       true,
	},
	FilteredSnakeCase: {
		Transform: strings.ToLower,
		Separator: "_",
		AtlInput:  true,
	},
	FilteredHyphenCase: {
		Transform: strings.ToLower,
		Separator: "-",
		AtlInput:  true,
	},
}
