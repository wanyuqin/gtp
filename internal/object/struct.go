package object

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Kind = string

const (
	Invalid Kind = "Invalid"
	Bool         = "bool"
	Int32        = "int32"
	Int64        = "int64"
	Uint32       = "uint32"
	Uint64       = "uint64"
	Float32      = "float"
	Float64      = "double"
	Array        = "repeated"
	Map          = "map"
	Slice        = "repeated"
	Bytes        = "bytes"
	String       = "string"
	Struct       = "struct"
)

var ClassKind = map[string]Kind{
	"string":  String,
	"int":     Int32,
	"int8":    Int32,
	"int16":   Int32,
	"int32":   Int32,
	"int64":   Int64,
	"uint":    Uint32,
	"uint8":   Uint32,
	"uint16":  Uint32,
	"uint32":  Uint32,
	"uint64":  Uint64,
	"slice":   Slice,
	"struct":  Struct,
	"map":     Map,
	"float32": Float32,
	"float64": Float64,
	"bool":    Bool,
}

type Object struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Class    Class  `json:"class"`
	JsonName string `json:"json_name"`
}

type Class struct {
	Kind    Kind `json:"kind"`
	Key     Kind `json:"key"`
	Value   Kind `json:"value"`
	Element Kind `json:"element"`
}

func ParseGo(fname string) (ObjectList, error) {
	objects := make([]Object, 0)
	f, err := parser.ParseFile(token.NewFileSet(), fname, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, v := range f.Decls {
		gendecl := genDeclTransform(v)
		if gendecl == nil {
			continue
		}
		for _, spec := range gendecl.Specs {
			switch spec.(type) {
			case *ast.TypeSpec:
				typeSpec := spec.(*ast.TypeSpec)
				tst := typeSpec.Type
				obj := Object{
					Name: typeSpec.Name.Name,
				}
				switch tst.(type) {
				case *ast.StructType:
					structType := tst.(*ast.StructType)
					fields := parseStructField(structType)
					obj.Fields = fields
				}

				objects = append(objects, obj)
			}

		}

	}

	return objects, nil
}

// parseStructField 处理字段
func parseStructField(st *ast.StructType) []Field {
	fields := make([]Field, 0)
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			f := Field{
				Name:    formatName(name.Name),
				Comment: formatComment(field.Comment.Text()),
			}
			c := Class{}
			switch field.Type.(type) {
			case *ast.Ident:
				if field.Type.(*ast.Ident).Name == "byte" {
					fmt.Println("byte cannot transform proto type")
					continue
				}
				kind, ok := ClassKind[field.Type.(*ast.Ident).Name]
				if !ok {
					fmt.Printf("%s not suppot proto type \n", field.Type.(*ast.Ident).Name)
					continue
				}
				c.Kind = kind
			case *ast.ArrayType:
				c.Kind = Slice
				if field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name == "byte" {
					c.Kind = Bytes
					goto end
				}
				if element, ok := ClassKind[field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name]; ok {
					c.Element = element
				} else {
					c.Element = field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name
				}
			case *ast.MapType:
				c.Kind = Map
				c.Key = ClassKind[field.Type.(*ast.MapType).Key.(*ast.Ident).Name]
				if value, ok := ClassKind[field.Type.(*ast.MapType).Value.(*ast.Ident).Name]; ok {
					c.Value = value
				} else {
					c.Value = field.Type.(*ast.MapType).Value.(*ast.Ident).Name
				}
			case *ast.InterfaceType:
				fmt.Println("interface cannot transform proto type")
				continue

			}

			f.JsonName = formatJsonTag(field.Tag)

		end:
			f.Class = c
			fields = append(fields, f)
		}
	}
	return fields
}

func formatJsonTag(tag *ast.BasicLit) string {
	t := ""
	index := strings.Index(tag.Value, "json:")
	if index >= 0 {
		index += len("json:") + 1
		tag.Value = tag.Value[index:]
		b := make([]byte, 0, len(tag.Value))
		for i := range tag.Value {
			if tag.Value[i] == '"' {
				break
			}
			b = append(b, tag.Value[i])
		}
		t = string(b)
	}

	return t
}

func formatName(n string) string {
	b := make([]byte, 0, len(n))
	for i := range n {
		if i == 0 {
			lower, _ := toLower(n[i])
			b = append(b, lower)
			continue
		}

		if lower, ok := toLower(n[i]); ok {
			b = append(b, '_', lower)
			continue
		}
		b = append(b, n[i])

	}
	return string(b)
}

func formatComment(comment string) string {
	return strings.ReplaceAll(comment, "\n", "")
}
func toLower(b byte) (byte, bool) {
	if b > 'A' && b < 'Z' {
		return b + 32, true
	}
	return b, false
}

func genDeclTransform(decl ast.Decl) *ast.GenDecl {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil
	}

	return genDecl
}

type ObjectList []Object

type Message struct {
	Name          string  `json:"name"`
	PathName      string  `json:"path_name"`
	Fields        []Field `json:"fields"`
	CreateRequest []Field `json:"create_request"`
	UpdateRequest []Field `json:"update_request"`
	DeleteRequest []Field `json:"delete_request"`
	GetRequest    []Field `json:"get_request"`
	ListRequest   []Field `json:"list_request"`
}

func (o ObjectList) BuildParam() []Message {
	msgs := make([]Message, 0, len(o))

	for _, v := range o {
		m := Message{
			Name:          v.Name,
			PathName:      formatPathName(v.Name),
			CreateRequest: v.Fields,
			UpdateRequest: v.Fields,
			Fields:        v.Fields,
		}
		msgs = append(msgs, m)
	}
	return msgs
}

func formatPathName(n string) string {
	b := make([]byte, 0, len(n))
	for i := range n {
		if i == 0 {
			lower, _ := toLower(n[i])
			b = append(b, lower)
			continue
		}

		if lower, ok := toLower(n[i]); ok {
			b = append(b, '-', lower)
			continue
		}
		b = append(b, n[i])

	}
	return string(b)
}

func (o *Object) BuildCreateOrUpdate() {

}
