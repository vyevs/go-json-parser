package lex

import (
	"bufio"
)

// attempts to read a string token (literal contained in double quotes)
// expects the beginning double quote to have been consumed already
// consumes all bytes up to and including the terminating double quote
// TODO: DOES NOT CURRENTLY SUPPORT ESCAPE SEQUENCES
func readStringBytes(r *bufio.Reader) ([]byte, bool) {
	buf := make([]byte, 0, 32)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return buf, false
		}
		if b == '"' {
			return buf, true
		}
		buf = append(buf, b)
	}
}
