package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"

	"github.com/fatih/structtag"
)

type Grep struct {
	fs     fs.FS
	output *bufio.Writer
}

func NewGrep(fs fs.FS, output *bufio.Writer) *Grep {
	return &Grep{
		fs:     fs,
		output: output,
	}
}

func (g *Grep) Execute() error {
	return fs.WalkDir(g.fs, ".", func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			panic(err)
		}

		b := &bytes.Buffer{}
		fmt.Fprintf(b, "# ./%s\n", path)

		hasStructWithEnvTag := false
		for _, node := range f.Decls {
			switch node.(type) {

			case *ast.GenDecl:
				genDecl := node.(*ast.GenDecl)
				for _, spec := range genDecl.Specs {
					switch spec.(type) {
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)

						bb := &bytes.Buffer{}
						fmt.Fprintf(bb, "## %s\n", typeSpec.Name.Name)

						hasEnvTag := false
						switch typeSpec.Type.(type) {
						case *ast.StructType:
							structType := typeSpec.Type.(*ast.StructType)
							for _, field := range structType.Fields.List {
								if field.Tag == nil {
									continue
								}

								tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
								if err != nil {
									panic(err)
								}

								envValue := ""
								envTag, err := tags.Get("env")
								if err != nil || envTag == nil || envTag.Value() == "" {
									continue
								}
								envValue = envTag.Value()
								if strings.Index(envValue, ",") >= 0 {
									envValue = strings.Split(envValue, ",")[0]
								}

								envDefaultValue := ""
								envDefaultTag, err := tags.Get("envDefault")
								if err == nil {
									envDefaultValue = envDefaultTag.Value()
								}

								fmt.Fprintf(bb, "%s=%s\n", envValue, envDefaultValue)
								hasEnvTag = true
							}
						}

						if hasEnvTag {
							b.WriteString(bb.String())
							hasStructWithEnvTag = true
						}
					}
				}
			}
		}

		if hasStructWithEnvTag {
			g.output.WriteString(b.String())
			g.output.Flush()
		}

		return nil
	})
}
