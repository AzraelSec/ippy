package lexer_test

import (
	"testing"

	"github.com/azraelsec/ippy/internal/lexer"
	"github.com/azraelsec/ippy/internal/token"
)

type tokenTestCase struct {
	expectedType    token.Type
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []tokenTestCase
	}{
		{input: "123", expectedTokens: []tokenTestCase{
			{token.NUMBER, "123"},
		}},
		{input: "1-2", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.DASH, "-"},
			{token.NUMBER, "2"},
		}},
		{input: "*", expectedTokens: []tokenTestCase{
			{token.ASTERISK, "*"},
		}},
		{input: "1,2-*", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.COMMA, ","},
			{token.NUMBER, "2"},
			{token.DASH, "-"},
			{token.ASTERISK, "*"},
		}},
		{input: "1,2,3", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.COMMA, ","},
			{token.NUMBER, "2"},
			{token.COMMA, ","},
			{token.NUMBER, "3"},
		}},

		// Edge cases
		{input: "255", expectedTokens: []tokenTestCase{
			{token.NUMBER, "255"},
		}},
		{input: "0", expectedTokens: []tokenTestCase{
			{token.NUMBER, "0"},
		}},
		{input: "000", expectedTokens: []tokenTestCase{
			{token.NUMBER, "000"},
		}},
		{input: "123456", expectedTokens: []tokenTestCase{
			{token.NUMBER, "123456"},
		}},

		// Whitespace handling
		{input: " 1 ", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
		}},
		{input: "1 - 2", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.DASH, "-"},
			{token.NUMBER, "2"},
		}},
		{input: " * ", expectedTokens: []tokenTestCase{
			{token.ASTERISK, "*"},
		}},
		{input: "1, 2, 3", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.COMMA, ","},
			{token.NUMBER, "2"},
			{token.COMMA, ","},
			{token.NUMBER, "3"},
		}},

		// Complex expressions
		{input: "1-10,20,30-40", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.DASH, "-"},
			{token.NUMBER, "10"},
			{token.COMMA, ","},
			{token.NUMBER, "20"},
			{token.COMMA, ","},
			{token.NUMBER, "30"},
			{token.DASH, "-"},
			{token.NUMBER, "40"},
		}},
		{input: "*,1-5,10", expectedTokens: []tokenTestCase{
			{token.ASTERISK, "*"},
			{token.COMMA, ","},
			{token.NUMBER, "1"},
			{token.DASH, "-"},
			{token.NUMBER, "5"},
			{token.COMMA, ","},
			{token.NUMBER, "10"},
		}},

		// Empty input
		{input: "", expectedTokens: []tokenTestCase{}},

		// Invalid characters (should produce ILLEGAL tokens)
		{input: "abc", expectedTokens: []tokenTestCase{
			{token.ILLEGAL, "a"},
			{token.ILLEGAL, "b"},
			{token.ILLEGAL, "c"},
		}},
		{input: "1@2", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.ILLEGAL, "@"},
			{token.NUMBER, "2"},
		}},
		{input: "1.2", expectedTokens: []tokenTestCase{
			{token.NUMBER, "1"},
			{token.ILLEGAL, "."},
			{token.NUMBER, "2"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			for _, expTkn := range tt.expectedTokens {
				tkn := l.NextToken()
				if tkn.Type != expTkn.expectedType {
					t.Fatalf("wrong token type expected=%q, got=%q", expTkn.expectedType, tkn.Type)
				}
				if tkn.Literal != expTkn.expectedLiteral {
					t.Fatalf("wrong token type expected=%q, got=%q", expTkn.expectedLiteral, tkn.Literal)
				}
			}
		})
	}
}

func TestLexer_EdgeCases(t *testing.T) {
	// Test EOF handling
	t.Run("EOF handling", func(t *testing.T) {
		l := lexer.New("")
		tkn := l.NextToken()
		if tkn.Type != token.EOF {
			t.Errorf("Expected EOF token, got %s", tkn.Type)
		}

		// Multiple calls to NextToken should continue returning EOF
		tkn = l.NextToken()
		if tkn.Type != token.EOF {
			t.Errorf("Expected EOF token on second call, got %s", tkn.Type)
		}
	})

	// Test that all tokens end with EOF
	t.Run("All inputs end with EOF", func(t *testing.T) {
		inputs := []string{"1", "*", "1-2", "1,2", "1-2,3-4"}

		for _, input := range inputs {
			l := lexer.New(input)
			var tokens []token.Token

			for {
				tkn := l.NextToken()
				tokens = append(tokens, tkn)
				if tkn.Type == token.EOF {
					break
				}
			}

			if len(tokens) == 0 {
				t.Errorf("No tokens generated for input %s", input)
				continue
			}

			lastToken := tokens[len(tokens)-1]
			if lastToken.Type != token.EOF {
				t.Errorf("Last token for input %s was not EOF, got %s", input, lastToken.Type)
			}
		}
	})
}

func TestLexer_TokenNew(t *testing.T) {
	// Test the token.New function
	tkn := token.New(token.NUMBER, "123")
	if tkn.Type != token.NUMBER {
		t.Errorf("Token type mismatch: expected %s, got %s", token.NUMBER, tkn.Type)
	}
	if tkn.Literal != "123" {
		t.Errorf("Token literal mismatch: expected '123', got '%s'", tkn.Literal)
	}
}

// Benchmark tests
func BenchmarkLexer_SimpleNumber(b *testing.B) {
	input := "123"

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		for {
			tkn := l.NextToken()
			if tkn.Type == token.EOF {
				break
			}
		}
	}
}

func BenchmarkLexer_ComplexExpression(b *testing.B) {
	input := "1-10,20,30-40,*,50-255"

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		for {
			tkn := l.NextToken()
			if tkn.Type == token.EOF {
				break
			}
		}
	}
}

func BenchmarkLexer_WithWhitespace(b *testing.B) {
	input := " 1 - 10 , 20 , 30 - 40 "

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		for {
			tkn := l.NextToken()
			if tkn.Type == token.EOF {
				break
			}
		}
	}
}
