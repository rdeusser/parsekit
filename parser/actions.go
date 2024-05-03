package parser

import (
	"github.com/rdeusser/parsekit/ast"
	"github.com/rdeusser/parsekit/token"
)

// Action is a type for parsing tokens into an AST.
type Action func(p *Parser, tokens []token.Token, pos int) (node ast.Node, consumed int, err error)

func ParseIdentifier(p *Parser, tokens []token.Token, pos int) (ast.Node, int, error) {
	ident := &ast.Identifier{
		Name:    tokens[pos].Literal,
		NamePos: tokens[pos].Start,
	}

	return ident, 1, nil
}
