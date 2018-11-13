package parser

import "github.com/dianelooney/graphql/scanner"

func (p *Parser) hasNext(token scanner.Token, literal string) bool {
	_, tkn, lit := p.sc.Peek()
	return tkn == token && lit == literal
}
