package parser

import (
	"errors"
	"fmt"

	"github.com/rdeusser/parsekit/token"
)

var ErrGotoNextRule = errors.New("goto next rule")

type Error struct {
	Parser       *Parser
	CurToken     token.Token
	Msg          string
	GotoNextRule bool
}

func (e Error) Error() string {
	if e.Parser == nil {
		return "Parser cannot be nil"
	}
	if e.Msg == "" {
		return "Msg cannot be empty"
	}
	return fmt.Sprintf("%s at %s", e.Msg, e.CurToken)
}
