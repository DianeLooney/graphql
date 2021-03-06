package scanner

import (
	"bytes"
	"regexp"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT
	WHITESPACE
	NEWLINE

	NAME

	BOOL
	INT
	FLOAT
	STRING
	BLOCKSTRING

	COMMA
	BANG
	DOLLAR
	LPAREN
	RPAREN
	ELLIPSIS
	COLON
	EQL
	AT
	LSQUARE
	RSQUARE
	LCURLY
	RCURLY
	BAR
	AMP
)

type Position struct {
	Line   int
	Offset int
}

type Scanner struct {
	src    []byte
	offset int
	data   []result
	idx    int
}
type scanFunc func(s *Scanner) (token Token, lit string)

var skipFuncs = [...]scanFunc{
	(*Scanner).scanWhitespace,
	(*Scanner).scanComment,
	(*Scanner).scanNewline,
	(*Scanner).scanComma,
}
var scanFuncs = [...]scanFunc{
	(*Scanner).scanEOF,
	(*Scanner).scanPunctuator,
	(*Scanner).scanFloat,
	(*Scanner).scanInt,
	(*Scanner).scanBoolean,
	(*Scanner).scanBlockString,
	(*Scanner).scanString,
	(*Scanner).scanName,
}

type result struct {
	pos Position
	tkn Token
	lit string
}

func (s *Scanner) Init(src []byte) {
	s.src = src
	s.offset = 0
	s.idx = 0
	s.data = make([]result, 0)
	for {
		pos, tkn, lit := s.scan()
		if tkn == EOF {
			break
		}

		s.data = append(s.data, result{pos, tkn, lit})
	}
}

func (s *Scanner) get(i int) (pos Position, token Token, lit string) {
	if i >= len(s.data) {
		token = EOF
	} else {
		res := s.data[i]
		pos, token, lit = res.pos, res.tkn, res.lit
	}

	return
}

func (s *Scanner) Peek() (pos Position, token Token, lit string) {
	return s.PeekN(0)
}
func (s *Scanner) PeekN(n int) (pos Position, token Token, lit string) {
	return s.get(s.idx + n)
}

func (s *Scanner) scan() (pos Position, token Token, lit string) {
	s.skipWhitespace()
	pos = Position{0, s.offset}
	for _, f := range scanFuncs {
		token, lit = f(s)

		if token == ILLEGAL && len(lit) == 0 {
			continue
		}

		s.consume(lit)
		return
	}

	lit = string(s.src[0])
	s.consume(lit)
	return
}

func (s *Scanner) Scan() (pos Position, token Token, lit string) {
	pos, token, lit = s.Peek()
	s.idx++

	return
}

func (s *Scanner) skipWhitespace() {
	done := false
	for !done {
		done = true
		for _, f := range skipFuncs {
			token, lit := f(s)
			if token == ILLEGAL && len(lit) == 0 {
				continue
			}

			s.consume(lit)
			done = false
		}
	}
}

func (s *Scanner) consume(lit string) {
	s.src = s.src[len(lit):]
	s.offset += len(lit)
}

func (s *Scanner) scanRegex(r *regexp.Regexp, t Token) (token Token, lit string) {
	match := r.Find(s.src)
	if match == nil {
		return ILLEGAL, ""
	}

	return t, string(match)
}

func (s *Scanner) scanEOF() (token Token, lit string) {
	if len(s.src) == 0 {
		return EOF, ""
	}
	return ILLEGAL, ""
}

var regexComment = regexp.MustCompile(`^(#[^\n]*)`)

func (s *Scanner) scanComment() (token Token, lit string) {
	return s.scanRegex(regexComment, COMMENT)
}

var regexWhitespace = regexp.MustCompile(`^([ \t]+)`)

func (s *Scanner) scanWhitespace() (token Token, lit string) {
	return s.scanRegex(regexWhitespace, WHITESPACE)
}

var regexName = regexp.MustCompile(`^([_A-Za-z][_0-9A-Za-z]*)`)

func (s *Scanner) scanName() (token Token, lit string) {
	return s.scanRegex(regexName, NAME)
}

var regexNewline = regexp.MustCompile(`^(\r|\n|\r\n)`)

func (s *Scanner) scanNewline() (token Token, lit string) {
	return s.scanRegex(regexNewline, NEWLINE)
}

var regexComma = regexp.MustCompile(`^(,)`)

func (s *Scanner) scanComma() (token Token, lit string) {
	return s.scanRegex(regexComma, COMMA)
}

var regexPunctuator = regexp.MustCompile(`^(\!|\$|\(|\)|\.\.\.|\:|\=|\@|\[|\]|\{|\||\}|&)`)
var punctuatorMap = map[string]Token{
	"!":   BANG,
	"$":   DOLLAR,
	"(":   LPAREN,
	")":   RPAREN,
	"...": ELLIPSIS,
	":":   COLON,
	"=":   EQL,
	"@":   AT,
	"[":   LSQUARE,
	"]":   RSQUARE,
	"{":   LCURLY,
	"|":   BAR,
	"}":   RCURLY,
	"&":   AMP,
}

func (s *Scanner) scanPunctuator() (token Token, lit string) {
	token, lit = s.scanRegex(regexPunctuator, Token(-1))
	if token != Token(-1) {
		return ILLEGAL, ""
	}
	token, ok := punctuatorMap[lit]
	if !ok {
		panic("Unhandled literal matched by scanPunctuator() (it should be added to punctuatorMap): '" + lit + "'")
	}

	return
}

var regexInt = regexp.MustCompile(`^(-?(0|[1-9][0-9]*))`)

func (s *Scanner) scanInt() (token Token, lit string) {
	return s.scanRegex(regexInt, INT)
}

var regexFloat = regexp.MustCompile(`^((-?(0|[1-9][0-9]*))(\.[0-9]*|[eE][+-]?[0-9]*|\.[0-9]*[eE][+-]?[0-9]*))`)

func (s *Scanner) scanFloat() (token Token, lit string) {
	return s.scanRegex(regexFloat, FLOAT)
}

var regexBoolean = regexp.MustCompile(`^(true|false)`)

func (s *Scanner) scanBoolean() (token Token, lit string) {
	return s.scanRegex(regexBoolean, BOOL)
}

var escapeSequences = map[byte]struct{}{
	'"':  struct{}{},
	'\\': struct{}{},
	'/':  struct{}{},
	'b':  struct{}{},
	'f':  struct{}{},
	'n':  struct{}{},
	'r':  struct{}{},
	't':  struct{}{},
}

func (s *Scanner) scanString() (token Token, lit string) {
	if s.src[0] != '"' {
		return ILLEGAL, ""
	}
	for i := 1; i < len(s.src); i++ {
		if s.src[i] == '"' {
			return STRING, string(s.src[:i+1])
		}

		if s.src[i] == '\n' {
			return ILLEGAL, string(s.src[:i])
		}

		if s.src[i] == '\r' && i+1 < len(s.src) && s.src[i+1] == '\n' {
			return ILLEGAL, string(s.src[:i])
		}

		if s.src[i] == '\\' && i+1 < len(s.src) {
			i++
			_, ok := escapeSequences[s.src[i]]
			if ok {
				continue
			}
			if s.src[i] != 'u' {
				return ILLEGAL, string(s.src[:i+2])
			}
			if i+4 < len(s.src) {
				return ILLEGAL, string(s.src[:i+5])
			}

			testChar := func(b byte) bool {
				if b >= 'a' && b <= 'f' {
					return true
				}
				if b >= 'A' && b <= 'F' {
					return true
				}
				if b >= '0' && b <= '9' {
					return true
				}
				return false
			}

			if testChar(s.src[i+1]) &&
				testChar(s.src[i+2]) &&
				testChar(s.src[i+3]) &&
				testChar(s.src[i+4]) {
				i += 4
				continue
			}
			return ILLEGAL, string(s.src[:i+1])
		}
	}
	return ILLEGAL, ""
}

var tripleQuote = []byte(`"""`)
var escapedTripleQuote = []byte(`\"""`)

func (s *Scanner) scanBlockString() (token Token, lit string) {
	if !bytes.HasPrefix(s.src, tripleQuote) {
		return ILLEGAL, ""
	}

	for i := 3; i < len(s.src); i++ {
		if bytes.HasPrefix(s.src[i:], tripleQuote) {
			return BLOCKSTRING, string(s.src[0 : i+3])
		}

		if bytes.HasPrefix(s.src[i:], escapedTripleQuote) {
			i += 3
		}
	}
	return ILLEGAL, ""
}
