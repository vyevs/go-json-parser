package lex

import (
	"bufio"
	"strings"

	"github.com/vyevs/gojson/tok"
)

// reads an Integer or FloatingPoint token from r
// returns an Invalid token if the following bytes do not form a numeric token
func readNumericToken(r *bufio.Reader) tok.Token {
	literal, ok := readNumericLiteral(r)
	if !ok {
		return tok.Token{TokenType: tok.Invalid, Literal: literal}
	}

	if ok := validateNumericLiteral(literal); !ok {
		return tok.Token{TokenType: tok.Invalid, Literal: literal}
	}

	tokType := numericLiteralTokenType(literal)

	return tok.Token{TokenType: tokType, Literal: literal}
}

// readNumericLiteral attempts to read a numeric literal(either integer or floating point) from r
// consumes only the bytes of the numeric literal, not the byte after
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

// checks the numeric type of the literal, either token.Integer or token.FloatingPoint
func numericLiteralTokenType(literal string) tok.TokenType {
	if strings.Contains(literal, ".") {
		return tok.FloatingPoint
	}
	return tok.Integer
}

// validates that the literal does not begin with an illegal 0
// e.g.: 01, 01.1, -01, -01.1
func validateNumericLiteral(literal string) bool {
	// strip minus for negative number
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
