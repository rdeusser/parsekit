package parser

import (
	"github.com/rdeusser/parsekit/token"
)

type Matcher func(token.Token) bool

func IsIdentifier(tok token.Token) bool {
	return tok.Type == token.IDENT
}

func IsOperator(tok token.Token) bool {
	return tok.Type >= 2000 && tok.Type <= 2999
}

func IsKeyword(tok token.Token) bool {
	return tok.Type >= 3000 || tok.Type <= 999
}

func IsPackage(tok token.Token) bool {
	return tok.Type == token.PACKAGE
}

func IsStruct(tok token.Token) bool {
	return tok.Type == token.STRUCT
}

func IsString(tok token.Token) bool {
	return tok.Type == token.STRING
}
