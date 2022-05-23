package enumer

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/maargenton/go-fileutils"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/maargenton/go-cli/pkg/strcase"
)

type Cmd struct {
	Format strcase.Format `opts:"-f,--format, default:filtered-hyphen-case" desc:"default string representation format"`
	Files  []string       `opts:"args, name:file" desc:"restrict processing to specified files if any"`
}

func (opts *Cmd) Run() error {
	pkg, err := loadPackage(".")
	if err != nil {
		return err
	}

	enums, err := ExtractEnums(pkg)
	if err != nil {
		return err
	}

	if len(opts.Files) == 0 {
		var files []string
		for _, t := range enums.Types {
			f := filepath.Base(t.Position.Filename)
			files = append(files, f)
		}
		opts.Files = strcase.UniqueStrings(files)
	}

	for _, f := range opts.Files {
		of := strings.ReplaceAll(f, ".go", "_enumer.go")
		otf := strings.ReplaceAll(f, ".go", "_enumer_test.go")
		var data = &PkgEnums{
			PkgName: enums.PkgName,
			PkgPath: enums.PkgPath,
		}

		for _, t := range enums.Types {
			if filepath.Base(t.Position.Filename) == f {
				ExtractValues(&t, opts.Format)
				data.Types = append(data.Types, t)
			}
		}

		if len(data.Types) == 0 {
			return fmt.Errorf("no enum type found in %v", f)
		}

		if err := applyTemplate(of, enumerTemplate, data); err != nil {
			return err
		}
		if data.PkgName != "main" {
			if err := applyTemplate(otf, enumerTestTemplate, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadPackage loads and parses the source code of a package names relative to
// the current working directory
func loadPackage(name string) (pkg *packages.Package, err error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedModule,
		Tests:      false,
		BuildFlags: []string{},
	}

	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return nil, fmt.Errorf("failed to load package '%v': %w", name, err)
	}
	if len(pkgs) < 1 {
		return nil, fmt.Errorf("no package loaded")
	}
	return pkgs[0], nil
}

func applyTemplate(of string, tmpl *template.Template, data interface{}) error {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fileutils.WriteFile(of, func(w io.Writer) error {
		output, err := imports.Process("", buf.Bytes(), &imports.Options{
			AllErrors: true, Comments: true, TabIndent: true, TabWidth: 8,
		})
		if err != nil {
			fmt.Printf("failed to format generated output: %v\n", err)
			fmt.Printf("'%v' saved unformatted\n", of)
			w.Write(buf.Bytes())
		} else {
			w.Write(output)
		}
		return nil
	})

	return nil
}
