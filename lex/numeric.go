package lex

import (
	"github.com/vyevs/gojson/tok"
)

// reads an Integer or FloatingPoint token from r
// returns an Invalid token if the following bytes do not form a numeric token
func (l *Lexer) readNumericBytes() tok.Type {
	if !l.readNumericLiteral() {
		return tok.Invalid
	}

	if ok := validateNumericBytes(l.lastTokenBytes); !ok {
		return tok.Invalid
	}

	tokType := numericLiteralType(l.lastTokenBytes)

	return tokType
}

// readNumericLiteral attempts to read a numeric literal(either integer or floating point) from r
// consumes only the bytes of the numeric literal, not the byte after
func (l *Lexer) readNumericLiteral() bool {
	b, err := l.r.ReadByte()
	if err != nil {
		return false
	} else if !isDigit(b) && b != '-' {
		return false
	}

	var seenPeriod bool
	l.lastTokenBytes = append(l.lastTokenBytes, b)
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			return true
		}
		if b == '.' {
			if seenPeriod {
				return false
			}
			seenPeriod = true
		} else if !isDigit(b) {
			break
		}
		l.lastTokenBytes = append(l.lastTokenBytes, b)
	}
	_ = l.r.UnreadByte()

	return true
}

// checks the numeric type of the literal, either token.Integer or token.FloatingPoint
func numericLiteralType(literal []byte) tok.Type {
	for _, b := range literal {
		if b == '.' {
			return tok.Float
		}
	}
	return tok.Integer
}

// validates that the literal does not begin with an illegal 0
// e.g.: 01, 01.1, -01, -01.1
func validateNumericBytes(literal []byte) bool {
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
