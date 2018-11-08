package scanner_test

import (
	"testing"

	"github.com/dianelooney/graphql/scanner"
)

func expectScanResult(t *testing.T, s *scanner.Scanner, token scanner.Token, lit string) {
	_, tkn, l := s.Scan()
	if tkn != token {
		t.Errorf("Expected to scan a '%v' token, but got '%v'\n", token, tkn)
	}
	if l != lit {
		t.Errorf("Expected to scan a literal '%v', but got '%v'\n", lit, l)
	}
}
func TestScan(t *testing.T) {
	src := `
"something\""
"""some
1234
"bananas"
thing"""
"""something\""""""
12345
0
-0
-1234
1234.
-1234.
1234e1
-1234e1
!$()...:=@[]{|}
	`
	s := &scanner.Scanner{}
	s.Init([]byte(src))
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.STRING, `"something\""`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.BLOCKSTRING, `"""some
1234
"bananas"
thing"""`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.BLOCKSTRING, `"""something\""""""`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.INT, `12345`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.INT, `0`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.INT, `-0`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.INT, `-1234`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.FLOAT, `1234.`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.FLOAT, `-1234.`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.FLOAT, `1234e1`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.FLOAT, `-1234e1`)
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.BANG, `!`)
	expectScanResult(t, s, scanner.DOLLAR, `$`)
	expectScanResult(t, s, scanner.LPAREN, `(`)
	expectScanResult(t, s, scanner.RPAREN, `)`)
	expectScanResult(t, s, scanner.ELLIPSIS, `...`)
	expectScanResult(t, s, scanner.COLON, `:`)
	expectScanResult(t, s, scanner.EQL, `=`)
	expectScanResult(t, s, scanner.AT, `@`)
	expectScanResult(t, s, scanner.LSQUARE, `[`)
	expectScanResult(t, s, scanner.RSQUARE, `]`)
	expectScanResult(t, s, scanner.LCURLY, `{`)
	expectScanResult(t, s, scanner.BAR, `|`)
	expectScanResult(t, s, scanner.RCURLY, `}`)

	s = &scanner.Scanner{}
	s.Init([]byte(`
"\"something
`))
	expectScanResult(t, s, scanner.NEWLINE, "\n")
	expectScanResult(t, s, scanner.ILLEGAL, `"\"something`)
}
