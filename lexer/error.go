package lexer

import (
	"fmt"
)

type Error struct {
	Lexer        *Lexer
	Msg          string
	GotoNextRule bool
}

func (e Error) Error() string {
	if e.Lexer == nil {
		return "Lexer cannot be nil"
	}
	if e.Msg == "" {
		return "Msg cannot be empty"
	}
	return fmt.Sprintf("%s at %s", e.Msg, e.Lexer.curPos)
}
