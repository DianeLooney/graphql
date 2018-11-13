package parser

import (
	"errors"

	"github.com/dianelooney/graphql/scanner"
)

func (p *Parser) hasNext(token scanner.Token, literal string) bool {
	_, tkn, lit := p.sc.Peek()
	return tkn == token && lit == literal
}
func (p *Parser) hasNextName(literal string) bool {
	return p.hasNext(scanner.NAME, literal)
}
func (p *Parser) hasNextTkn(token scanner.Token) bool {
	_, tkn, _ := p.sc.Peek()
	return tkn == token
}
func (p *Parser) consumeNameLiteral(literal string) {
	if !p.hasNext(scanner.NAME, literal) {
		p.errors = append(p.errors, errors.New("expected to find the name "+literal))
		return
	}

	p.sc.Scan()
}
func (p *Parser) consumeName() string {
	return p.consumeToken(scanner.NAME)
}
func (p *Parser) consumeToken(tkn scanner.Token) string {
	if !p.hasNextTkn(tkn) {
		p.errors = append(p.errors, errors.New("expected to find a different token"))
		return ""
	}

	_, _, lit := p.sc.Scan()
	return lit
}
func (p *Parser) consumeString() string {
	if !p.hasNextTkn(scanner.STRING) &&
		!p.hasNextTkn(scanner.BLOCKSTRING) {
		p.errors = append(p.errors, errors.New("expected to find a string"))
		return ""
	}

	if p.hasNextTkn(scanner.STRING) {
		_, _, lit := p.sc.Scan()
		val := lit[1 : len(lit)-1]
		return val
	}
	if p.hasNextTkn(scanner.BLOCKSTRING) {
		_, _, lit := p.sc.Scan()
		val := lit[3 : len(lit)-3]
		return val
	}

	panic("parser bug")
}
