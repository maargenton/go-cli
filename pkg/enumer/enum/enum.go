package enum

type Value struct {
	Name     string
	GoName   string
	AltNames []string
	Value    interface{}
}

type Type interface {
	String() string
	Set(v string) error
	EnumValues() []Value
}
