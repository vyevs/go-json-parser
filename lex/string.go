package lex

import (
	"bufio"
	"strings"
	"vlad/json/token"
)

func readStringToken(r *bufio.Reader) token.Token {
	literal, ok := readStringLiteral(r)
	tokenType := token.String
	if !ok {
		tokenType = token.Invalid
	}
	return token.Token{TokenType: tokenType, Literal: literal}
}

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
