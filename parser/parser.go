package parser

import (
	"github.com/dianelooney/graphql/parser/ast"
	"github.com/dianelooney/graphql/scanner"
)

type Parser struct {
	sc scanner.Scanner
}

func (p *Parser) Init(src []byte) {
	p.sc = scanner.Scanner{}
	p.sc.Init(src)
}
func (p *Parser) Parse() (d *ast.Document, err error) {

	return
}
func (p *Parser) scanDocument() (d ast.Document, err error) {
	for {
		defn, err := p.scanDefinition()
		if err != nil {
			return d, err
		}

	}
}
func (p *Parser) scanDefinition() (d *ast.Definition, err error) {

}

func (p *Parser) skipWhiteSpace() {
	for {
		_, t, _ := p.sc.Peek()
		if t == scanner.WHITESPACE || t == scanner.NEWLINE {
			p.sc.Scan()
		} else {
			return
		}
	}
}
func (p *Parser) peek() (position scanner.Position, token scanner.Token, literal string) {
	p.skipWhiteSpace()
	return p.sc.Peek()
}
func (p *Parser) scan() (position scanner.Position, token scanner.Token, literal string) {
	p.skipWhiteSpace()
	return p.sc.Scan()
}
