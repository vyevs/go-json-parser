package lex

import (
	"bufio"
	"strings"

	"github.com/vyevs/json/token"
)

// attempts to read a string literal token from r
// expects the beginning double quote to have already been consumed
func readStringToken(r *bufio.Reader) token.Token {
	literal, ok := readStringLiteral(r)
	tokenType := token.String
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

// attempts to read a string token (literal contained in double quotes)
// expects the beginning double quote to have been consumed already
// consumes all bytes up to and including the terminating double quote
// TODO: DOES NOT CURRENTLY SUPPORT ESCAPE SEQUENCES
func readStringLiteral(r *bufio.Reader) (string, bool) {
	var builder strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			return builder.String(), false
		}
		if b == '"' {
			return builder.String(), true
		}
		builder.WriteByte(b)
	}
}
