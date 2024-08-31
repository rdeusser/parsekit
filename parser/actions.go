package parser

import (
	"fmt"
	"unicode"

	"github.com/rdeusser/parsekit/ast"
	"github.com/rdeusser/parsekit/token"
)

// Action is a type for parsing tokens into an AST.
type Action func(p *Parser, tok token.Token) (ast.Node, error)

func ParseIdentifier(p *Parser, tok token.Token) (ast.Node, error) {
	ident := &ast.Identifier{
		Name: tok.Literal,
		Pos:  tok.Start,
	}

	return ident, nil
}

func ParsePackage(p *Parser, tok token.Token) (ast.Node, error) {
	pkg := &ast.Package{
		Token: tok.Start,
	}

	name, err := ParseIdentifier(p, p.Next())
	if err != nil {
		return nil, err
	}

	pkg.Name = name.(*ast.Identifier)

	return pkg, nil
}

func ParseStruct(p *Parser, tok token.Token) (ast.Node, error) {
	if tok.Type != token.STRUCT {
		return nil, fmt.Errorf("expected struct, got %s", tok)
	}

	node := &ast.Struct{
		Token: tok.Start,
	}

	tok = p.Next()

	if unicode.IsUpper(rune(tok.Literal[0])) {
		node.Public = true
	}

	name, err := ParseIdentifier(p, tok)
	if err != nil {
		return nil, err
	}

	node.Name = name.(*ast.Identifier)

	tok = p.Next() // move past struct name
	tok = p.Next() // move past {

	return node, nil
}
