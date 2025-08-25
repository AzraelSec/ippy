# ippy

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**ippy** is a fast and flexible IPv4 pattern matching library for Go. It allows you to define complex IP address patterns using a simple expression syntax and efficiently match IPv4 addresses against those patterns.

## Features

- **Flexible Pattern Syntax**: Support for ranges, wildcards, and comma-separated values in IP expressions
- **High Performance**: Uses bit vectors for efficient pattern matching with O(1) lookup time
- **Simple API**: Easy-to-use interface with parse-once, match-many semantics
- **IP Generation**: Generate all IPs that match a given pattern with iterator support
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Comprehensive Testing**: Well-tested with extensive unit tests
- **Command Line Tool**: Includes a validator tool for testing patterns

## Installation

```bash
go get github.com/azraelsec/ippy
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/azraelsec/ippy/pkg/ipexpr"
)

func main() {
    // Parse a pattern
    expr, err := ipexpr.Parse("192.168.1.*")
    if err != nil {
        log.Fatal(err)
    }

    // Test IP addresses
    ips := []string{"192.168.1.1", "192.168.1.255", "10.0.0.1"}

    for _, ip := range ips {
        matches, err := expr.Matches(ip)
        if err != nil {
            log.Printf("Error matching %s: %v", ip, err)
            continue
        }

        if matches {
            fmt.Printf("✓ %s matches the pattern\n", ip)
        } else {
            fmt.Printf("✗ %s does not match the pattern\n", ip)
        }
    }
}
```

## Pattern Syntax

The pattern syntax supports flexible expressions for each octet of an IPv4 address:

### Basic Patterns

| Pattern | Description               | Example Matches                                                       |
| ------- | ------------------------- | --------------------------------------------------------------------- |
| `*`     | Matches any value (0-255) | `192.168.1.*` matches `192.168.1.1`, `192.168.1.255`                  |
| `123`   | Matches exact value       | `192.168.1.123` matches only `192.168.1.123`                          |
| `1-10`  | Matches range (inclusive) | `192.168.1.1-10` matches `192.168.1.1` through `192.168.1.10`         |
| `1,3,5` | Matches multiple values   | `192.168.1.1,3,5` matches `192.168.1.1`, `192.168.1.3`, `192.168.1.5` |

### Complex Patterns

You can combine different pattern types within a single octet:

```go
// Mixed patterns
"192.168.1.1,10-20,100,200-255"  // Matches 1, 10-20, 100, and 200-255
"10.0-50.*.1-254"                // Matches 10.0-50.anything.1-254
"172.16,20.0.1"                  // Matches 172.16.0.1 and 172.20.0.1
```

## API Reference

### Functions

#### `Parse(expr string) (*IPExpr, error)`

Parses an IP pattern expression and returns an IPExpr instance.

```go
expr, err := ipexpr.Parse("192.168.1.*")
if err != nil {
    // Handle parsing error
}
```

**Parameters:**

- `expr string`: IPv4 pattern expression (e.g., "192.168.1.\*")

**Returns:**

- `*IPExpr`: Parsed expression ready for matching
- `error`: Parsing error if the pattern is invalid

#### `(ie IPExpr) Matches(ip string) (bool, error)`

Tests whether an IP address matches the parsed pattern.

```go
matches, err := expr.Matches("192.168.1.100")
if err != nil {
    // Handle matching error (invalid IP format, etc.)
}
```

**Parameters:**

- `ip string`: IPv4 address to test (e.g., "192.168.1.1")

**Returns:**

- `bool`: Whether the IP matches the pattern
- `error`: Matching error if IP format is invalid

#### `(ie IPExpr) Generate() iter.Seq2[int, ip.IPv4]`

Generates all IP addresses that match the pattern using Go's iterator interface.

```go
expr, _ := ipexpr.Parse("192.168.1.1-3")
for i, ip := range expr.Generate() {
    fmt.Printf("%d: %s\n", i, ip)
}
// Output:
// 0: 192.168.1.1
// 1: 192.168.1.2
// 2: 192.168.1.3
```

**Returns:**

- `iter.Seq2[int, ip.IPv4]`: Iterator yielding index and IP address pairs

## Command Line Tool

The library includes a command-line validator tool:

```bash
# Build the tool
go build -o ippy-validator ./cmd/ippy-validator

# Test if an IP matches a pattern
./ippy-validator -pattern "192.168.1.*" -ip "192.168.1.100"
# Output: ip matches the given pattern

./ippy-validator -pattern "192.168.1.*" -ip "10.0.0.1"
# Output: ip does not match the given pattern
```

### Installation via go install

```bash
go install github.com/azraelsec/ippy/cmd/ippy-validator@latest
```

## Examples

### Corporate Network Matching

```go
// Match corporate network ranges
patterns := []string{
    "10.*.*.*",           // Private Class A
    "172.16-31.*.*",      // Private Class B
    "192.168.*.*",        // Private Class C
    "203.0.113.1-100",    // Specific public range
}

for _, pattern := range patterns {
    expr, err := ipexpr.Parse(pattern)
    if err != nil {
        log.Printf("Invalid pattern %s: %v", pattern, err)
        continue
    }

    // Test against your IPs
    testIP := "172.20.1.50"
    if matches, _ := expr.Matches(testIP); matches {
        fmt.Printf("%s matches corporate network pattern: %s\n", testIP, pattern)
        break
    }
}
```

### Load Balancer Health Check

```go
// Define health check IP ranges
healthCheckPattern := "10.0.1.1-10,20-30,100"

expr, err := ipexpr.Parse(healthCheckPattern)
if err != nil {
    log.Fatal(err)
}

func isHealthCheckIP(clientIP string) bool {
    matches, err := expr.Matches(clientIP)
    if err != nil {
        log.Printf("Error checking health check IP %s: %v", clientIP, err)
        return false
    }
    return matches
}
```

### Security Filtering

```go
// Block suspicious IP ranges
suspiciousPatterns := []string{
    "0.0.0.0",           // Invalid source
    "127.*.*.*",         // Loopback
    "169.254.*.*",       // Link-local
    "224-255.*.*.*",     // Multicast and reserved
}

func isIPSuspicious(ip string) bool {
    for _, pattern := range suspiciousPatterns {
        expr, err := ipexpr.Parse(pattern)
        if err != nil {
            continue
        }

        if matches, _ := expr.Matches(ip); matches {
            return true
        }
    }
    return false
}
```

### IP Range Generation

```go
// Generate all IPs in a pattern
expr, err := ipexpr.Parse("192.168.1.1-5,10")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Generated IPs:")
for i, ip := range expr.Generate() {
    fmt.Printf("%d: %s\n", i, ip)
}
// Output:
// 0: 192.168.1.1
// 1: 192.168.1.2
// 2: 192.168.1.3
// 3: 192.168.1.4
// 4: 192.168.1.5
// 5: 192.168.1.10
```

## Performance

The library uses bit vectors for efficient pattern matching, providing:

- **O(1)** time complexity for matching operations
- **O(n)** preprocessing time for pattern compilation where n is the number of intervals
- **Constant memory usage** per octet (256 bits = 32 bytes)
- **Efficient generation** with iterator-based IP enumeration

## Architecture

The library consists of several internal components:

- **Lexer**: Tokenizes IP pattern expressions into tokens (numbers, ranges, wildcards, commas)
- **Parser**: Parses tokens into interval structures representing valid ranges
- **Bit Vector**: Uses 256-bit vectors (32 bytes) per octet for O(1) membership testing
- **IP Parser**: Validates and parses IPv4 addresses into octets
- **IPExpr**: High-level API that orchestrates the components and provides matching/generation

### Key Design Decisions

- **Bit Vectors over BST**: Replaced binary search trees with bit vectors for constant-time lookups
- **Iterator-based Generation**: Uses Go 1.23+ iterators for memory-efficient IP generation
- **Parse-once Semantics**: Patterns are parsed once and can be reused for multiple matches
- **Zero Dependencies**: Pure Go implementation using only standard library

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

_Built with ❤️ in Go_
