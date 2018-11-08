package ast

type Document struct {
	Definitions []Definition
}
type Definition struct {
	*ExecutableDefinition
	// *TypeSystemDefinition
	// *TypeSystemExtension
}
type ExecutableDefinition struct {
	*OperationDefinition
	*FragmentDefinition
}
type OperationDefinition struct {
	*OperationType
	*Name
	VariableDefinitions []VariableDefinition
	Directives          []Directive
	SelectionSet        []Selection
}
type OperationType int

const (
	QUERY OperationType = iota
	MUTATION
	SUBSCRIPTION
)

type Selection struct {
	*Field
	*FragmentSpread
	*InlineFragment
}

type Field struct {
	*Alias
	Name
	Arguments    []Argument
	Directives   []Directive
	SelectionSet []Selection
}
type Alias struct {
	Name
}
type Name string

type Argument struct {
	Name
	Value
}
type Value struct {
	*Variable
	*IntValue
	*FloatValue
	*StringValue
	*BooleanValue
	*NullValue
	*EnumValue
	*ListValue
	*ObjectValue
}
type BooleanValue struct {
	bool
}
type NullValue struct {
}
type EnumValue struct {
	Name
}
type ListValue struct {
	Values []Value
}
type IntValue struct {
	int
}
type FloatValue struct {
	float64
}
type StringValue struct {
	string
}
type ObjectValue struct {
	ObjectFields []ObjectField
}
type ObjectField struct {
	Name
	Value
}
type VariableDefinition struct {
	Variable
	Type
	DefaultValue *Value
}
type Variable struct {
	Name
}
type Type struct {
	*Name
	List    *Type
	NonNull *Type
}
type Directive struct {
	Name
	Arguments []Argument
}
type FragmentSpread struct {
	Name
	Directives []Directive
}
type InlineFragment struct {
	OnType     *Name
	Directives []Directive
	Selections []Selection
}
type FragmentDefinition struct {
	Name
	OnType     *Name
	Directives []Directive
	Selections []Selection
}
