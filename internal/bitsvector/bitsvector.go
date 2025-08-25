// Package bitsvector provides a compact bit vector implementation for representing
// sets of byte values (0-255) using a 32-byte array where each bit corresponds
// to whether a specific byte value is present in the set.
package bitsvector

import (
	"github.com/azraelsec/ippy/internal/parser"
)

var AllSet = OctetBits{
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF,
}

type OctetBits [32]byte

func New(its []parser.Interval) OctetBits {
	if len(its) == 1 && its[0][0] == its[0][1] && its[0][0] == 255 {
		asc := AllSet
		return asc
	}

	ob := &OctetBits{}
	for _, it := range its {
		start, end := int(it[0]), int(it[1])
		for i := start; i <= end; i++ {
			ob.set(byte(i))
		}
	}
	return *ob
}

func (o *OctetBits) set(n byte) {
	o[n/8] |= 1 << (n % 8)
}

func (o OctetBits) Test(n byte) bool {
	return o[int(n)/8]&(1<<(n%8)) != 0
}
