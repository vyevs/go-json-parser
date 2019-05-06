package lex

import (
	"bufio"

	"github.com/vyevs/json/tok"
)

// attempts to read "null" and return the corresponding token
// returns an Token with TokenType Invalid if not successful
// upon success consumes ONLY the "null" string from the reader
func readNullToken(r *bufio.Reader) tok.Token {
	literal, ok := readNullLiteral(r)
	tokenType := tok.Null
	if !ok {
		tokenType = tok.Invalid
	}
	return tok.Token{TokenType: tokenType, Literal: literal}
}

// attempts to read "null" from the reader
// bool return value indicates whether this was successful
func readNullLiteral(r *bufio.Reader) (string, bool) {
	str, err := readNByteStr(r, 4)
	if err != nil {
		return str, false
	} else if str != "null" {
		return str, false
	}
	return str, true
}
