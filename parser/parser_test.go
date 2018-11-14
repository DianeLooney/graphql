package parser_test

import (
	"strings"
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
	doc := p.Parse()
	for _, sc := range doc.Scalars {
		desc := "<nil>"
		if sc.Description != nil {
			desc = "'" + *sc.Description + "'"
		}
		t.Logf("Scalar %v (description %v, directiveCount %v)", sc.Name, desc, len(sc.Directives))
	}
	for _, obj := range doc.ObjectTypes {
		desc := "<nil>"
		if obj.Description != nil {
			desc = "'" + *obj.Description + "'"
		}
		t.Logf("Object %v (description %v, directiveCount %v)", obj.Name, desc, len(obj.Directives))
	}
	for _, intf := range doc.Interfaces {
		desc := "<nil>"
		if intf.Description != nil {
			desc = "'" + *intf.Description + "'"
		}
		t.Logf("Interface %v (description %v, directiveCount %v)", intf.Name, desc, len(intf.Directives))
	}
	for _, u := range doc.Unions {
		desc := "<nil>"
		if u.Description != nil {
			desc = "'" + *u.Description + "'"
		}
		t.Logf("Union %v (description %v, directiveCount %v, types %v)", u.Name, desc, len(u.Directives), strings.Join(u.Types, " | "))
	}
	for _, e := range doc.Enums {
		desc := "<nil>"
		if e.Description != nil {
			desc = "'" + *e.Description + "'"
		}
		t.Logf("Enum %v (description %v, directiveCount %v)", e.Name, desc, len(e.Directives))
	}
	for _, input := range doc.Inputs {
		desc := "<nil>"
		if input.Description != nil {
			desc = "'" + *input.Description + "'"
		}
		t.Logf("Input %v (description %v, directiveCount %v)", input.Name, desc, len(input.Directives))
	}
	for _, dir := range doc.Directives {
		desc := "<nil>"
		if dir.Description != nil {
			desc = "'" + *dir.Description + "'"
		}
		t.Logf("Directive %v (description %v, onCount %v)", dir.Name, desc, dir.Locations)
	}
	for _, frag := range doc.Fragments {
		t.Logf("Fragment %v", frag.Name)
	}
	for _, err := range p.Errors() {
		t.Error(err)
	}
}
