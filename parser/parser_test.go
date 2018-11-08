package parser_test

import (
	"testing"

	"github.com/dianelooney/graphql/parser"
)

func TestSchemaParser(t *testing.T) {
	src := `
	query Something(simple) {
		field1
		field2 {
			something
			more
			complex(than: that) {
				...isGood
			}
		}
	}
	fragment isGood {
		to
		test
	}
`
	p := &parser.SchemaParser{}
	p.Init([]byte(src))
	_, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

}
