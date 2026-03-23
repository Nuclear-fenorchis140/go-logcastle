package logcastle

import (
	"bufio"
	"io"
)

// Scanner is a high-performance line scanner
type Scanner struct {
	scanner *bufio.Scanner
}

// NewScanner creates a scanner with optimized buffer size for log processing.
// Supports lines up to 1MB (MaxScanTokenSize).
func NewScanner(r io.Reader) *Scanner {
	scanner := bufio.NewScanner(r)
	// Use larger buffer for performance
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return &Scanner{scanner: scanner}
}

// Scan advances to the next line
func (s *Scanner) Scan() bool {
	return s.scanner.Scan()
}

// Bytes returns the current line as bytes
func (s *Scanner) Bytes() []byte {
	return s.scanner.Bytes()
}

// Err returns any error encountered
func (s *Scanner) Err() error {
	return s.scanner.Err()
}
