package lex

import (
	"bufio"

	"github.com/vyevs/json/token"
)

// attempts to read "null" and return the corresponding token
// returns an Token with TokenType Invalid if not successful
// upon success consumes ONLY the "null" string from the reader
func readNullToken(r *bufio.Reader) token.Token {
	literal, ok := readNullLiteral(r)
	tokenType := token.Null
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

// attemps to read "null" from the reader
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
