package ipexpr_test

import (
	"net"
	"testing"

	"github.com/azraelsec/ippy/internal/ip"
	"github.com/azraelsec/ippy/pkg/ipexpr"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		// Valid expressions
		{
			name:    "valid simple expression",
			expr:    "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "valid wildcard expression",
			expr:    "192.168.1.*",
			wantErr: false,
		},
		{
			name:    "valid range expression",
			expr:    "192.168.1.1-10",
			wantErr: false,
		},
		{
			name:    "valid comma-separated expression",
			expr:    "192.168.1.1,2,3",
			wantErr: false,
		},
		{
			name:    "valid complex expression",
			expr:    "192.168.1-5.1,10-20,100",
			wantErr: false,
		},
		{
			name:    "valid all wildcards",
			expr:    "*.*.*.*",
			wantErr: false,
		},
		{
			name:    "valid edge values",
			expr:    "0.0.0.0",
			wantErr: false,
		},
		{
			name:    "valid max values",
			expr:    "255.255.255.255",
			wantErr: false,
		},
		{
			name:    "valid full range",
			expr:    "0-255.0-255.0-255.0-255",
			wantErr: false,
		},

		// Invalid expressions - wrong number of octets
		{
			name:    "invalid - too few octets",
			expr:    "192.168.1",
			wantErr: true,
		},
		{
			name:    "invalid - too many octets",
			expr:    "192.168.1.1.1",
			wantErr: true,
		},
		{
			name:    "invalid - empty expression",
			expr:    "",
			wantErr: true,
		},

		// Invalid expressions - malformed octets
		{
			name:    "invalid - malformed range",
			expr:    "192.168.1.1-",
			wantErr: true,
		},
		{
			name:    "invalid - out of range value",
			expr:    "192.168.1.256",
			wantErr: true,
		},
		{
			name:    "invalid - non-numeric",
			expr:    "192.168.1.abc",
			wantErr: true,
		},
		{
			name:    "invalid - empty octet",
			expr:    "192.168..1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ipexpr.Parse(tt.expr)

			if tt.wantErr && err == nil {
				t.Errorf("Parse() expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}
		})
	}
}

func TestIPExpr_Matches(t *testing.T) {
	tests := []struct {
		name string
		expr string
		ip   string
		want bool
	}{
		// Exact matches
		{
			name: "exact match - success",
			expr: "192.168.1.1",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "exact match - failure",
			expr: "192.168.1.1",
			ip:   "192.168.1.2",
			want: false,
		},

		// Wildcard matches
		{
			name: "wildcard last octet - success",
			expr: "192.168.1.*",
			ip:   "192.168.1.100",
			want: true,
		},
		{
			name: "wildcard last octet - success at boundary",
			expr: "192.168.1.*",
			ip:   "192.168.1.255",
			want: true,
		},
		{
			name: "wildcard last octet - success at zero",
			expr: "192.168.1.*",
			ip:   "192.168.1.0",
			want: true,
		},
		{
			name: "wildcard last octet - failure",
			expr: "192.168.1.*",
			ip:   "192.168.2.100",
			want: false,
		},
		{
			name: "wildcard first octet - success",
			expr: "*.168.1.1",
			ip:   "10.168.1.1",
			want: true,
		},
		{
			name: "wildcard multiple octets - success",
			expr: "192.*.*.1",
			ip:   "192.100.50.1",
			want: true,
		},
		{
			name: "wildcard all octets - success",
			expr: "*.*.*.*",
			ip:   "1.2.3.4",
			want: true,
		},

		// Range matches
		{
			name: "range match - within range",
			expr: "192.168.1.1-10",
			ip:   "192.168.1.5",
			want: true,
		},
		{
			name: "range match - at start boundary",
			expr: "192.168.1.1-10",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "range match - at end boundary",
			expr: "192.168.1.1-10",
			ip:   "192.168.1.10",
			want: true,
		},
		{
			name: "range match - below range",
			expr: "192.168.1.5-10",
			ip:   "192.168.1.4",
			want: false,
		},
		{
			name: "range match - above range",
			expr: "192.168.1.1-10",
			ip:   "192.168.1.11",
			want: false,
		},
		{
			name: "range match - different octet",
			expr: "192.168.1.1-10",
			ip:   "192.168.2.5",
			want: false,
		},

		// Comma-separated matches
		{
			name: "comma-separated - first value",
			expr: "192.168.1.1,5,10",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "comma-separated - middle value",
			expr: "192.168.1.1,5,10",
			ip:   "192.168.1.5",
			want: true,
		},
		{
			name: "comma-separated - last value",
			expr: "192.168.1.1,5,10",
			ip:   "192.168.1.10",
			want: true,
		},
		{
			name: "comma-separated - no match",
			expr: "192.168.1.1,5,10",
			ip:   "192.168.1.2",
			want: false,
		},

		// Complex patterns
		{
			name: "complex pattern - range and comma",
			expr: "192.168.1.1-5,10,20-25",
			ip:   "192.168.1.3",
			want: true,
		},
		{
			name: "complex pattern - exact match in comma list",
			expr: "192.168.1.1-5,10,20-25",
			ip:   "192.168.1.10",
			want: true,
		},
		{
			name: "complex pattern - range match",
			expr: "192.168.1.1-5,10,20-25",
			ip:   "192.168.1.22",
			want: true,
		},
		{
			name: "complex pattern - no match",
			expr: "192.168.1.1-5,10,20-25",
			ip:   "192.168.1.15",
			want: false,
		},

		// Multiple octets with patterns
		{
			name: "multiple octet patterns - success",
			expr: "192.168.1-5.1,10-20",
			ip:   "192.168.3.15",
			want: true,
		},
		{
			name: "multiple octet patterns - third octet no match",
			expr: "192.168.1-5.1,10-20",
			ip:   "192.168.6.15",
			want: false,
		},
		{
			name: "multiple octet patterns - fourth octet no match",
			expr: "192.168.1-5.1,10-20",
			ip:   "192.168.3.25",
			want: false,
		},

		// Edge cases
		{
			name: "edge case - all zeros",
			expr: "0.0.0.0",
			ip:   "0.0.0.0",
			want: true,
		},
		{
			name: "edge case - all 255s",
			expr: "255.255.255.255",
			ip:   "255.255.255.255",
			want: true,
		},
		{
			name: "edge case - full range match",
			expr: "0-255.0-255.0-255.0-255",
			ip:   "127.0.0.1",
			want: true,
		},
		{
			name: "edge case - single value ranges",
			expr: "192-192.168-168.1-1.1-1",
			ip:   "192.168.1.1",
			want: true,
		},

		// Test with spaces (should be handled by lexer)
		{
			name: "spaces in comma list",
			expr: "192.168.1.1, 2, 3",
			ip:   "192.168.1.2",
			want: true,
		},
		{
			name: "spaces in wildcard range",
			expr: "192.168.1. *",
			ip:   "192.168.1.100",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipExpr, err := ipexpr.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}

			got, _ := ipExpr.Matches(tt.ip)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPExpr_Generate(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want []ip.IPv4
	}{
		// Single IP address
		{
			name: "single ip",
			expr: "192.168.1.1",
			want: []ip.IPv4{net.IPv4(192, 168, 1, 1)},
		},
		// Simple ranges
		{
			name: "range on the last octet",
			expr: "192.168.1.1-3",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 2),
				net.IPv4(192, 168, 1, 3),
			},
		},
		{
			name: "range on the third octet",
			expr: "10.0.1-3.1",
			want: []ip.IPv4{
				net.IPv4(10, 0, 1, 1),
				net.IPv4(10, 0, 2, 1),
				net.IPv4(10, 0, 3, 1),
			},
		},
		{
			name: "range on the second octet",
			expr: "172.16-17.0.1",
			want: []ip.IPv4{
				net.IPv4(172, 16, 0, 1),
				net.IPv4(172, 17, 0, 1),
			},
		},
		{
			name: "range on the first octet",
			expr: "10-11.0.0.1",
			want: []ip.IPv4{
				net.IPv4(10, 0, 0, 1),
				net.IPv4(11, 0, 0, 1),
			},
		},

		// Comma-separated values
		{
			name: "comma-separated last octet",
			expr: "192.168.1.1,3,5",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 3),
				net.IPv4(192, 168, 1, 5),
			},
		},
		{
			name: "comma-separated third octet",
			expr: "10.0.1,3,5.1",
			want: []ip.IPv4{
				net.IPv4(10, 0, 1, 1),
				net.IPv4(10, 0, 3, 1),
				net.IPv4(10, 0, 5, 1),
			},
		},

		// Wildcards
		{
			name: "wildcard last octet - limited by Generate implementation",
			expr: "192.168.1.*",
			want: func() []ip.IPv4 {
				var ips []ip.IPv4
				for i := 0; i <= 255; i++ {
					ips = append(ips, net.IPv4(192, 168, 1, byte(i)))
				}
				return ips
			}(),
		},
		{
			name: "wildcard third octet - first few values",
			expr: "10.0.*.1",
			want: func() []ip.IPv4 {
				var ips []ip.IPv4
				for i := 0; i <= 255; i++ {
					ips = append(ips, net.IPv4(10, 0, byte(i), 1))
				}
				return ips
			}(),
		},

		// Mixed patterns with ranges and commas
		{
			name: "range and comma in last octet",
			expr: "192.168.1.1-3,10,20-21",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 2),
				net.IPv4(192, 168, 1, 3),
				net.IPv4(192, 168, 1, 10),
				net.IPv4(192, 168, 1, 20),
				net.IPv4(192, 168, 1, 21),
			},
		},
		{
			name: "mixed patterns multiple octets",
			expr: "10.0-1.1,3.1-2",
			want: []ip.IPv4{
				net.IPv4(10, 0, 1, 1),
				net.IPv4(10, 0, 1, 2),
				net.IPv4(10, 0, 3, 1),
				net.IPv4(10, 0, 3, 2),
				net.IPv4(10, 1, 1, 1),
				net.IPv4(10, 1, 1, 2),
				net.IPv4(10, 1, 3, 1),
				net.IPv4(10, 1, 3, 2),
			},
		},

		// Edge cases
		{
			name: "all zeros",
			expr: "0.0.0.0",
			want: []ip.IPv4{net.IPv4(0, 0, 0, 0)},
		},
		{
			name: "all 255s",
			expr: "255.255.255.255",
			want: []ip.IPv4{net.IPv4(255, 255, 255, 255)},
		},
		{
			name: "single value range",
			expr: "192.168.1.5-5",
			want: []ip.IPv4{net.IPv4(192, 168, 1, 5)},
		},
		{
			name: "range at boundaries",
			expr: "192.168.1.254-255",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 254),
				net.IPv4(192, 168, 1, 255),
			},
		},
		{
			name: "range at zero boundary",
			expr: "192.168.1.0-1",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 0),
				net.IPv4(192, 168, 1, 1),
			},
		},

		// Complex multi-octet patterns
		{
			name: "ranges in multiple octets",
			expr: "10.0-1.0-1.1-2",
			want: []ip.IPv4{
				net.IPv4(10, 0, 0, 1),
				net.IPv4(10, 0, 0, 2),
				net.IPv4(10, 0, 1, 1),
				net.IPv4(10, 0, 1, 2),
				net.IPv4(10, 1, 0, 1),
				net.IPv4(10, 1, 0, 2),
				net.IPv4(10, 1, 1, 1),
				net.IPv4(10, 1, 1, 2),
			},
		},
		{
			name: "comma separated in multiple octets",
			expr: "192.168.1,2.1,2",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 2),
				net.IPv4(192, 168, 2, 1),
				net.IPv4(192, 168, 2, 2),
			},
		},

		// Mixed patterns across octets
		{
			name: "wildcard and specific values",
			expr: "192.168.*.1",
			want: func() []ip.IPv4 {
				var ips []ip.IPv4
				for i := 0; i <= 255; i++ {
					ips = append(ips, net.IPv4(192, 168, byte(i), 1))
				}
				return ips
			}(),
		},

		// Small test cases for verification
		{
			name: "two values in different octets",
			expr: "10.0,1.0.1",
			want: []ip.IPv4{
				net.IPv4(10, 0, 0, 1),
				net.IPv4(10, 1, 0, 1),
			},
		},

		// Additional edge cases and ordering tests
		{
			name: "reverse order range",
			expr: "192.168.1.5-3", // Note: this might be invalid depending on parser
			want: []ip.IPv4{},     // Empty if invalid, or properly ordered if valid
		},
		{
			name: "overlapping comma and range",
			expr: "192.168.1.1-3,2,4",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 2),
				net.IPv4(192, 168, 1, 3),
				net.IPv4(192, 168, 1, 4),
			},
		},
		{
			name: "duplicate values in comma list",
			expr: "192.168.1.1,1,2",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 2),
			},
		},
		{
			name: "large range in first octet",
			expr: "1-3.168.1.1",
			want: []ip.IPv4{
				net.IPv4(1, 168, 1, 1),
				net.IPv4(2, 168, 1, 1),
				net.IPv4(3, 168, 1, 1),
			},
		},
		{
			name: "comma with wide spacing",
			expr: "192.168.1.1,50,100,200,255",
			want: []ip.IPv4{
				net.IPv4(192, 168, 1, 1),
				net.IPv4(192, 168, 1, 50),
				net.IPv4(192, 168, 1, 100),
				net.IPv4(192, 168, 1, 200),
				net.IPv4(192, 168, 1, 255),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipExpr, err := ipexpr.Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}

			for i, iip := range ipExpr.Generate() {
				if !tt.want[i].Equal(iip) {
					t.Errorf("Generate(%s)[%d] = %s, want %s", tt.expr, i, iip, tt.want[i])
				}
			}
		})
	}
}

// Test IPExpr methods separately
func TestIPExpr_MatchesMethod(t *testing.T) {
	// Test that we can call the Matches method directly
	ipExpr, err := ipexpr.Parse("192.168.1.*")
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Test various combinations
	testCases := []struct {
		ip   string
		want bool
	}{
		{"192.168.1.0", true},
		{"192.168.1.1", true},
		{"192.168.1.255", true},
		{"192.168.2.1", false},
		{"10.168.1.1", false},
	}

	for _, tc := range testCases {
		got, _ := ipExpr.Matches(tc.ip)
		if got != tc.want {
			t.Errorf("Matches(%s) = %v, want %v", tc.ip, got, tc.want)
		}
	}
}

// Benchmark tests
func BenchmarkParse_Simple(b *testing.B) {
	expr := "192.168.1.1"

	for b.Loop() {
		if _, err := ipexpr.Parse(expr); err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParse_Complex(b *testing.B) {
	expr := "192.168.1-5.1,10-20,100,200-255"

	for b.Loop() {
		if _, err := ipexpr.Parse(expr); err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkIPExpr_Matches_Simple(b *testing.B) {
	ipExpr, err := ipexpr.Parse("192.168.1.*")
	if err != nil {
		b.Fatalf("Parse failed: %v", err)
	}

	for b.Loop() {
		_, _ = ipExpr.Matches("192.168.1.100")
	}
}

func BenchmarkIPExpr_Matches_Complex(b *testing.B) {
	ipExpr, err := ipexpr.Parse("10-50.20-200.1,5,10-20,50-100,150-200,250.1-254")
	if err != nil {
		b.Fatalf("Parse failed: %v", err)
	}

	for b.Loop() {
		_, _ = ipExpr.Matches("25.100.75.100")
	}
}
