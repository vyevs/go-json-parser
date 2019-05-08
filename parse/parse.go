package parse

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/vyevs/gojson/lex"
	"github.com/vyevs/gojson/tok"
)

var errInvalidJSON = errors.New("Invalid JSON")

// Parse reads the bytes in r and returns a JSON doc (if valid)
// or an error on some JSON syntax error
func Parse(r io.Reader) (interface{}, error) {
	l := lex.New(r)

	tt := l.ReadToken()

	switch tt {
	case tok.OpeningCurlyBrace:
		return parseValue(l, tt)
	case tok.Invalid:
		return nil, errInvalidJSON
	default:
		return parseSingleValueDoc(l, tt)
	}
}

func parseSingleValueDoc(l *lex.Lexer, ct tok.Type) (interface{}, error) {
	v, err := parseValue(l, ct)
	if err != nil {
		return nil, err
	}
	eof := l.ReadToken()
	if eof != tok.EOF {
		return nil, errInvalidJSON
	}
	return v, nil
}

// parsing an object consists of reading a key followed by a colon
// followed by a value repeatedly until we encounter a closing curly brace
func parseObject(l *lex.Lexer) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	var seenValue bool
	for t := l.ReadToken(); t != tok.ClosingCurlyBrace; t = l.ReadToken() {
		if seenValue {
			if t != tok.Comma {
				return nil, errInvalidJSON
			}
			t = l.ReadToken()
		}
		if t != tok.String {
			return nil, errInvalidJSON
		}
		key := string(l.GetTokenBytes())

		t = l.ReadToken()
		if t != tok.Colon {
			return nil, errInvalidJSON
		}

		v, err := parseValue(l, l.ReadToken())
		if err != nil {
			return nil, err
		}
		seenValue = true

		if _, ok := out[key]; ok {
			return nil, fmt.Errorf("Found duplicate key %q", key)
		}
		out[key] = v
	}
	return out, nil
}

// parseValue expects ct to contain the first token representing a value
// e.g.: "[" for array, "{" for object, str for string value
// the comma before a value (if any) should already be consumed by the calling func
func parseValue(l *lex.Lexer, ct tok.Type) (interface{}, error) {
	switch ct {
	case tok.String:
		return string(l.GetTokenBytes()), nil
	case tok.Integer:
		return parseInteger(l.GetTokenBytes())
	case tok.Float:
		return parseFloatingPoint(l.GetTokenBytes())
	case tok.OpeningCurlyBrace:
		return parseObject(l)
	case tok.OpeningSquareBracket:
		return parseArray(l)
	case tok.True:
		return true, nil
	case tok.False:
		return false, nil
	case tok.Null:
		return nil, nil
	}
	return nil, fmt.Errorf("parseValue() received unknown token %v", ct)
}

// parseArray starts parsing AFTER the opening square bracket has already been consumed
func parseArray(l *lex.Lexer) ([]interface{}, error) {
	out := make([]interface{}, 0, 8)
	var seenValue bool

	// array tokens are read in pairs after the 1st value
	// 1st token should be a comma followed by a value EXCEPT for the 1st
	// value in the array
	for t := l.ReadToken(); t != tok.ClosingSquareBracket; t = l.ReadToken() {
		if seenValue {
			if t != tok.Comma {
				return nil, errInvalidJSON
			}
			t = l.ReadToken()
		}
		seenValue = true
		v, err := parseValue(l, t)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func parseInteger(lit []byte) (int, error) {
	v, err := strconv.Atoi(string(lit))
	if err != nil {
		return 0, fmt.Errorf("Invalid value found: %q", lit)
	}
	return v, nil
}

func parseFloatingPoint(lit []byte) (float64, error) {
	v, err := strconv.ParseFloat(string(lit), 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid value found: %q", lit)
	}
	return v, nil
}
