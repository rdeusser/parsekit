package ast

import (
	"github.com/rdeusser/parsekit/token"
)

// Node represents a node in an AST.
type Node interface {
	Start() token.Position
	End() token.Position
}

type Statement interface {
	Node
	StatementNode()
}

type Declaration interface {
	Node
	DeclarationNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type File struct {
	Nodes []Node
}

func (x *File) Start() token.Position { return x.Nodes[0].Start() }
func (x *File) End() token.Position   { return x.Nodes[len(x.Nodes)-1].End() }

type DeclarationStatement struct {
	Declaration Declaration
}

func (x *DeclarationStatement) Start() token.Position { return x.Declaration.Start() }
func (x *DeclarationStatement) End() token.Position   { return x.Declaration.End() }

type ExpressionStatement struct {
	Expression Expression
}

func (x *ExpressionStatement) Start() token.Position { return x.Expression.Start() }
func (x *ExpressionStatement) End() token.Position   { return x.Expression.End() }

type Identifier struct {
	Name    string
	NamePos token.Position
}

func (x *Identifier) Start() token.Position { return x.NamePos }
func (x *Identifier) End() token.Position {
	return token.Position{
		Pos:    x.NamePos.Pos + len(x.Name),
		Line:   x.NamePos.Line,
		Column: x.NamePos.Column + len(x.Name),
	}
}

func (x *DeclarationStatement) StatementNode() {}
func (x *ExpressionStatement) StatementNode()  {}

func (x *Identifier) ExpressionNode() {}
