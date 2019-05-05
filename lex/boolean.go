package lex

import (
	"bufio"
	"fmt"

	"github.com/vyevs/json/token"
)

func readBooleanToken(r *bufio.Reader) token.Token {
	literal, ok := readBooleanLiteral(r)
	tokenType := token.Boolean
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

func readBooleanLiteral(r *bufio.Reader) (string, bool) {
	str, err := readNByteString(r, 4)
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
