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

}
