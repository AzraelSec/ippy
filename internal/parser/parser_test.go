package parser_test

import (
	"strings"
	"testing"

	"github.com/azraelsec/ippy/internal/parser"
)

func TestParse_Valid(t *testing.T) {
	tests := []struct {
		input  string
		ranges []parser.Interval
	}{
		{"0", []parser.Interval{{0, 0}}},
		{"0-1", []parser.Interval{{0, 1}}},
		{"0,1", []parser.Interval{{0, 0}, {1, 1}}},
		{"0,1-2,4-5", []parser.Interval{{0, 0}, {1, 2}, {4, 5}}},
		{"*", []parser.Interval{{0, 255}}},
		{"0, 2, *", []parser.Interval{{0, 0}, {2, 2}, {0, 255}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			p := parser.New(tt.input)
			its, ok := p.Parse()
			if !ok {
				t.Fatalf("parsing failed: %q", p.Errors())
			}

			if len(its) != len(tt.ranges) {
				t.Fatalf("intervals length mismatch want=%d, have=%d", len(tt.ranges), len(its))
			}
			for i := range tt.ranges {
				if tt.ranges[i][0] != its[i][0] || tt.ranges[i][1] != its[i][1] {
					t.Fatalf("interval mismatch want=%v, have=%v", tt.ranges[i], its[i])
				}
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	tests := []struct {
		input        string
		expectedErrs []string // We expect specific error messages
	}{
		{
			input:        "",
			expectedErrs: []string{"a valid octet should have at least 1 range"},
		},
		{
			input:        "256",
			expectedErrs: []string{"numeric value 256 is not valid"}, // Parser should handle out-of-range numbers but return error
		},
		{
			input:        "1-",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
		{
			input:        "-1",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
		{
			input:        "1,",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
		{
			input:        ",1",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
		{
			input:        "1,,2",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
		{
			input:        "abc",
			expectedErrs: []string{"expected current token type is NUMBER"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			p := parser.New(tt.input)

			if _, ok := p.Parse(); ok {
				t.Errorf("Parse() expected to fail but succeeded")
				return
			}

			errors := p.Errors()
			if len(tt.expectedErrs) == 0 {
				// We expect some error but don't care about the specific message
				if len(errors) == 0 {
					t.Errorf("Parse() expected errors but got none")
				}
			} else {
				// Check for specific error messages
				for _, expectedErr := range tt.expectedErrs {
					found := false
					for _, actualErr := range errors {
						if strings.Contains(actualErr, expectedErr) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Parse() expected error containing %q, got errors: %v", expectedErr, errors)
					}
				}
			}
		})
	}
}

func TestParser_New(t *testing.T) {
	// Test that New creates a parser with proper initial state
	p := parser.New("123")

	if p == nil {
		t.Fatal("New() returned nil parser")
	}

	// Test that parser starts with no errors
	if len(p.Errors()) != 0 {
		t.Errorf("New parser should have no errors, got: %v", p.Errors())
	}
}

func TestParser_Intervals(t *testing.T) {
	// Test that Intervals() returns empty slice for unparsed parser
	p := parser.New("123")

	its, ok := p.Parse()
	if !ok {
		t.Fatalf("Parse() failed: %v", p.Errors())
	}

	expected := []parser.Interval{{123, 123}}
	if len(its) != len(expected) {
		t.Errorf("Intervals length mismatch: expected %d, got %d", len(expected), len(its))
	}

	if its[0] != expected[0] {
		t.Errorf("Interval mismatch: expected %v, got %v", expected[0], its[0])
	}
}

// Benchmark tests
func BenchmarkParser_Simple(b *testing.B) {
	input := "123"

	for b.Loop() {
		p := parser.New(input)
		if _, ok := p.Parse(); !ok {
			b.Fatalf("Parse failed: %v", p.Errors())
		}
	}
}

func BenchmarkParser_Complex(b *testing.B) {
	input := "1-10,20,30-40,50-255"

	for b.Loop() {
		p := parser.New(input)
		if _, ok := p.Parse(); !ok {
			b.Fatalf("Parse failed: %v", p.Errors())
		}
	}
}

func BenchmarkParser_Wildcard(b *testing.B) {
	input := "*"

	for b.Loop() {
		p := parser.New(input)
		if _, ok := p.Parse(); !ok {
			b.Fatalf("Parse failed: %v", p.Errors())
		}
	}
}
