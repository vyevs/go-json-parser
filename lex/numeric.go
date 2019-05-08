package lex

import (
	"bufio"

	"github.com/vyevs/gojson/tok"
)

// reads an Integer or FloatingPoint token from r
// returns an Invalid token if the following bytes do not form a numeric token
func readNumericBytes(r *bufio.Reader) ([]byte, tok.Type) {
	bytes, ok := readNumericLiteral(r)
	if !ok {
		return nil, tok.Invalid
	}

	if ok := validateNumericBytes(bytes); !ok {
		return nil, tok.Invalid
	}

	tokType := numericLiteralType(bytes)

	return bytes, tokType
}

// readNumericLiteral attempts to read a numeric literal(either integer or floating point) from r
// consumes only the bytes of the numeric literal, not the byte after
func readNumericLiteral(r *bufio.Reader) ([]byte, bool) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, false
	} else if !isDigit(b) && b != '-' {
		return nil, false
	}

	var seenPeriod bool
	buf := make([]byte, 0, 8)
	buf = append(buf, b)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return buf, true
		}
		if b == '.' {
			if seenPeriod {
				return nil, false
			}
			seenPeriod = true
		} else if !isDigit(b) {
			break
		}
		buf = append(buf, b)
	}
	_ = r.UnreadByte()

	return buf, true
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
