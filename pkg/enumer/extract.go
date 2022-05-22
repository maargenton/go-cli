package enumer

import (
	"go/constant"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/maargenton/go-cli/pkg/strcase"
)

func ExtractEnums(pkg *packages.Package) (*PkgEnums, error) {
	var enums = &PkgEnums{
		PkgName: pkg.Name,
		PkgPath: pkg.PkgPath,
	}

	var pkgTypes = make(map[*types.Named][]*types.Const)
	var scope = pkg.Types.Scope()
	var names = scope.Names()

	for _, name := range names {
		var def = scope.Lookup(name)
		if def, ok := def.(*types.Const); ok {
			if t, ok := def.Type().(*types.Named); ok {
				var p = t.Obj().Pkg()
				if p.Path() != pkg.PkgPath {
					continue
				}
				pkgTypes[t] = append(pkgTypes[t], def)
			}
		}
	}

	for t, vv := range pkgTypes {
		sort.Slice(vv, func(i, j int) bool {
			return vv[i].Pos() < vv[j].Pos()
		})

		var hasCustomFormat = false
		var hasParseCustomFormat = false
		var isBitfield = false

		var n = t.NumMethods()
		for i := 0; i < n; i++ {
			switch t.Method(i).Name() {
			case "customFormat":
				hasCustomFormat = true
			case "parseCustomFormat":
				hasParseCustomFormat = true
			case "bitfield":
				isBitfield = true
			}
		}

		var defs = make([]EnumValueDef, 0, len(vv))
		for _, v := range vv {
			if v.Val().Kind() != constant.Int {
				continue
			}
			var value = EnumValueDef{
				Name: v.Name(),
				def:  v,
			}

			if v, ok := constant.Int64Val(v.Val()); ok {
				value.Value = v
			}
			if v, ok := constant.Uint64Val(v.Val()); ok {
				value.Value = v
			}
			defs = append(defs, value)
		}

		enums.Types = append(enums.Types, EnumType{
			Name:         t.Obj().Name(),
			Definitions:  defs,
			CustomFormat: hasCustomFormat && hasParseCustomFormat,
			Bitfield:     isBitfield,
			Position:     pkg.Fset.Position(t.Obj().Pos()),
			def:          t,
		})
	}

	return enums, nil
}

func ExtractValues(t *EnumType, f strcase.Format) error {
	var m = make(map[interface{}][]EnumValueDef)
	var v []interface{}
	for _, d := range t.Definitions {
		if _, ok := m[d.Value]; !ok {
			v = append(v, d.Value)
		}
		m[d.Value] = append(m[d.Value], d)
	}

	var typeParts = strcase.Split(t.Name)

	for _, vv := range v {
		defs := m[vv]
		value := EnumValue{
			GoName: defs[0].Name,
			Value:  vv,
		}

		var parts = strcase.Split(value.GoName)
		var altParts = strcase.FilterParts(parts, typeParts)

		value.Name = f.ApplySlice(parts, altParts)
		value.AltNames = append(value.AltNames, value.Name)

		for _, ff := range strcase.AllFormats {
			value.AltNames = append(value.AltNames, ff.ApplySlice(parts, altParts))
		}

		for i, dd := range defs {
			if i == 0 {
				continue
			}

			var parts = strcase.Split(dd.Name)
			var altParts = strcase.FilterParts(parts, typeParts)
			for _, ff := range strcase.AllFormats {
				value.AltNames = append(value.AltNames, ff.ApplySlice(parts, altParts))
			}
		}
		value.AltNames = strcase.UniqueStrings(value.AltNames)

		for _, n := range value.AltNames {
			value.LowerCaseNames = append(value.LowerCaseNames, strings.ToLower(n))
		}
		value.LowerCaseNames = strcase.UniqueStrings(value.LowerCaseNames)

		t.Values = append(t.Values, value)
	}
	return nil
}
