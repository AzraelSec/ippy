package ip_test

import (
	"net"
	"testing"

	"github.com/azraelsec/ippy/internal/ip"
)

func TestParseIP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    net.IP
		wantErr bool
	}{
		// Valid IP addresses
		{
			name:    "valid IP - basic",
			input:   "192.168.1.1",
			want:    net.IPv4(192, 168, 1, 1),
			wantErr: false,
		},
		{
			name:    "valid IP - localhost",
			input:   "127.0.0.1",
			want:    net.IPv4(127, 0, 0, 1),
			wantErr: false,
		},
		{
			name:    "valid IP - all zeros",
			input:   "0.0.0.0",
			want:    net.IPv4(0, 0, 0, 0),
			wantErr: false,
		},
		{
			name:    "valid IP - all 255s",
			input:   "255.255.255.255",
			want:    net.IPv4(255, 255, 255, 255),
			wantErr: false,
		},
		{
			name:    "valid IP - mixed values",
			input:   "10.0.1.255",
			want:    net.IPv4(10, 0, 1, 255),
			wantErr: false,
		},
		{
			name:    "valid IP - private class A",
			input:   "10.1.2.3",
			want:    net.IPv4(10, 1, 2, 3),
			wantErr: false,
		},
		{
			name:    "valid IP - private class B",
			input:   "172.16.0.1",
			want:    net.IPv4(172, 16, 0, 1),
			wantErr: false,
		},
		{
			name:    "valid IP - private class C",
			input:   "192.168.0.1",
			want:    net.IPv4(192, 168, 0, 1),
			wantErr: false,
		},

		// Invalid IP addresses - wrong number of octets
		{
			name:    "invalid IP - too few octets (3)",
			input:   "192.168.1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - too few octets (2)",
			input:   "192.168",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - too few octets (1)",
			input:   "192",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - too many octets",
			input:   "192.168.1.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - way too many octets",
			input:   "192.168.1.1.1.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},

		// Invalid IP addresses - out of range values
		{
			name:    "invalid IP - octet too large (256)",
			input:   "192.168.1.256",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - octet way too large",
			input:   "192.168.1.999",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - first octet too large",
			input:   "256.168.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - middle octet too large",
			input:   "192.256.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - multiple octets too large",
			input:   "256.256.256.256",
			want:    net.IPv4zero,
			wantErr: true,
		},

		// Invalid IP addresses - non-numeric values
		{
			name:    "invalid IP - alphabetic characters",
			input:   "192.168.1.abc",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - mixed alphanumeric",
			input:   "192.168.1.1a",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - special characters",
			input:   "192.168.1.@",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - spaces",
			input:   "192.168.1. 1",
			want:    net.IPv4zero,
			wantErr: true,
		},

		// Invalid IP addresses - negative values
		{
			name:    "invalid IP - negative value",
			input:   "192.168.1.-1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - negative sign only",
			input:   "192.168.1.-",
			want:    net.IPv4zero,
			wantErr: true,
		},

		// Invalid IP addresses - empty values
		{
			name:    "invalid IP - empty string",
			input:   "",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - empty octet",
			input:   "192.168..1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - trailing dot",
			input:   "192.168.1.1.",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - leading dot",
			input:   ".192.168.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},
		{
			name:    "invalid IP - multiple dots",
			input:   "192..168.1.1",
			want:    net.IPv4zero,
			wantErr: true,
		},

		// Edge cases with leading zeros
		{
			name:    "valid IP - leading zeros",
			input:   "192.168.001.001",
			want:    net.IPv4(192, 168, 1, 1),
			wantErr: false,
		},
		{
			name:    "valid IP - single zero",
			input:   "192.168.0.0",
			want:    net.IPv4(192, 168, 0, 0),
			wantErr: false,
		},

		// Invalid hexadecimal (should be treated as invalid since we expect decimal)
		{
			name:    "invalid IP - hexadecimal notation",
			input:   "0xFF.0xFF.0xFF.0xFF",
			want:    net.IPv4zero,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ip.Parse(tt.input)

			if tt.wantErr && err == nil {
				t.Errorf("ParseIP() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ParseIP() unexpected error: %v", err)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test the IPv4 type alias functionality
func TestIPv4Type(t *testing.T) {
	// Test that we can create and use IPv4 values
	ipv4 := net.IPv4(192, 168, 1, 1).To4()

	// Test array access
	if ipv4[0] != 192 {
		t.Errorf("IPv4[0] = %d, want 192", ipv4[0])
	}
	if ipv4[1] != 168 {
		t.Errorf("IPv4[1] = %d, want 168", ipv4[1])
	}
	if ipv4[2] != 1 {
		t.Errorf("IPv4[2] = %d, want 1", ipv4[2])
	}
	if ipv4[3] != 1 {
		t.Errorf("IPv4[3] = %d, want 1", ipv4[3])
	}

	// Test that length is 4
	if len(ipv4) != 4 {
		t.Errorf("len(IPv4) = %d, want 4", len(ipv4))
	}
}

// Benchmark tests
func BenchmarkParseIP_Valid(b *testing.B) {
	vip := "192.168.1.100"

	for b.Loop() {
		if _, err := ip.Parse(vip); err != nil {
			b.Fatalf("ParseIP failed: %v", err)
		}
	}
}

func BenchmarkParseIP_Invalid(b *testing.B) {
	invalidIP := "192.168.1.256"

	for b.Loop() {
		if _, err := ip.Parse(invalidIP); err == nil {
			b.Fatal("ParseIP should have failed")
		}
	}
}
