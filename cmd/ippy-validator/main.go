package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azraelsec/ippy/pkg/ipexpr"
)

func main() {
	matching := flag.String("pattern", "", "IPv4 pattern to validate the ip against")
	ip := flag.String("ip", "", "IPv4 value to validate")
	flag.Parse()

	if *matching == "" {
		fmt.Fprintln(os.Stderr, "error: -pattern flag is required")
		flag.Usage()
		os.Exit(2)
	}
	if *ip == "" {
		fmt.Fprintln(os.Stderr, "error: -ip flag is required")
		flag.Usage()
		os.Exit(2)
	}

	ipexpr, err := ipexpr.Parse(*matching)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot compile: %s\n", err.Error())
		os.Exit(1)
	}

	matches, err := ipexpr.Matches(*ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot match: %s\n", err.Error())
		os.Exit(1)
	}

	if matches {
		fmt.Println("ip matches the given pattern")
	} else {
		fmt.Println("ip does not match the given pattern")
	}
}
