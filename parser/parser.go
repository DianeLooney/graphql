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
	doc.Types = make(map[string]ast.TypeDef)
	doc.Directives = make(map[string]ast.DirectiveDef)
	doc.Fragments = make(map[string]ast.FragmentDef)

	for {
		var desc *string
		if p.hasNextTkn(scanner.EOF) {
			break
		}

		if p.hasNextName("schema") {
			schema := p.parseSchema()
			doc.Schema = &schema
			continue
		} else if p.hasNextName("query") ||
			p.hasNextName("mutation") ||
			p.hasNextName("subscription") {
			op := p.parseOperationDef()
			doc.Operations[*op.Name] = op
			continue
		} else if p.hasNextTkn(scanner.LCURLY) {
			op := p.parseOperationDef()
			doc.Operations[*op.Name] = op
			continue
		} else if p.hasNextName("fragment") {
			frag := p.parseFragmentDef()
			doc.Fragments[frag.Name] = frag
			continue
		}

		desc = p.parseDescription()
		if p.hasNextName("scalar") {
			scalar := p.parseScalarTypeDefinition(desc)
			doc.Types[scalar.Name] = ast.TypeDef{ScalarDef: &scalar}
		} else if p.hasNextName("type") {
			obj := p.parseObjectTypeDefinition(desc)
			doc.Types[obj.Name] = ast.TypeDef{ObjectTypeDef: &obj}
		} else if p.hasNextName("interface") {
			intf := p.parseInterfaceTypeDef(desc)
			doc.Types[intf.Name] = ast.TypeDef{InterfaceDef: &intf}
		} else if p.hasNextName("union") {
			union := p.parseUnionDef(desc)
			doc.Types[union.Name] = ast.TypeDef{UnionDef: &union}
		} else if p.hasNextName("enum") {
			enum := p.parseEnumDef(desc)
			doc.Types[enum.Name] = ast.TypeDef{EnumDef: &enum}
		} else if p.hasNextName("input") {
			input := p.parseInputDef(desc)
			doc.Types[input.Name] = ast.TypeDef{InputDef: &input}
		} else if p.hasNextName("directive") {
			dir := p.parseDirectiveDef(desc)
			doc.Directives[dir.Name] = dir
		} else {
			_, _, lit := p.sc.Scan()
			p.errors = append(p.errors, errors.New("unknown: "+lit))
		}
	}

	return
}

func (p *Parser) parseOperationDef() (op ast.Operation) {
	if p.hasNextName("query") ||
		p.hasNextName("mutation") ||
		p.hasNextName("subscription") {
		op.OpType = p.consumeName()
		if p.hasNextTkn(scanner.NAME) {
			name := p.consumeName()
			op.Name = &name
		}
		if p.hasNextTkn(scanner.LPAREN) {
			p.consumeToken(scanner.LPAREN)
			for {
				if p.hasNextTkn(scanner.RPAREN) || p.hasNextTkn(scanner.EOF) {
					break
				}
				op.Variables = append(op.Variables, p.parseVariableDef())
			}
			p.consumeToken(scanner.RPAREN)
		}
		op.Directives = p.parseDirectives()
	} else {
		op.OpType = "query"
	}
	op.SelectionSet = p.parseSelectionSet()

	return
}

func (p *Parser) parseSelection() (sel ast.Selection) {
	if !p.hasNextTkn(scanner.ELLIPSIS) {
		field := p.parseField()
		sel.Field = &field
	} else {
		p.consumeToken(scanner.ELLIPSIS)
		if p.hasNextName("on") {
			frag := p.parseInlineFragment()
			sel.InlineFragment = &frag
		} else {
			frag := p.parseFragmentSpread()
			sel.FragmentSpread = &frag
		}
	}

	return
}
func (p *Parser) parseFragmentDef() (frag ast.FragmentDef) {
	p.consumeNameLiteral("fragment")
	frag.Name = p.consumeName()
	p.consumeNameLiteral("on")
	frag.Type = p.consumeName()
	frag.Directives = p.parseDirectives()
	frag.SelectionSet = p.parseSelectionSet()

	return
}
func (p *Parser) parseField() (field ast.Field) {
	n1 := p.consumeName()
	if p.hasNextTkn(scanner.COLON) {
		p.consumeToken(scanner.COLON)
		field.Alias = &n1
		field.Name = p.consumeName()
	} else {
		field.Name = n1
	}
	if p.hasNextTkn(scanner.LPAREN) {
		field.Arguments = p.parseArguments()
	}

	field.Directives = p.parseDirectives()
	if p.hasNextTkn(scanner.LCURLY) {
		field.SelectionSet = p.parseSelectionSet()
	}

	return
}
func (p *Parser) parseInlineFragment() (frag ast.InlineFragment) {
	// ... is already consumed, as we have to distinguish
	// between InlineFragment and FragmentSpread
	p.consumeNameLiteral("on")
	inline := ast.InlineFragment{}
	if p.hasNextTkn(scanner.NAME) {
		n := p.consumeName()
		inline.Type = &n
	}
	inline.Directives = p.parseDirectives()
	inline.SelectionSet = p.parseSelectionSet()

	return
}
func (p *Parser) parseFragmentSpread() (frag ast.FragmentSpread) {
	frag.Name = p.consumeName()
	frag.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseSelectionSet() (selections []ast.Selection) {
	p.consumeToken(scanner.LCURLY)
	for {
		if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
			break
		}

		selections = append(selections, p.parseSelection())
	}
	p.consumeToken(scanner.RCURLY)

	return
}

func (p *Parser) parseVariableDef() (vari ast.VariableDef) {
	p.consumeToken(scanner.DOLLAR)
	vari.Name = p.consumeName()
	p.consumeToken(scanner.COLON)
	if p.hasNextTkn(scanner.EQL) {
		p.consumeToken(scanner.EQL)
		v := p.parseValue()
		vari.DefaultValue = &v
	}
	vari.Directives = p.parseDirectives()

	return
}

func (p *Parser) parseSchema() (schema ast.Schema) {
	p.consumeNameLiteral("schema")
	schema.Directives = p.parseDirectives()
	p.consumeToken(scanner.LCURLY)
	for {
		if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
			break
		}

		schema.RootOperationTypeDefs = append(schema.RootOperationTypeDefs, p.parseRootOpTypeDefinition())
	}
	p.consumeToken(scanner.RCURLY)

	return
}
func (p *Parser) parseRootOpTypeDefinition() (def ast.RootOperationTypeDef) {
	def.OpType = p.consumeName()
	p.consumeToken(scanner.COLON)
	def.NamedType = p.consumeName()

	return
}
func (p *Parser) parseScalarTypeDefinition(desc *string) (scalar ast.ScalarDef) {
	scalar.Description = desc
	p.consumeNameLiteral("scalar")
	scalar.Name = p.consumeName()
	scalar.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseObjectTypeDefinition(desc *string) (obj ast.ObjectTypeDef) {
	obj.Description = desc
	p.consumeNameLiteral("type")
	obj.Name = p.consumeName()
	if p.hasNextName("implements") {
		obj.ImplementsInterface = p.parseImplements()
	}
	obj.Directives = p.parseDirectives()

	if !p.hasNextTkn(scanner.LCURLY) {
		return
	}
	obj.Fields = p.parseFieldDefs()

	return
}
func (p *Parser) parseInterfaceTypeDef(desc *string) (intf ast.InterfaceDef) {
	intf.Description = desc
	p.consumeNameLiteral("interface")
	intf.Name = p.consumeName()
	intf.Directives = p.parseDirectives()

	if !p.hasNextTkn(scanner.LCURLY) {
		return
	}
	intf.Fields = p.parseFieldDefs()

	return
}
func (p *Parser) parseEnumDef(desc *string) (enum ast.EnumDef) {
	enum.Description = desc
	p.consumeNameLiteral("enum")
	enum.Name = p.consumeName()
	enum.Directives = p.parseDirectives()
	if !p.hasNextTkn(scanner.LCURLY) {
		return
	}
	p.consumeToken(scanner.LCURLY)
	for {
		if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
			break
		}
		enum.Values = append(enum.Values, p.parseEnumValueDef())
	}
	p.consumeToken(scanner.RCURLY)

	return
}
func (p *Parser) parseEnumValueDef() (val ast.EnumValueDef) {
	val.Description = p.parseDescription()
	val.Name = p.consumeName()
	if val.Name == "true" || val.Name == "false" || val.Name == "" {
		p.errors = append(p.errors, errors.New("invalid enum value"))
		val.Name = ""
	}
	val.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseFieldDefs() (fields []ast.FieldDef) {
	p.consumeToken(scanner.LCURLY)
	for {
		if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
			break
		}

		fields = append(fields, p.parseFieldDef())
	}
	p.consumeToken(scanner.RCURLY)

	return
}
func (p *Parser) parseInputDef(desc *string) (input ast.InputDef) {
	input.Description = desc
	p.consumeNameLiteral("input")
	input.Name = p.consumeName()
	input.Directives = p.parseDirectives()
	if !p.hasNextTkn(scanner.LCURLY) {
		return
	}
	p.consumeToken(scanner.LCURLY)

	for {
		if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
			break
		}
		input.Fields = append(input.Fields, p.parseInputValueDef())
	}
	p.consumeToken(scanner.RCURLY)

	return
}
func (p *Parser) parseInputValueDef() (val ast.InputValueDef) {
	val.Description = p.parseDescription()
	val.Name = p.consumeName()
	p.consumeToken(scanner.COLON)
	val.Type = p.parseType()
	if p.hasNextTkn(scanner.EQL) {
		p.consumeToken(scanner.EQL)
		v := p.parseValue()
		val.DefaultValue = &v
	}
	val.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseFieldDef() (field ast.FieldDef) {
	field.Description = p.parseDescription()
	field.Name = p.consumeName()
	field.Arguments = p.parseArgumentsDefn()
	p.consumeToken(scanner.COLON)
	field.Type = p.parseType()
	field.Directives = p.parseDirectives()

	return
}
func (p *Parser) parseDirectiveDef(desc *string) (dir ast.DirectiveDef) {
	dir.Description = desc
	p.consumeNameLiteral("directive")
	p.consumeToken(scanner.AT)
	dir.Name = p.consumeName()
	dir.Arguments = p.parseArgumentsDefn()
	p.consumeNameLiteral("on")
	if p.hasNextTkn(scanner.BAR) {
		p.consumeToken(scanner.BAR)
	}
	for {
		location := p.consumeName()
		_, isExec := ast.ExecutableDirectiveLocations[location]
		_, isType := ast.TypeSystemDirectiveLocations[location]
		if !isExec && !isType {
			p.errors = append(p.errors, errors.New("unepected directive location"))
		} else {
			dir.Locations = append(dir.Locations, location)
		}
		if !p.hasNextTkn(scanner.BAR) {
			break
		}
		p.consumeToken(scanner.BAR)
	}

	return
}
func (p *Parser) parseType() (t ast.Type) {
	if p.hasNextTkn(scanner.LSQUARE) {
		p.consumeToken(scanner.LSQUARE)
		in := p.parseType()
		t.ListType = &in
		p.consumeToken(scanner.RSQUARE)
	} else if p.hasNextTkn(scanner.NAME) {
		n := p.consumeName()
		t.Name = &n
	}

	if p.hasNextTkn(scanner.BANG) {
		var nullType ast.Type
		nullType = t
		t = ast.Type{NonNullType: &nullType}
	}

	return
}
func (p *Parser) parseImplements() (implements []string) {
	p.consumeNameLiteral("implements")
	if p.hasNextTkn(scanner.AMP) {
		p.consumeToken(scanner.AMP)
	}
	for {
		n := p.consumeName()
		implements = append(implements, n)
		if !p.hasNextTkn(scanner.AMP) {
			break
		}
		p.consumeToken(scanner.AMP)
	}
	return
}
func (p *Parser) parseUnionDef(desc *string) (union ast.UnionDef) {
	union.Description = desc
	p.consumeNameLiteral("union")
	union.Name = p.consumeName()
	union.Directives = p.parseDirectives()
	if !p.hasNextTkn(scanner.EQL) {
		return
	}
	p.consumeToken(scanner.EQL)
	if p.hasNextTkn(scanner.BAR) {
		p.consumeToken(scanner.BAR)
	}
	for {
		n := p.consumeName()
		union.Types = append(union.Types, n)
		if !p.hasNextTkn(scanner.BAR) {
			break
		}
		p.consumeToken(scanner.BAR)
	}

	return
}
func (p *Parser) parseArgumentsDefn() (args []ast.InputValueDef) {
	if !p.hasNextTkn(scanner.LPAREN) {
		return
	}
	p.consumeToken(scanner.LPAREN)

	for {
		if p.hasNextTkn(scanner.RPAREN) || p.hasNextTkn(scanner.EOF) {
			break
		}
		args = append(args, p.parseInputValueDef())
	}
	p.consumeToken(scanner.RPAREN)

	return
}
func (p *Parser) parseDescription() (description *string) {
	if p.hasNextTkn(scanner.STRING) || p.hasNextTkn(scanner.BLOCKSTRING) {
		desc := p.consumeString()
		return &desc
	}

	return nil
}
func (p *Parser) parseDirectives() (directives []ast.Directive) {
	for {
		if !p.hasNextTkn(scanner.AT) {
			break
		}

		directives = append(directives, p.parseDirective())
	}

	return
}
func (p *Parser) parseDirective() (directive ast.Directive) {
	p.consumeToken(scanner.AT)
	directive.Name = p.consumeName()

	if p.hasNextTkn(scanner.LPAREN) {
		directive.Arguments = p.parseArguments()
	}

	return
}
func (p *Parser) parseArguments() (arguments map[string]ast.Value) {
	arguments = make(map[string]ast.Value)

	if !p.hasNextTkn(scanner.LPAREN) {
		p.errors = append(p.errors, errors.New("expected left paren to start argument list"))
		return
	}
	p.sc.Scan()

	for {
		if p.hasNextTkn(scanner.RPAREN) || p.hasNextTkn(scanner.EOF) {
			break
		}

		name, value := p.parseArgument()
		arguments[name] = value
	}
	p.consumeToken(scanner.RPAREN)
	return
}
func (p *Parser) parseArgument() (name string, value ast.Value) {
	name = p.consumeName()
	p.consumeToken(scanner.COLON)
	value = p.parseValue()

	return
}
func (p *Parser) parseObjectField() (name string, value ast.Value) {
	name = p.consumeName()
	p.consumeToken(scanner.COLON)
	value = p.parseValue()

	return
}
func (p *Parser) parseValue() (value ast.Value) {
	_, tkn, _ := p.sc.Peek()
	switch tkn {
	case scanner.DOLLAR:
		p.sc.Scan()
		name := p.consumeName()
		value.Variable = &name
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
		fallthrough
	case scanner.BLOCKSTRING:
		str := p.consumeString()
		value.String = &str
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
			if p.hasNextTkn(scanner.RSQUARE) || p.hasNextTkn(scanner.EOF) {
				break
			}

			value.List = append(value.List, p.parseValue())
		}
		p.consumeToken(scanner.RSQUARE)
	case scanner.LCURLY:
		p.sc.Scan()
		value.Object = make(map[string]ast.Value)
		for {
			if p.hasNextTkn(scanner.RCURLY) || p.hasNextTkn(scanner.EOF) {
				break
			}

			name, value := p.parseObjectField()
			value.Object[name] = value
		}
		p.consumeToken(scanner.RCURLY)
	default:
		p.sc.Scan()
		p.errors = append(p.errors, errors.New("unexpected token"))
	}

	return
}
