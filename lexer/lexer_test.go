package lexer

import (
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/parsekit/token"
)

func TestLexer(t *testing.T) {
	tests := map[string]struct {
		input   string
		config  Config
		wantErr assert.ErrorAssertionFunc
		want    autogold.Value
	}{
		"identifier": {
			"foo",
			Config{
				Rules: []Rule{
					{Name: "LexIdentifier", Match: IsLetter, Action: LexIdentifier},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(4),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    3,
						Line:   1,
						Column: 4,
					},
					Literal: "foo",
				},
			}),
		},
		"char": {
			"'f'",
			Config{
				Rules: []Rule{
					{Name: "LexChar", Match: IsSingleQuote, Action: LexChar},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(6),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    3,
						Line:   1,
						Column: 4,
					},
					Literal: "'f'",
				},
			}),
		},
		"newline char": {
			"'\n'",
			Config{
				Rules: []Rule{
					{Name: "LexChar", Match: IsSingleQuote, Action: LexChar},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(6),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    3,
						Line:   1,
						Column: 4,
					},
					Literal: "'\n'",
				},
			}),
		},
		"too many characters in char": {
			"'foo'",
			Config{
				Rules: []Rule{
					{Name: "LexChar", Match: IsSingleQuote, Action: LexChar},
				},
			},
			assert.Error,
			autogold.Expect([]token.Token{}),
		},
		"string": {
			"\"foo\"",
			Config{
				Rules: []Rule{
					{Name: "LexString", Match: IsDoubleQuote, Action: LexString},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(5),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    5,
						Line:   1,
						Column: 6,
					},
					Literal: `"foo"`,
				},
			}),
		},
		"string with newline": {
			"\"foo\n\"",
			Config{
				Rules: []Rule{
					{Name: "LexString", Match: IsDoubleQuote, Action: LexString},
				},
			},
			assert.Error,
			autogold.Expect([]token.Token{}),
		},
		"raw string": {
			"`hello world`",
			Config{
				Rules: []Rule{
					{Name: "LexRawString", Match: IsBackQuote, Action: LexRawString},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(5),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    13,
						Line:   1,
						Column: 14,
					},
					Literal: "`hello world`",
				},
			}),
		},
		"raw string with newlines": {
			"`hello\nworld`",
			Config{
				Rules: []Rule{
					{Name: "LexRawString", Match: IsBackQuote, Action: LexRawString},
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(5),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    13,
						Line:   2,
						Column: 8,
					},
					Literal: "`hello\nworld`",
				},
			}),
		},
		"operator": {
			"(",
			Config{
				Rules: []Rule{
					{Name: "LexOperator", Match: IsOperator, Action: LexOperator},
				},
				Operators: map[string]token.TokenType{
					"(":   2001,
					"+":   2002,
					"<<=": 2003,
				},
			},
			assert.NoError,
			autogold.Expect([]token.Token{
				{
					Type: token.TokenType(2001),
					Start: token.Position{
						Line:   1,
						Column: 1,
					},
					End: token.Position{
						Pos:    1,
						Line:   1,
						Column: 2,
					},
					Literal: "(",
				},
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			lexer := New(tt.config)
			tokens, err := lexer.Lex(tt.input)
			tt.wantErr(t, err)
			tt.want.Equal(t, tokens)
		})
	}
}
