// Package parser provides functionality for parsing IP octet expressions into intervals.
// It takes tokenized input from the lexer and converts it into structured interval representations
// that can be used for IP address range validation and processing.
package parser

import (
	"fmt"
	"strconv"

	"github.com/azraelsec/ippy/internal/lexer"
	"github.com/azraelsec/ippy/internal/token"
)

type Interval [2]byte

type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) expectCurrIs(t token.Type) bool {
	if !p.currTokenIs(t) {
		p.currError(t)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) currError(t token.Type) {
	msg := fmt.Sprintf("expected current token type is %s, found %s", t, p.currToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseExpr() ([]Interval, bool) {
	var intervals []Interval
	for !p.currTokenIs(token.EOF) {
		interval, ok := p.parseTerm()
		if !ok {
			return []Interval{}, false
		}
		intervals = append(intervals, interval)

		if !p.peekTokenIs(token.EOF) {
			p.expectCurrIs(token.COMMA)
		}
	}
	return intervals, true
}

func (p *Parser) Parse() ([]Interval, bool) {
	intervals, ok := p.parseExpr()
	if !ok {
		return []Interval{}, false
	}
	if len(intervals) == 0 {
		msg := "a valid octet should have at least 1 range"
		p.errors = append(p.errors, msg)
		return []Interval{}, false
	}
	return intervals, true
}

func (p *Parser) parseTerm() (Interval, bool) {
	// TODO: handle x-* and *-x intervals
	if p.currTokenIs(token.ASTERISK) {
		p.nextToken()
		return Interval{0, 255}, true
	}

	start, ok := p.parseNumber()
	if !ok {
		p.numberParsingError()
		return Interval{}, false
	}

	if !p.currTokenIs(token.DASH) {
		return Interval{start, start}, true
	}

	p.nextToken()
	end, ok := p.parseNumber()
	if !ok {
		p.numberParsingError()
		return Interval{}, false
	}

	return Interval{start, end}, true
}

func (p *Parser) parseNumber() (uint8, bool) {
	if !p.currTokenIs(token.NUMBER) {
		p.currError(token.NUMBER)
		return 0, false
	}

	num, err := strconv.Atoi(p.currToken.Literal)
	if err != nil || num < 0 || num > 255 {
		p.numberParsingError()
		return 0, false
	}

	p.nextToken()
	return uint8(num), true
}

func (p *Parser) numberParsingError() {
	msg := fmt.Sprintf("numeric value %s is not valid", p.currToken.Literal)
	p.errors = append(p.errors, msg)
}

func New(s string) *Parser {
	p := &Parser{
		l:      lexer.New(s),
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}
