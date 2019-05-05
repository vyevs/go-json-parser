package lex

import (
	"bufio"
	"fmt"

	"github.com/vyevs/json/token"
)

// attempts to read "true" or "false" and returns the corresponding token
// if not successful in reading the bool literal, returns an invalid token
// with the literal encountered instead
func readBoolToken(r *bufio.Reader) token.Token {
	literal, ok := readBoolLiteral(r)
	tokenType := token.Boolean
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

// attempts to read either "true" or "false" from the reader
// bool return value indicates whether this was successful
// upon success consumes only the bytes of "true"/"false"
func readBoolLiteral(r *bufio.Reader) (string, bool) {
	str, err := readNByteStr(r, 4)
	if err != nil {
		return str, false
	}
	if str == "true" {
		return "true", true
	}

	b, err := r.ReadByte()
	if err != nil {
		return str, false
	}
	str = fmt.Sprintf("%s%c", str, b)
	if str == "false" {
		return "false", true
	}

	return str, false
}
