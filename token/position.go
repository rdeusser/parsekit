package token

import (
	"fmt"
	"strconv"
)

// Position is the position of a token.
type Position struct {
	Pos    int // character position
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)
}

// IsValid reports whether the position is valid.
func (p Position) IsValid() bool { return p.Pos >= 0 && p.Line > 0 && p.Column > 0 }

// String returns the string form of a position.
func (p Position) String() string {
	s := ""
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += strconv.Itoa(p.Line)
		if p.Column != 0 {
			s += fmt.Sprintf(":%d", p.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
