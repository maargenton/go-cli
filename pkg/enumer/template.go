package enumer

import (
	"text/template"
)

var funcMap = template.FuncMap{}

var enumerTemplate = template.Must(
	template.New("template_enumer.go").
		Funcs(funcMap).
		Parse(enumerTemplateStr))

var enumerTemplateStr = `
// GENERATED CODE -- DO NOT EDIT

package {{.PkgName}}

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maargenton/go-cli/pkg/enumer/enum"
)

{{range .Types}}
// ---------------------------------------------------------------------------
// {{.Name}}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	{{range .Values -}}
	_ = x[{{.GoName}}-({{.Value}})]
	{{end -}}
	{{if .Bitfield -}}
	_ = {{.Name}}(0).bitfield
	{{end -}}
}

var _ enum.Type = (*{{.Name}})(nil)

var {{.Name}}Values = []enum.Value{
	{{range .Values -}}
	{
		Name:     {{.Name | printf "%q"}},
		GoName:   {{.GoName | printf "%q"}},
		AltNames: []string{ {{range .AltNames -}}{{. | printf "%q"}},{{end -}} },
		Value:    {{.GoName}},
	},
	{{end -}}
}

func (v {{.Name}}) EnumValues() []enum.Value {
	return {{.Name}}Values
}


func (v {{.Name}}) String() string {
	switch v {
	{{range .Values -}}
	case {{.GoName}}:
		return {{.Name | printf "%q"}}
	{{ end -}}
	}
	return "{{.Name}}(" + strconv.FormatInt(int64(v), 10) + ")"
}

func Parse{{.Name}} ( s string ) ( {{.Name}}, error ) {
	switch strings.ToLower(s) {
	{{range .Values -}}

	case {{range $i, $v := .LowerCaseNames -}}{{if $i}} ,{{end}}{{$v | printf "%q"}}{{end -}}:
		return {{.GoName}}, nil
	{{ end -}}
	}
	return 0, fmt.Errorf("invalid {{.Name}} value '%v'", s)
}

func (v *{{.Name}}) Set(s string) error {
	vv, err := Parse{{.Name}}(s)
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

func (v {{.Name}}) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

func (v *{{.Name}}) UnmarshalText(text []byte) error {
	return v.Set(string(text))
}

// {{.Name}}
// ---------------------------------------------------------------------------
{{end}}
`

var enumerTestTemplate = template.Must(
	template.New("template_enumer_test.go").
		Funcs(funcMap).
		Parse(enumerTestTemplateStr))

var enumerTestTemplateStr = `
// GENERATED CODE -- DO NOT EDIT

package {{.PkgName}}_test

import (
	"math"
	"testing"
)

{{range .Types}}
// ---------------------------------------------------------------------------
// {{.Name}}

func Test{{.Name}}Enummer(t *testing.T) {
	var l = [] {{ $.PkgName }}.{{ .Name }} {
		{{range .Values -}}
		{{ $.PkgName }}.{{ .GoName }},
		{{end -}}
	}

	for _, v := range l {
		var vv {{ $.PkgName }}.{{ .Name }}
		var err = vv.Set(v.String())
		if err != nil {
			t.Errorf("failed to parse %v", v.String())
		}
		if v != vv {
			t.Errorf("%v != %v", v, vv)
		}
	}

	{{if .UnderlyingMaxIntLiteral -}}
	var v = {{ $.PkgName }}.{{ .Name }}({{.UnderlyingMaxIntLiteral}})
	{{else -}}
	var v {{ $.PkgName }}.{{ .Name }}
	{{end -}}
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

// {{.Name}}
// ---------------------------------------------------------------------------
{{end}}

`
