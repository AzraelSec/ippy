// Package lexer provides lexical analysis functionality for tokenizing IP octet expressions.
// It converts input strings into a sequence of tokens that can be parsed by the parser package.
// The lexer supports numbers, dashes, asterisks, commas, and handles whitespace appropriately.
package lexer

import "github.com/azraelsec/ippy/internal/token"

const nul = 0

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(s string) *Lexer {
	l := &Lexer{input: s}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = nul
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tkn token.Token

	l.skipWhiteSpaces()

	switch l.ch {
	case '-':
		tkn = token.New(token.DASH, string(l.ch))
	case '*':
		tkn = token.New(token.ASTERISK, string(l.ch))
	case ',':
		tkn = token.New(token.COMMA, string(l.ch))
	case nul:
		tkn = token.New(token.EOF, "")
	default:
		if isDigit(l.ch) {
			tkn = token.New(token.NUMBER, l.readNumber())
			return tkn
		}
		tkn = token.New(token.ILLEGAL, string(l.ch))
	}

	l.readChar()
	return tkn
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch == ' ' {
		l.readChar()
	}
}
