// Package token defines the token types and structures used by the lexer
// for tokenizing IP octet expressions. It provides constants for different
// token types like numbers, dashes, asterisks, and commas, along with a
// Token struct to represent individual tokens with their type and literal value.
package token

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	NUMBER   = "NUMBER"
	DASH     = "DASH"
	ASTERISK = "ASTERISK"
	COMMA    = "COMMA"
)

type Type = string

type Token struct {
	Type    Type
	Literal string
}

func New(tp Type, l string) Token {
	return Token{
		Type:    tp,
		Literal: l,
	}
}
