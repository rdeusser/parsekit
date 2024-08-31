package parser

import (
	"errors"
	"fmt"

	"github.com/rdeusser/parsekit"
	"github.com/rdeusser/parsekit/ast"
	"github.com/rdeusser/parsekit/lexer"
	"github.com/rdeusser/parsekit/token"
)

// Parser is a generic parser implementation.
type Parser struct {
	l      *lexer.Lexer
	config Config
	pos    int
	tokens []token.Token
	logger parsekit.Logger
}

// Rule is a parser rule with a name, matcher, and an action to take if that matcher matches
// some input.
type Rule struct {
	Name   string
	Match  Matcher
	Action Action
}

// Config configures the parser to respond to the provided rules.
type Config struct {
	Rules []Rule
}

// Option sets options on parsers.
type Option func(*Parser)

// WithLogger sets a logger to use in debugging. Only debug logs are used.
func WithLogger(logger parsekit.Logger) Option {
	return func(p *Parser) {
		p.logger = logger
	}
}

// New constructs a new Parser.
func New(l *lexer.Lexer, config Config, options ...Option) *Parser {
	parser := &Parser{
		l:      l,
		config: config,
		tokens: make([]token.Token, 0),
		logger: parsekit.DefaultLogger,
	}

	for _, option := range options {
		option(parser)
	}

	return parser
}

func (p *Parser) Parse(input string) (file *ast.File, err error) {
	file = &ast.File{
		Nodes: make([]ast.Node, 0),
	}

	p.tokens, err = p.l.Lex(input)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	for p.pos < len(p.tokens) {
		curToken := p.tokens[p.pos]
		matched := false
		for _, rule := range p.config.Rules {
			p.logger.Debug("Attempting to match %q with token %q", rule.Name, curToken)

			if rule.Match(curToken) {
				p.logger.Debug("Running action %q", rule.Name)

				node, err := rule.Action(p, curToken)
				var perr Error
				if errors.As(err, &perr) {
					if perr.GotoNextRule {
						p.logger.Debug("Received an error from %q, moving to next rule", rule.Name)
						continue
					} else {
						_ = perr.Error()
						return nil, perr
					}
				} else if errors.Is(err, ErrGotoNextRule) {
					p.logger.Debug("Moving to next rule")
					continue
				} else if err != nil {
					return nil, err
				}

				file.Nodes = append(file.Nodes, node)
				matched = true
				break
			}
		}

		if !matched {
			// TODO(rdeusser): add output with line numbers and an up arrow at position.
			return nil, fmt.Errorf("parser error: no rule to handle token %q", input[curToken.Start.Pos:curToken.End.Pos])
		}

		p.Next()
	}

	return file, nil
}

func (p *Parser) Lookahead(n int) []token.Token {
	if p.pos+n >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos : p.pos+n]
}

func (p *Parser) Lookbehind(n int) []token.Token {
	if p.pos-n < 0 {
		return nil
	}
	return p.tokens[p.pos-n : p.pos]
}

func (p *Parser) Next() token.Token {
	p.pos++
	if p.pos >= len(p.tokens) {
		return token.NoToken
	}
	return p.tokens[p.pos]
}

func (p *Parser) Backup() token.Token {
	p.pos--
	if p.pos < 0 {
		return token.NoToken
	}
	return p.tokens[p.pos]
}
