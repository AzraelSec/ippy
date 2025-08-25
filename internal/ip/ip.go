// Package ip provides basic IPv4 address parsing functionality.
package ip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type IPv4 = net.IP

func Parse(ip string) (IPv4, error) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return net.IPv4zero, fmt.Errorf("invalid ip: %s", ip)
	}

	octets := [4]byte{}
	for i, os := range parts {
		octet, err := strconv.ParseUint(os, 10, 8)
		if err != nil {
			return net.IPv4zero, err
		}
		octets[i] = uint8(octet)
	}

	return net.IPv4(octets[0], octets[1], octets[2], octets[3]).To4(), nil
}
