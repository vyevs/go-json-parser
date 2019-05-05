package lex

import (
	"bufio"
	"strings"
	"vlad/json/token"
)

func readNumericToken(r *bufio.Reader) token.Token {
	literal, ok := readNumericLiteral(r)
	if !ok {
		return token.Token{TokenType: token.Invalid, Literal: literal}
	}

	if ok := validateNumericLiteral(literal); !ok {
		return token.Token{TokenType: token.Invalid, Literal: literal}
	}

	tokType := numericLiteralTokenType(literal)

	return token.Token{TokenType: tokType, Literal: literal}
}

func readNumericLiteral(r *bufio.Reader) (string, bool) {
	b, err := r.ReadByte()
	if err != nil {
		return "", false
	} else if !isDigit(b) && b != '-' {
		return string(b), false
	}

	var seenPeriod bool
	var builder strings.Builder
	builder.WriteByte(b)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return builder.String(), true
		}
		if b == '.' {
			if seenPeriod {
				builder.WriteByte(b)
				return builder.String(), false
			}
			seenPeriod = true
		} else if !isDigit(b) {
			break
		}
		builder.WriteByte(b)
	}
	_ = r.UnreadByte()

	return builder.String(), true
}

func numericLiteralTokenType(literal string) token.TokenType {
	if strings.Contains(literal, ".") {
		return token.FloatingPoint
	}
	return token.Integer
}

// validates that the literal does not begin with an illegal 0
// e.g.: 01, 01.1, -01, -01.1
func validateNumericLiteral(literal string) bool {
	if literal[0] == '-' {
		literal = literal[1:]
	}

	if len(literal) > 1 && literal[0] == '0' && literal[1] != '.' {
		return false
	}
	if literal[len(literal)-1] == '.' {
		return false
	}
	return true
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
