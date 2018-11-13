package parser_test

import (
	"testing"

	parser "github.com/dianelooney/graphql/parsers"
)

func TestSchemaParser(t *testing.T) {

	src := `
	schema {
		query: apples
		mutation: bananas
		subscription: carrots
	}
	`
	p := parser.Parser{}
	p.Init([]byte(src))
	s := p.ParseSchema()
	t.Log(s.Directives)
	t.Log(*s.Mutation)
	t.Log(*s.Query)
	t.Log(*s.Subscription)
}
