package lexer

import (
	"fmt"

	"github.com/rdeusser/parsekit/token"
)

// Action is used to lex input.
type Action func(l *Lexer, ch rune) (tok token.Token, err error)

// LexIdentifier lexes an identifier.
func LexIdentifier(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.IDENT)

	// This can't be a number if you use the parsekit command-line tool because the matcher is set
	// to letters only (as the starting character). This is here to demonstrate how to pass errors
	// to the lexer and what you want the lexer to do with it (i.e. move to the next rule or bail).
	if IsNumber(ch) {
		return l.EndRule(tok, Error{Lexer: l, Msg: "first character of an identifier can't be a number", GotoNextRule: true})
	}

	for IsLetter(ch) || IsNumber(ch) {
		ch = l.Next()
	}

	return l.EndRule(tok, nil)
}

func LexChar(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.CHAR)

	ch = l.Next()
	ch = l.Next()

	if !IsSingleQuote(ch) {
		return l.EndRule(tok, Error{Lexer: l, Msg: "chars cannot be more than one character"})
	}

	ch = l.Next()

	return l.EndRule(tok, nil)
}

func LexString(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.STRING)

	prevCh := ch
	for {
		prevCh = ch
		ch = l.Next()

		if IsNewline(ch) {
			return l.EndRule(tok, Error{Lexer: l, Msg: "literal newlines aren't valid in a string"})
		}

		if IsDoubleQuote(ch) {
			if prevCh == '\\' { // is the quote escaped?
				prevCh = ch
				ch = l.Next()
			} else {
				ch = l.Next()
				break
			}
		}
	}

	tok, err := l.EndRule(tok, nil)

	return tok, err
}

func LexRawString(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.STRING)

	for {
		ch = l.Next()
		if IsBackQuote(ch) {
			ch = l.Next()
			break
		}
	}

	return l.EndRule(tok, nil)
}

func LexNumber(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.NUMBER)

	for IsNumber(ch) {
		ch = l.Next()
		if ch == '.' {
			tok.Type = token.FLOAT
			ch = l.Next()
		}
	}

	return l.EndRule(tok, nil)
}

func LexOperator(l *Lexer, ch rune) (token.Token, error) {
	tok := l.StartRule(token.ILLEGAL)
	maxLen := longestOperator(l.config.Operators)
	op := l.Lookahead(maxLen)

	for {
		tok.Type = l.LookupToken(string(op))
		op = l.Lookahead(maxLen)
		if maxLen < 0 {
			return l.EndRule(tok, fmt.Errorf("invalid operator: %s", string(ch)))
		}
		if tok.Type == token.ILLEGAL {
			maxLen--
			continue
		}
		break
	}

	for i := 0; i < len(op); i++ {
		_ = l.Next()
	}

	return l.EndRule(tok, nil)
}

func longestOperator(m map[string]token.TokenType) int {
	result := ""
	for s := range m {
		if len(s) > len(result) {
			result = s
		}
	}
	return len(result)
}
