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
	Pos  token.Position
	Name string
}

func (x *Identifier) Start() token.Position { return x.Pos }
func (x *Identifier) End() token.Position {
	return token.Position{
		Pos:    x.Pos.Pos + len(x.Name),
		Line:   x.Pos.Line,
		Column: x.Pos.Column + len(x.Name),
	}
}

type Package struct {
	Token token.Position // position of 'package'
	Name  *Identifier
}

func (x *Package) Start() token.Position { return x.Token }
func (x *Package) End() token.Position   { return x.Name.End() }

type Struct struct {
	Token          token.Position // position of 'pub' or 'struct'
	Public         bool
	Name           *Identifier
	TypeParameters *TypeParameters
	Body           *Block
}

func (x *Struct) Start() token.Position { return x.Token }
func (x *Struct) End() token.Position   { return x.Body.End() }

type TypeParameters struct {
	Lbrack token.Position // position of '['
	List   []*TypeParameter
	Rbrack token.Position // position of ']'
}

func (x *TypeParameters) Start() token.Position { return x.Lbrack }
func (x *TypeParameters) End() token.Position   { return x.Rbrack }

type TypeParameter struct {
	Name *Identifier
	Type *Identifier
}

func (x *TypeParameter) Start() token.Position { return x.Name.Start() }
func (x *TypeParameter) End() token.Position   { return x.Type.End() }

type Block struct {
	Lbrace     token.Position
	Statements []Statement
	Rbrace     token.Position
}

func (x *Block) Start() token.Position { return x.Lbrace }
func (x *Block) End() token.Position   { return x.Rbrace }

func (x *DeclarationStatement) StatementNode() {}
func (x *ExpressionStatement) StatementNode()  {}

func (x *Package) DeclarationNode() {}
func (x *Struct) DeclarationNode()  {}

func (x *Identifier) ExpressionNode() {}
func (x *Block) ExpressionNode()      {}
