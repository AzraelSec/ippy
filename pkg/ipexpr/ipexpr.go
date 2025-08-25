// Package ipexpr implements a flexible IPv4 pattern matching and generation system.
//
// This package allows you to define complex IPv4 address patterns using a simple
// expression syntax and efficiently match IP addresses against those patterns.
// It supports ranges (1-10), wildcards (*), comma-separated values (1,3,5),
// and combinations thereof in each octet of an IPv4 address.
package ipexpr

import (
	"fmt"
	"iter"
	"net"
	"strings"

	"github.com/azraelsec/ippy/internal/bitsvector"
	"github.com/azraelsec/ippy/internal/ip"
	"github.com/azraelsec/ippy/internal/parser"
)

type IPExpr struct {
	octets [4]bitsvector.OctetBits
}

func (ie IPExpr) Matches(i string) (bool, error) {
	ip, err := ip.Parse(i)
	if err != nil {
		return false, err
	}

	for i, octet := range ip {
		if !ie.octets[i].Test(octet) {
			return false, nil
		}
	}
	return true, nil
}

func (ie IPExpr) Generate() iter.Seq2[int, ip.IPv4] {
	i := 0
	counter := [4]int{}
	return func(yield func(int, ip.IPv4) bool) {
		for counter[0] != 0 && counter[1] != 0 && counter[2] != 0 && counter[3] != 0 {
			ip := net.IPv4(byte(counter[0]), byte(counter[1]), byte(counter[2]), byte(counter[3]))

			for i := range 4 {
				carry := false
				for {
					counter[3-i] = (counter[3-i] + 1) % 256
					if counter[3-i] == 0 {
						carry = true
					}
					if ie.octets[3-i].Test(byte(counter[3-i])) {
						break
					}
				}
				if !carry {
					break
				}
			}

			if !yield(i, ip) {
				return
			}

			i++
		}
	}
}

func Parse(expr string) (*IPExpr, error) {
	parts := strings.Split(expr, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid ip expression: %s", expr)
	}

	ip := &IPExpr{}
	for i, part := range parts {
		bv, err := parseOctet(part)
		if err != nil {
			return nil, err
		}
		ip.octets[i] = bv
	}
	return ip, nil
}

func parseOctet(o string) (bitsvector.OctetBits, error) {
	// NOTE: shall we return parsing errors instead of a generic message?
	its, ok := parser.New(o).Parse()
	if !ok {
		return bitsvector.OctetBits{}, fmt.Errorf("invalid octet format in %s", o)
	}
	return bitsvector.New(its), nil
}
