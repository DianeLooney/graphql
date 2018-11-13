package parser

import (
	"errors"
	"strconv"

	"github.com/dianelooney/graphql/ast"
	"github.com/dianelooney/graphql/scanner"
)

type Parser struct {
	sc     scanner.Scanner
	errors []error
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) Init(src []byte) {
	p.sc = scanner.Scanner{}
	p.sc.Init(src)
}

func (p *Parser) Parse() (doc ast.Document) {
	doc.Scalars = make([]ast.ScalarDef, 0)
	for {
		var desc *string
		_, tkn, lit := p.sc.Peek()
		if tkn == scanner.EOF {
			break
		}
		desc = p.parseDescription()
		if tkn == scanner.NAME {
			switch lit {
			case "schema":
				if desc != nil {
					p.errors = append(p.errors, errors.New("unexpected description string given on schema"))
				}
				schema := p.parseSchema()
				doc.Schema = &schema
			case "scalar":
				scalar := p.parseScalarTypeDefinition(desc)
				doc.Scalars = append(doc.Scalars, scalar)
			case "type":
				obj := p.parseObjectTypeDefinition(desc)
				doc.ObjectTypes = append(doc.ObjectTypes, obj)
			}
		}
	}

	return
}

func (p *Parser) parseSchema() (schema ast.Schema) {
	_, tkn, lit := p.sc.Scan()
	if tkn != scanner.NAME || lit != "schema" {
		p.errors = append(p.errors, errors.New("expected token 'schema'"))
	}
	schema.Directives = p.parseDirectives()
	_, tkn, lit = p.sc.Peek()
	if tkn != scanner.LCURLY {
		p.errors = append(p.errors, errors.New("expected a block definining root types, it was: "+lit))
		return
	}
	p.sc.Scan()
	for {
		_, tkn, _ := p.sc.Peek()
		if tkn == scanner.RCURLY {
			break
		}
		if tkn == scanner.EOF {
			p.errors = append(p.errors, errors.New("unexpected EOF"))
			break
		}

		_, tkn, opName := p.sc.Scan()
		if tkn != scanner.NAME {
			p.errors = append(p.errors, errors.New("expected a name"))
			continue
		}

		if opName != "query" &&
			opName != "mutation" &&
			opName != "subscription" {
			p.errors = append(p.errors, errors.New("expected name to be query|mutation|subscription"))
		}
		_, tkn, _ = p.sc.Peek()
		if tkn != scanner.COLON {
			p.errors = append(p.errors, errors.New("expected a colon"))
			continue
		}
		p.sc.Scan()

		_, tkn, namedType := p.sc.Peek()
		if tkn != scanner.NAME {
			p.errors = append(p.errors, errors.New("expected a name"))
			continue
		}
		p.sc.Scan()

		switch opName {
		case "query":
			schema.Query = &namedType
		case "mutation":
			schema.Mutation = &namedType
		case "subscription":
			schema.Subscription = &namedType
		}
	}

	return
}
func (p *Parser) parseScalarTypeDefinition(desc *string) (scalar ast.ScalarDef) {
	scalar.Description = desc
	p.sc.Scan()
	_, tknName, litName := p.sc.Scan()
	if tknName != scanner.NAME {
		p.errors = append(p.errors, errors.New("expected a name for the scalar"))
	}
	scalar.Name = litName
	scalar.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseObjectTypeDefinition(desc *string) (obj ast.ObjectTypeDef) {
	obj.Description = desc
	p.sc.Scan()
	_, nameTkn, name := p.sc.Scan()
	if nameTkn != scanner.NAME {
		p.errors = append(p.errors, errors.New("expected a name for object type"))
	}
	obj.Name = name

	args := p.parseArgumentsDefn()
	return
}
func (p *Parser) parseArgumentsDefn() (args []ast.ArgumentDef) {
	_, tkn, _ := p.sc.Peek()
	if tkn != scanner.LPAREN {
		return
	}

	for {
		arg := ast.ArgumentDef{}
		arg.Description = p.parseDescription()
		if tkn == scanner.EOF {
			p.errors = append(p.errors, errors.New("unexpected EOF"))
			return
		}

		//_, tkn, lit := p.sc.Scan()

	}
}
func (p *Parser) parseDescription() (description *string) {
	_, tkn, _ := p.sc.Peek()

	if tkn == scanner.STRING {
		_, _, lit := p.sc.Scan()
		val := lit[1 : len(lit)-1]
		description = &val
	}
	if tkn == scanner.BLOCKSTRING {
		_, _, lit := p.sc.Scan()
		val := lit[3 : len(lit)-3]
		description = &val
	}
	return
}
func (p *Parser) parseDirectives() (directives []ast.Directive) {
	for {
		_, tkn, _ := p.sc.Peek()
		if tkn != scanner.AT {
			return
		}

		directives = append(directives, p.parseDirective())
	}
}
func (p *Parser) parseDirective() (directive ast.Directive) {
	p.sc.Scan()
	_, nameToken, name := p.sc.Scan()
	if nameToken != scanner.NAME {
		p.errors = append(p.errors, errors.New("expected directive name"))
	}
	directive.Name = name

	_, next, _ := p.sc.Peek()
	if next == scanner.LPAREN {
		directive.Arguments = p.parseArguments()
	}

	return
}
func (p *Parser) parseArguments() (arguments map[string]ast.Value) {
	arguments = make(map[string]ast.Value)

	_, tkn, _ := p.sc.Peek()
	if tkn == scanner.LPAREN {
		p.errors = append(p.errors, errors.New("expected left paren to start argument list"))
		return
	}
	p.sc.Scan()

	for {
		_, tkn, _ := p.sc.Peek()
		if tkn == scanner.RPAREN {
			p.sc.Scan()
			return
		}
		if tkn == scanner.EOF {
			p.errors = append(p.errors, errors.New("unexpected EOF"))
		}

		name, value := p.parseArgument()
		arguments[name] = value
	}
}
func (p *Parser) parseArgument() (name string, value ast.Value) {
	_, tkn, name := p.sc.Scan()
	if tkn != scanner.NAME {
		p.errors = append(p.errors, errors.New("expected argument name"))
	}
	if _, tkn, _ := p.sc.Scan(); tkn != scanner.COLON {
		p.errors = append(p.errors, errors.New("expected argument colon"))
	}
	value = p.parseValue()

	return
}
func (p *Parser) parseObjectField() (name string, value ast.Value) {
	_, tkn, name := p.sc.Scan()
	if tkn != scanner.NAME {
		p.errors = append(p.errors, errors.New("expected object field name"))
	}
	if _, tkn, _ := p.sc.Scan(); tkn != scanner.COLON {
		p.errors = append(p.errors, errors.New("expected object field colon"))
	}
	value = p.parseValue()

	return
}
func (p *Parser) parseValue() (value ast.Value) {
	_, tkn, _ := p.sc.Peek()
	switch tkn {
	case scanner.DOLLAR:
		p.sc.Scan()
		_, tkn, lit := p.sc.Scan()
		if tkn != scanner.NAME {
			p.errors = append(p.errors, errors.New("expected a name to follow $"))
			break
		}
		value.Variable = &lit
	case scanner.INT:
		_, _, lit := p.sc.Scan()
		val, err := strconv.Atoi(lit)
		if err != nil {
			p.errors = append(p.errors, errors.New("not an integer"))
		}
		value.Int = &val
	case scanner.FLOAT:
		_, _, lit := p.sc.Scan()
		val, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			p.errors = append(p.errors, errors.New("not a float"))
		}
		value.Float = &val
	case scanner.STRING:
		_, _, lit := p.sc.Scan()
		val := lit[1 : len(lit)-1]
		value.String = &val
	case scanner.BLOCKSTRING:
		_, _, lit := p.sc.Scan()
		val := lit[3 : len(lit)-3]
		value.String = &val
	case scanner.BOOL:
		_, _, lit := p.sc.Scan()
		b, _ := strconv.ParseBool(lit)
		value.Bool = &b
	case scanner.NAME:
		_, _, lit := p.sc.Scan()
		if lit == "null" {
			value.IsNull = true
			break
		}
		value.Enum = &lit
	case scanner.LSQUARE:
		p.sc.Scan()
		value.List = make([]ast.Value, 0)
		for {
			_, tkn, _ := p.sc.Peek()
			if tkn == scanner.RSQUARE {
				break
			}
			if tkn == scanner.EOF {
				p.errors = append(p.errors, errors.New("unexpected EOF"))
			}

			value.List = append(value.List, p.parseValue())
		}
	case scanner.LCURLY:
		p.sc.Scan()
		value.Object = make(map[string]ast.Value)
		for {
			_, tkn, _ := p.sc.Peek()
			if tkn == scanner.RCURLY {
				break
			}
			if tkn == scanner.EOF {
				p.errors = append(p.errors, errors.New("unexpected EOF"))
			}

			name, value := p.parseObjectField()
			value.Object[name] = value
		}
	default:
		p.sc.Scan()
		p.errors = append(p.errors, errors.New("unexpected token"))
	}

	return
}
