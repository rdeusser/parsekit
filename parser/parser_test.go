package parser

import (
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/parsekit/ast"
	"github.com/rdeusser/parsekit/lexer"
	"github.com/rdeusser/parsekit/token"
)

func TestParser(t *testing.T) {
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
					{Name: "ParseIdentifier", Match: IsIdentifier, Action: ParseIdentifier},
				},
			},
			assert.NoError,
			autogold.Expect(&ast.File{
				Nodes: []ast.Node{
					&ast.Identifier{
						Name: "foo",
						NamePos: token.Position{
							Line:   1,
							Column: 1,
						},
					},
				},
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			l := lexer.New(lexer.DefaultConfig)
			p := New(l, tt.config)
			got, err := p.Parse(tt.input)
			tt.wantErr(t, err)
			tt.want.Equal(t, got)
		})
	}
}
