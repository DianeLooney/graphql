package ast

type Node struct {
}
type Document struct {
	Schema      *Schema
	Scalars     map[string]ScalarDef
	ObjectTypes map[string]ObjectTypeDef
	Interfaces  map[string]InterfaceDef
	Unions      map[string]UnionDef
	Enums       map[string]EnumDef
	Inputs      map[string]InputDef
	Directives  map[string]DirectiveDef
}
type Schema struct {
	Node

	Directives            []Directive
	RootOperationTypeDefs []RootOperationTypeDef
}
type RootOperationTypeDef struct {
	Node

	OpType    string
	NamedType string
}
type Directive struct {
	Node

	Name string
	Arguments
}
type Arguments map[string]Value
type Value struct {
	Node

	Variable *string
	Int      *int
	Float    *float64
	String   *string
	Bool     *bool
	IsNull   bool
	Enum     *string
	List     []Value
	Object   map[string]Value
}
type ScalarDef struct {
	Node

	Name        string
	Description *string
	Directives  []Directive
}
type ObjectTypeDef struct {
	Node

	Name                string
	Description         *string
	ImplementsInterface []string
	Directives          []Directive
	Fields              []FieldDef
}
type FieldDef struct {
	Node

	Name        string
	Description *string
	Arguments   []InputValueDef
	Type
	Directives []Directive
}

type Type struct {
	Name        *string
	ListType    *Type
	NonNullType *Type
}

type InterfaceDef struct {
	Node

	Description *string
	Name        string
	Directives  []Directive
	Fields      []FieldDef
}

type UnionDef struct {
	Node

	Description *string
	Name        string
	Directives  []Directive
	Types       []string
}
type EnumDef struct {
	Node

	Description *string
	Name        string
	Directives  []Directive
	Values      []EnumValueDef
}
type EnumValueDef struct {
	Node

	Description *string
	Name        string
	Directives  []Directive
}
type InputDef struct {
	Node

	Description *string
	Name        string
	Directives  []Directive
	Fields      []InputValueDef
}
type InputValueDef struct {
	Node

	Description *string
	Name        string
	Type
	DefaultValue *Value
	Directives   []Directive
}

var ExecutableDirectiveLocations = map[string]bool{
	"QUERY":               true,
	"MUTATION":            true,
	"SUBSCRIPTION":        true,
	"FIELD":               true,
	"FRAGMENT_DEFINITION": true,
	"FRAGMENT_SPREAD":     true,
	"INLINE_FRAGMENT":     true,
	"VARIABLE_DEFINITION": true,
}
var TypeSystemDirectiveLocations = map[string]bool{
	"SCHEMA":                 true,
	"SCALAR":                 true,
	"OBJECT":                 true,
	"FIELD_DEFINITION":       true,
	"ARGUMENT_DEFINITION":    true,
	"INTERFACE":              true,
	"UNION":                  true,
	"ENUM":                   true,
	"ENUM_VALUE":             true,
	"INPUT_OBJECT":           true,
	"INPUT_FIELD_DEFINITION": true,
}

type DirectiveDef struct {
	Node

	Description *string
	Name        string
	Arguments   []InputValueDef
	Locations   []string
}
