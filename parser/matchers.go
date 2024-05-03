package parser

import (
	"github.com/rdeusser/parsekit/token"
)

type Matcher func(token.Token) bool

func IsIdentifier(tok token.Token) bool {
	return tok.Type == token.IDENT
}

func IsKeyword(tok token.Token) bool {
	return tok.Type >= 3000
}

func IsString(tok token.Token) bool {
	return tok.Type == token.STRING
}
