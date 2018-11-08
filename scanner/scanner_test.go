package scanner_test

import (
	"testing"

	"github.com/dianelooney/graphql/scanner"
)

func TestScan(t *testing.T) {
	src := `
	"somethun\""
	"""some




	thing"""
	"""something\""""""
	`
	s := scanner.Scanner{}
	s.Init([]byte(src))
	for {
		pos, tkn, lit := s.Scan()

		if tkn == 0 {
			t.Logf("Scan[ILLEGAL]: '%v' %v\n", lit, pos.Offset)
		} else if tkn == scanner.NEWLINE {
			t.Logf("Scan[%v]: newline at %v\n", tkn, pos.Offset)
		} else {
			t.Logf("Scan[%v]: '%v' at %v\n", tkn, lit, pos.Offset)
		}

		if tkn == scanner.EOF {
			break
		}
	}
}
