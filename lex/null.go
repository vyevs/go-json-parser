package lex

import (
	"bufio"
	"vlad/json/token"
)

func readNullToken(r *bufio.Reader) token.Token {
	literal, ok := readNullLiteral(r)
	tokenType := token.Null
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

func readNullLiteral(r *bufio.Reader) (string, bool) {
	str, err := readNByteString(r, 4)
	if err != nil {
		return str, false
	} else if str != "null" {
		return str, false
	}
	return str, true
}
