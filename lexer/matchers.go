package lexer

import "unicode"

type Matcher func(rune) bool

func IsLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func IsDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func IsNumber(ch rune) bool {
	return unicode.IsNumber(ch)
}

func IsSingleQuote(ch rune) bool {
	return ch == '\''
}

func IsDoubleQuote(ch rune) bool {
	return ch == '"'
}

func IsBackQuote(ch rune) bool {
	return ch == '`'
}

func IsOperator(ch rune) bool {
	return IsSymbol(ch) || IsPunct(ch)
}

func IsSymbol(ch rune) bool {
	return unicode.IsSymbol(ch)
}

func IsPunct(ch rune) bool {
	return unicode.IsPunct(ch)
}

func IsWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func IsSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func IsNewline(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func IsEOF(ch rune) bool {
	return ch == 0 || ch == rune(0)
}
