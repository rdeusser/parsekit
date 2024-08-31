package lexer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rdeusser/parsekit"
	"github.com/rdeusser/parsekit/internal/loopdetector"
	"github.com/rdeusser/parsekit/token"
)

var eof = rune(0)

// Lexer is a generic lexer implementation.
type Lexer struct {
	input        string
	curPos       token.Position
	prevPos      token.Position
	config       Config
	logger       parsekit.Logger
	loopDetector *loopdetector.Detector
	mu           sync.Mutex
}

// Rule is a lexer rule with a name, matcher, and an action to take if that matcher matches
// some input.
type Rule struct {
	Name   string
	Match  Matcher
	Action Action
}

// Config configures the lexer to respond to the provided rules and user-defined operators and keywords.
type Config struct {
	SkipWhitespace bool
	Rules          []Rule
	Operators      map[string]token.TokenType
	Keywords       map[string]token.TokenType
}

// Option sets options on lexers.
type Option func(*Lexer)

// WithLogger sets a logger to use in debugging. Only debug logs are used.
func WithLogger(logger parsekit.Logger) Option {
	return func(l *Lexer) {
		l.logger = logger
	}
}

// New creates a new Lexer from a lexer config and options.
func New(config Config, options ...Option) *Lexer {
	lexer := &Lexer{
		curPos: token.Position{
			Line:   1,
			Column: 1,
		},
		config:       config,
		logger:       parsekit.DefaultLogger,
		loopDetector: loopdetector.New(),
		mu:           sync.Mutex{},
	}

	for _, option := range options {
		option(lexer)
	}

	return lexer
}

// Lex lexes the input and returns a slice of tokens, or an error.
func (l *Lexer) Lex(input string) ([]token.Token, error) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			<-ticker.C
			l.loopDetector.Detect(&l.mu, l.curPos.Pos)
			if l.loopDetector.IsLooping() {
				panic(fmt.Sprintf("lexer error: detected an infinite loop: %s", l.curPos))
			}
			l.logger.Debug("No loop detected")
		}
	}()

	l.input = input
	tokens := make([]token.Token, 0)
	for l.curPos.Pos < len(l.input) {
		matched := false
		ch := l.currentChar()

		if l.config.SkipWhitespace {
			for IsWhitespace(ch) {
				ch = l.Next()
			}
		}

		if IsEOF(ch) {
			break
		}

		for _, rule := range l.config.Rules {
			l.logger.Debug("Attempting to match %q with char %q", rule.Name, ch)

			if rule.Match(ch) {
				l.logger.Debug("Running action %q", rule.Name)

				tok, err := rule.Action(l, ch)
				var lerr Error
				if errors.As(err, &lerr) {
					if lerr.GotoNextRule {
						l.logger.Debug("Received an error from %q, moving to next rule", rule.Name)
						if l.prevPos.IsValid() {
							_ = l.Prev()
						}
						continue
					} else {
						_ = lerr.Error()
						return nil, lerr
					}
				} else if err != nil {
					return nil, err
				}

				if !tok.Start.IsValid() || !tok.End.IsValid() {
					return nil, fmt.Errorf("lexer error: start and/or end position is invalid (did you forget to start or end the rule?)")
				}

				if tok.Type == token.ILLEGAL {
					return nil, fmt.Errorf("lexer error: illegal token: %s: %q", tok, l.input[tok.Start.Pos:tok.End.Pos])
				}

				tokens = append(tokens, tok)
				matched = true
				break
			}
		}

		if !matched {
			// TODO(rdeusser): add output with line numbers and an up arrow at position.
			return nil, fmt.Errorf("lexer error: no rule to handle character at %s", l.curPos)
		}
	}

	return tokens, nil
}

// Lookahead how many runes ahead.
func (l *Lexer) Lookahead(n int) []rune {
	if n < 0 {
		return nil
	}
	if l.curPos.Pos+n >= len(l.input) {
		return []rune(l.input[l.curPos.Pos:])
	}
	return []rune(l.input[l.curPos.Pos : l.curPos.Pos+n])
}

func (l *Lexer) Next() rune {
	if l.curPos.Pos < len(l.input) {
		l.prevPos = l.curPos
	}

	if l.curPos.Line == 0 {
		l.curPos.Line = 1
	}

	ch := l.currentChar()
	if l.curPos.Pos+1 >= len(l.input) {
		ch = eof
	} else {
		ch = rune(l.input[l.curPos.Pos+1])
	}

	l.curPos.Pos++
	l.curPos.Column++

	if rune(l.input[l.prevPos.Pos]) != '\\' && rune(l.input[l.prevPos.Pos]) != '\'' {
		if ch == '\n' || ch == '\r' {
			l.curPos.Column = 1
			l.curPos.Line++
		}
	}

	return ch
}

func (l *Lexer) Prev() rune {
	l.curPos = l.prevPos
	return l.currentChar()
}

// StartRule is a helper method for starting a lex rule.
func (l *Lexer) StartRule(tokenType token.TokenType) token.Token {
	return token.Token{
		Type:  tokenType,
		Start: l.currentPos(),
	}
}

// EndRule is a helper method for ending a lex rule.
func (l *Lexer) EndRule(tok token.Token, err error) (token.Token, error) {
	tok.End = l.currentPos()
	if tok.End.Pos >= len(l.input) {
		tok.Literal = l.input[tok.Start.Pos:len(l.input)]
	} else {
		tok.Literal = l.input[tok.Start.Pos:tok.End.Pos]
	}
	// If the user already set the type, we shouldn't try to look it up because we can't look up things like strings. Only operators and keywords.
	if typ := l.LookupToken(tok.Literal); typ != token.ILLEGAL {
		tok.Type = typ
	}
	return tok, err
}

func (l *Lexer) LookupToken(literal string) token.TokenType {
	if t, ok := l.config.Operators[literal]; ok {
		return t
	}
	if t, ok := l.config.Keywords[literal]; ok {
		return t
	}
	return token.ILLEGAL
}

func (l *Lexer) currentChar() rune {
	if l.curPos.Pos >= len(l.input) {
		return eof
	} else {
		return rune(l.input[l.curPos.Pos])
	}
}

func (l *Lexer) currentPos() token.Position {
	return l.curPos
}
