package ast

type Node struct {
}
type Document struct {
	*Schema
	Scalars     map[string]ScalarDef
	ObjectTypes map[string]ObjectTypeDef
	Interfaces  map[string]InterfaceDef
	Unions      map[string]UnionDef
	Enums       map[string]EnumDef
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
	Arguments   []ArgumentDef
	Type
	Directives []Directive
}
type ArgumentDef struct {
	Node

	Name string
	Type

	Description  *string
	DefaultValue Value
	Directives   []Directive
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
