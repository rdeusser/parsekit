package token

import (
	"fmt"
)

// TokenType is a type of token. Duh?
type TokenType int

const (
	ILLEGAL    TokenType = iota // ILLEGAL
	EOF                         // EOF
	COMMENT                     // COMMENT
	WHITESPACE                  // WHITESPACE

	// Identifiers and basic type literals.
	IDENT  // IDENT
	STRING // STRING
	CHAR   // CHAR
	NUMBER // NUMBER
	FLOAT  // FLOAT

	// 0-999 is reserved for parsekit.

	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// User-defined identifiers and basic type literals.
	LiteralStart = 1000

	// User-defined operators.
	OperatorStart = 2000

	// User-defined keywords.
	KeywordStart = 3000
)

type Token struct {
	Type    TokenType
	Start   Position
	End     Position
	Literal string
}

var NoToken = Token{}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %d, Start: %s, End: %s}", t.Type, t.Start, t.End)
}

func (t Token) Pos() string {
	if t.Start.Line == t.End.Line {
		return fmt.Sprintf("%d:%d:%d", t.Start.Line, t.Start.Column, t.End.Column)
	}
	return fmt.Sprintf("%s:%s", t.Start, t.End)
}
