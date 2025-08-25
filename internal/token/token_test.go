package token_test

import (
	"testing"

	"github.com/azraelsec/ippy/internal/token"
)

func TestTokenTypes(t *testing.T) {
	// Test that all token types are defined correctly
	tokenTypes := []token.Type{
		token.EOF,
		token.ILLEGAL,
		token.NUMBER,
		token.DASH,
		token.ASTERISK,
		token.COMMA,
	}

	expectedValues := []string{
		"EOF",
		"ILLEGAL",
		"NUMBER",
		"DASH",
		"ASTERISK",
		"COMMA",
	}

	if len(tokenTypes) != len(expectedValues) {
		t.Fatalf("Token types length mismatch: expected %d, got %d", len(expectedValues), len(tokenTypes))
	}

	for i, tokenType := range tokenTypes {
		if string(tokenType) != expectedValues[i] {
			t.Errorf("Token type mismatch at index %d: expected %s, got %s", i, expectedValues[i], string(tokenType))
		}
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		tokenType token.Type
		literal   string
	}{
		{
			name:      "NUMBER token",
			tokenType: token.NUMBER,
			literal:   "123",
		},
		{
			name:      "DASH token",
			tokenType: token.DASH,
			literal:   "-",
		},
		{
			name:      "ASTERISK token",
			tokenType: token.ASTERISK,
			literal:   "*",
		},
		{
			name:      "COMMA token",
			tokenType: token.COMMA,
			literal:   ",",
		},
		{
			name:      "EOF token",
			tokenType: token.EOF,
			literal:   "",
		},
		{
			name:      "ILLEGAL token",
			tokenType: token.ILLEGAL,
			literal:   "@",
		},
		{
			name:      "Empty literal",
			tokenType: token.NUMBER,
			literal:   "",
		},
		{
			name:      "Multi-character literal",
			tokenType: token.NUMBER,
			literal:   "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := token.New(tt.tokenType, tt.literal)

			if tok.Type != tt.tokenType {
				t.Errorf("Token type mismatch: expected %s, got %s", tt.tokenType, tok.Type)
			}

			if tok.Literal != tt.literal {
				t.Errorf("Token literal mismatch: expected %s, got %s", tt.literal, tok.Literal)
			}
		})
	}
}

func TestToken_Struct(t *testing.T) {
	// Test that Token struct has expected fields
	tok := token.Token{
		Type:    token.NUMBER,
		Literal: "123",
	}

	if tok.Type != token.NUMBER {
		t.Errorf("Direct field access failed: expected %s, got %s", token.NUMBER, tok.Type)
	}

	if tok.Literal != "123" {
		t.Errorf("Direct field access failed: expected '123', got '%s'", tok.Literal)
	}
}

func TestToken_ZeroValue(t *testing.T) {
	// Test zero value of Token
	var tok token.Token

	if tok.Type != "" {
		t.Errorf("Zero value Type should be empty string, got %s", tok.Type)
	}

	if tok.Literal != "" {
		t.Errorf("Zero value Literal should be empty string, got %s", tok.Literal)
	}
}

// Benchmark tests
func BenchmarkToken_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = token.New(token.NUMBER, "123")
	}
}

func BenchmarkToken_NewLongLiteral(b *testing.B) {
	longLiteral := "123456789012345678901234567890"

	for i := 0; i < b.N; i++ {
		_ = token.New(token.NUMBER, longLiteral)
	}
}
