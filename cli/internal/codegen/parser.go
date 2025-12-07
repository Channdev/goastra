/*
 * GoAstra CLI - Go Parser
 *
 * Parses Go source files to extract struct definitions
 * for TypeScript code generation.
 */
package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

/*
 * GoParser extracts type definitions from Go source files.
 */
type GoParser struct {
	schemaPath string
}

/*
 * TypeDef represents a parsed Go type definition.
 */
type TypeDef struct {
	Name   string
	Fields []FieldDef
	Doc    string
}

/*
 * FieldDef represents a struct field definition.
 */
type FieldDef struct {
	Name     string
	Type     string
	JSONName string
	Optional bool
	Doc      string
}

/*
 * NewGoParser creates a new parser instance.
 */
func NewGoParser(schemaPath string) *GoParser {
	return &GoParser{schemaPath: schemaPath}
}

/*
 * Parse reads all Go files in the schema path and extracts types.
 */
func (p *GoParser) Parse() ([]TypeDef, error) {
	var types []TypeDef

	files, err := filepath.Glob(filepath.Join(p.schemaPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob schema files: %w", err)
	}

	fset := token.NewFileSet()

	for _, file := range files {
		fileDefs, err := p.parseFile(fset, file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}
		types = append(types, fileDefs...)
	}

	return types, nil
}

func (p *GoParser) parseFile(fset *token.FileSet, filename string) ([]TypeDef, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var types []TypeDef

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			typeDef := TypeDef{
				Name: typeSpec.Name.Name,
			}

			if genDecl.Doc != nil {
				typeDef.Doc = genDecl.Doc.Text()
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				fieldDef := FieldDef{
					Name: field.Names[0].Name,
					Type: p.typeToString(field.Type),
				}

				if field.Tag != nil {
					jsonTag := p.extractJSONTag(field.Tag.Value)
					if jsonTag != "" && jsonTag != "-" {
						fieldDef.JSONName = strings.Split(jsonTag, ",")[0]
						fieldDef.Optional = strings.Contains(jsonTag, "omitempty")
					}
				}

				if field.Doc != nil {
					fieldDef.Doc = field.Doc.Text()
				}

				typeDef.Fields = append(typeDef.Fields, fieldDef)
			}

			types = append(types, typeDef)
		}
	}

	return types, nil
}

func (p *GoParser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", p.typeToString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + p.typeToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.typeToString(t.Key), p.typeToString(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

func (p *GoParser) extractJSONTag(tag string) string {
	tag = strings.Trim(tag, "`")
	structTag := reflect.StructTag(tag)
	return structTag.Get("json")
}
