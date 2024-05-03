package parser

import (
	"errors"
	"fmt"

	"github.com/rdeusser/parsekit"
	"github.com/rdeusser/parsekit/ast"
	"github.com/rdeusser/parsekit/lexer"
)

// Parser is a generic parser implementation.
type Parser struct {
	l      *lexer.Lexer
	config Config
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
		logger: parsekit.DefaultLogger,
	}

	for _, option := range options {
		option(parser)
	}

	return parser
}

func (p *Parser) Parse(input string) (*ast.File, error) {
	file := &ast.File{
		Nodes: make([]ast.Node, 0),
	}

	tokens, err := p.l.Lex(input)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	pos := 0

	for pos < len(tokens) {
		curToken := tokens[pos]
		matched := false
		for _, rule := range p.config.Rules {
			p.logger.Debug("Attempting to match %q with token %q", rule.Name, curToken)

			if rule.Match(tokens[pos]) {
				p.logger.Debug("Running action %q", rule.Name)

				node, consumed, err := rule.Action(p, tokens, pos)
				var perr Error
				if errors.As(err, &perr) {
					if perr.GotoNextRule {
						p.logger.Debug("Received an error from %q, moving to next rule", rule.Name)
						continue
					} else {
						_ = perr.Error()
						return nil, perr
					}
				} else if err != nil {
					return nil, err
				}

				pos += consumed
				file.Nodes = append(file.Nodes, node)
				matched = true
				break
			}
		}

		if !matched {
			// TODO(rdeusser): add output with line numbers and an up arrow at position.
			return nil, fmt.Errorf("parser error: no rule to handle token %s", curToken)
		}
	}

	return file, nil
}
