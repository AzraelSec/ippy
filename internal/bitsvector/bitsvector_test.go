package bitsvector

import (
	"testing"

	"github.com/azraelsec/ippy/internal/parser"
)

func TestNew_SingleInterval(t *testing.T) {
	tests := []struct {
		name      string
		intervals []parser.Interval
		testByte  byte
		expected  bool
	}{
		{
			name:      "single byte interval",
			intervals: []parser.Interval{{10, 10}},
			testByte:  10,
			expected:  true,
		},
		{
			name:      "single byte interval - test outside",
			intervals: []parser.Interval{{10, 10}},
			testByte:  11,
			expected:  false,
		},
		{
			name:      "range interval",
			intervals: []parser.Interval{{10, 15}},
			testByte:  12,
			expected:  true,
		},
		{
			name:      "range interval - start boundary",
			intervals: []parser.Interval{{10, 15}},
			testByte:  10,
			expected:  true,
		},
		{
			name:      "range interval - end boundary",
			intervals: []parser.Interval{{10, 15}},
			testByte:  15,
			expected:  true,
		},
		{
			name:      "range interval - outside range",
			intervals: []parser.Interval{{10, 15}},
			testByte:  16,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := New(tt.intervals)
			result := ob.Test(tt.testByte)
			if result != tt.expected {
				t.Errorf("Test(%d) = %v, expected %v", tt.testByte, result, tt.expected)
			}
		})
	}
}

func TestNew_MultipleIntervals(t *testing.T) {
	intervals := []parser.Interval{{1, 5}, {10, 15}, {20, 25}}
	ob := New(intervals)

	// Test bytes within intervals
	testCases := []struct {
		testByte byte
		expected bool
	}{
		{1, true},   // first interval start
		{3, true},   // first interval middle
		{5, true},   // first interval end
		{7, false},  // between intervals
		{10, true},  // second interval start
		{12, true},  // second interval middle
		{15, true},  // second interval end
		{18, false}, // between intervals
		{20, true},  // third interval start
		{22, true},  // third interval middle
		{25, true},  // third interval end
		{30, false}, // after all intervals
		{0, false},  // before all intervals
	}

	for _, tc := range testCases {
		t.Run(string(rune(tc.testByte)), func(t *testing.T) {
			result := ob.Test(tc.testByte)
			if result != tc.expected {
				t.Errorf("Test(%d) = %v, expected %v", tc.testByte, result, tc.expected)
			}
		})
	}
}

func TestNew_AllSetSpecialCase(t *testing.T) {
	// Test the special case where interval is [255, 255]
	intervals := []parser.Interval{{255, 255}}
	ob := New(intervals)

	// This should return the AllSet constant
	if ob != AllSet {
		t.Error("New with interval [255, 255] should return AllSet")
	}

	// Test that all bits are set
	for i := 0; i <= 255; i++ {
		if !ob.Test(byte(i)) {
			t.Errorf("AllSet should have bit %d set", i)
		}
	}
}

func TestNew_EmptyIntervals(t *testing.T) {
	intervals := []parser.Interval{}
	ob := New(intervals)

	// All bits should be unset
	for i := 0; i <= 255; i++ {
		if ob.Test(byte(i)) {
			t.Errorf("Empty intervals should not have bit %d set", i)
		}
	}
}

func TestNew_FullRange(t *testing.T) {
	intervals := []parser.Interval{{0, 255}}
	ob := New(intervals)

	// All bits should be set
	for i := 0; i <= 255; i++ {
		if !ob.Test(byte(i)) {
			t.Errorf("Full range should have bit %d set", i)
		}
	}
}

func TestSet(t *testing.T) {
	ob := &OctetBits{}

	// Test setting various bits
	testBits := []byte{0, 1, 7, 8, 15, 16, 31, 63, 127, 255}

	for _, bit := range testBits {
		ob.set(bit)
		if !ob.Test(bit) {
			t.Errorf("After setting bit %d, Test(%d) should return true", bit, bit)
		}
	}
}

func TestSet_BitManipulation(t *testing.T) {
	ob := &OctetBits{}

	// Test setting bits in the same byte
	ob.set(0) // First bit of first byte
	ob.set(1) // Second bit of first byte
	ob.set(7) // Last bit of first byte

	// Check that the correct bits are set
	if ob[0] != 0x83 { // 10000011 in binary
		t.Errorf("Expected first byte to be 0x83, got 0x%02x", ob[0])
	}

	// Test setting bits in different bytes
	ob.set(8)  // First bit of second byte
	ob.set(16) // First bit of third byte

	if ob[1] != 0x01 { // 00000001 in binary
		t.Errorf("Expected second byte to be 0x01, got 0x%02x", ob[1])
	}

	if ob[2] != 0x01 { // 00000001 in binary
		t.Errorf("Expected third byte to be 0x01, got 0x%02x", ob[2])
	}
}

func TestTest_AllBits(t *testing.T) {
	// Test with AllSet
	for i := 0; i <= 255; i++ {
		if !AllSet.Test(byte(i)) {
			t.Errorf("AllSet should have bit %d set", i)
		}
	}

	// Test with empty OctetBits
	empty := &OctetBits{}
	for i := 0; i <= 255; i++ {
		if empty.Test(byte(i)) {
			t.Errorf("Empty OctetBits should not have bit %d set", i)
		}
	}
}

func TestTest_EdgeCases(t *testing.T) {
	ob := &OctetBits{}

	// Test edge cases: first and last bits
	ob.set(0)
	ob.set(255)

	if !ob.Test(0) {
		t.Error("Test(0) should return true after setting bit 0")
	}

	if !ob.Test(255) {
		t.Error("Test(255) should return true after setting bit 255")
	}

	// Test that other bits are not set
	if ob.Test(1) {
		t.Error("Test(1) should return false when only bits 0 and 255 are set")
	}

	if ob.Test(254) {
		t.Error("Test(254) should return false when only bits 0 and 255 are set")
	}
}

func TestOctetBits_BitPatterns(t *testing.T) {
	tests := []struct {
		name     string
		setBits  []byte
		testBits []struct {
			bit      byte
			expected bool
		}
	}{
		{
			name:    "alternating pattern",
			setBits: []byte{0, 2, 4, 6, 8, 10, 12, 14},
			testBits: []struct {
				bit      byte
				expected bool
			}{
				{0, true},
				{1, false},
				{2, true},
				{3, false},
				{4, true},
				{5, false},
			},
		},
		{
			name:    "byte boundaries",
			setBits: []byte{7, 8, 15, 16, 31, 32},
			testBits: []struct {
				bit      byte
				expected bool
			}{
				{7, true},
				{8, true},
				{15, true},
				{16, true},
				{31, true},
				{32, true},
				{6, false},
				{9, false},
				{14, false},
				{17, false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := &OctetBits{}

			// Set the specified bits
			for _, bit := range tt.setBits {
				ob.set(bit)
			}

			// Test the expected results
			for _, testCase := range tt.testBits {
				result := ob.Test(testCase.bit)
				if result != testCase.expected {
					t.Errorf("Test(%d) = %v, expected %v", testCase.bit, result, testCase.expected)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkNew_SingleInterval(b *testing.B) {
	intervals := []parser.Interval{{10, 20}}

	for i := 0; i < b.N; i++ {
		New(intervals)
	}
}

func BenchmarkNew_MultipleIntervals(b *testing.B) {
	intervals := []parser.Interval{{1, 10}, {20, 30}, {40, 50}, {60, 70}}

	for i := 0; i < b.N; i++ {
		New(intervals)
	}
}

func BenchmarkNew_AllSet(b *testing.B) {
	intervals := []parser.Interval{{255, 255}}

	for i := 0; i < b.N; i++ {
		New(intervals)
	}
}

func BenchmarkSet(b *testing.B) {
	ob := &OctetBits{}

	for i := 0; i < b.N; i++ {
		ob.set(byte(i % 256))
	}
}

func BenchmarkTest(b *testing.B) {
	ob := AllSet

	for i := 0; i < b.N; i++ {
		ob.Test(byte(i % 256))
	}
}

func BenchmarkTest_EmptyBits(b *testing.B) {
	ob := &OctetBits{}

	for i := 0; i < b.N; i++ {
		ob.Test(byte(i % 256))
	}
} 