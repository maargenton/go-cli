package enumer

import (
	"go/token"
	"go/types"
)

// PkgEnums records all enumerated type found in one package
type PkgEnums struct {
	PkgName string
	PkgPath string
	Types   []EnumType
}

// EnumType records the details about a defined enum type and its associated
// values
type EnumType struct {
	Name         string
	Bitfield     bool
	CustomFormat bool
	Definitions  []EnumValueDef
	Position     token.Position
	def          *types.Named

	Values []EnumValue
}

// EnumValueDef records the definition of an enumerated value
type EnumValueDef struct {
	Name  string
	Value interface{}
	def   *types.Const
}

// EnumValue records all the details of an enumerated value necessary to produce
// the generated code output, including all alternate names for the definition
// of one value.
type EnumValue struct {
	GoName string
	Value  interface{}

	Name           string
	AltNames       []string
	LowerCaseNames []string
}
