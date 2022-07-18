package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(os.Stdout)

	filepath.WalkDir(cwd, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			panic(err)
		}

		b := &bytes.Buffer{}
		fmt.Fprintf(b, "# .%s\n", strings.TrimPrefix(path, cwd))

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

								envTag, err := tags.Get("env")
								if err != nil || envTag == nil || envTag.Value() == "" {
									continue
								}
								envDefaultTag, err := tags.Get("envDefault")
								if err != nil || envTag == nil {
									continue
								}

								fmt.Fprintf(bb, "%s=%s\n", envTag.Value(), envDefaultTag.Value())
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
			w.WriteString(b.String())
			w.Flush()
		}

		return nil
	})

}
