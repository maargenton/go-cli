package enumer

var templateText = `
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

// {{.Name}}
// ---------------------------------------------------------------------------
{{end}}

`
