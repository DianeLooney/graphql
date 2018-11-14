package parser_test

import (
	"testing"

	parser "github.com/dianelooney/graphql/parser"
)

func TestSchemaParser(t *testing.T) {

	src := `
	schema {
		query: Query
		mutation: Mutation
		subscription: Subscription
	}

	scalar Sc

	"scalar desc"
	scalar ScWithDescription

	scalar ScWithDirectives
	@someDirective
	@someDirective2

	type Query {}
	type Mutation {}
	type Subscription {}

	type Obj {}

	interface Interface {}
	"interface description" interface InterfaceWithDescription {}
	union Union = | T1
	"union description" union UnionWithDescription = T1
	union Union2 = T1 | T2 | T3

	enum Enum
	"enum description" enum EnumWithDescription
	enum EnumWithValues {
		Value
		"value description" ValueWithDescription
	}

	input Input
	"input description" input InputWithDescription
	input InputWithValues {
		x: Int
		"description of input value"
		y: String
		z: int = 0
	}

	directive @Directive on SCALAR
	"directive description" directive @DirectiveWithDescription on ENUM
	directive @DirectiveOnMulti on UNION | INPUT_OBJECT | FIELD

	fragment Frag1 on Obj {
		x
	}
	fragment Frag2 on Obj {
		y(r: "s") {
			t
			l: n
			e
		}
	}
	`
	p := parser.Parser{}
	p.Init([]byte(src))
	p.Parse()
	for _, err := range p.Errors() {
		t.Error(err)
	}
}
