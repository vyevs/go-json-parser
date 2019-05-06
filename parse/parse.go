package parse

import (
	"fmt"
	"io"
	"strconv"

	"github.com/vyevs/gojson/lex"
	"github.com/vyevs/gojson/tok"
)

// Parse reads the bytes in r and returns a JSON doc (if valid)
// or an error on some JSON syntax error
func Parse(r io.Reader) (interface{}, error) {
	l := lex.New(r)

	t := l.ReadToken()

	switch t.TokenType {
	case tok.OpeningCurlyBrace:
		return parseValue(l, t)
	case tok.Invalid:
		return nil, fmt.Errorf("Found invalid token: %s", t.Literal)
	default:
		return parseSingleValueDoc(l, t)
	}
}

func parseSingleValueDoc(l lex.Lexer, ct tok.Token) (interface{}, error) {
	v, err := parseValue(l, ct)
	if err != nil {
		return nil, err
	}
	eof := l.ReadToken()
	if eof.TokenType != tok.EOF {
		return nil, fmt.Errorf("Expected end of document, found: %q", eof.Literal)
	}
	return v, nil
}

// parsing an object consists of reading a key followed by a colon
// followed by a value repeatedly until we encounter a closing curly brace
func parseObject(l lex.Lexer) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	var seenValue bool
	for t := l.ReadToken(); t.TokenType != tok.ClosingCurlyBrace; t = l.ReadToken() {
		if seenValue {
			if t.TokenType != tok.Comma {
				return nil, fmt.Errorf("Expected comma(%q) got %q", ",", t.Literal)
			}
			t = l.ReadToken()
		}
		if t.TokenType != tok.String {
			return nil, fmt.Errorf("Expected key, got: %q", t.Literal)
		}
		key := t.Literal

		t = l.ReadToken()
		if t.TokenType != tok.Colon {
			return nil, fmt.Errorf("Expected colon (%q): got: %q", ":", t.Literal)
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
func parseValue(l lex.Lexer, ct tok.Token) (interface{}, error) {
	switch ct.TokenType {
	case tok.String:
		return ct.Literal, nil
	case tok.Integer:
		return parseInteger(ct.Literal)
	case tok.FloatingPoint:
		return parseFloatingPoint(ct.Literal)
	case tok.OpeningCurlyBrace:
		return parseObject(l)
	case tok.OpeningSquareBracket:
		return parseArray(l)
	case tok.Boolean:
		return parseBool(ct.Literal)
	case tok.Null:
		return nil, nil
	}
	return nil, fmt.Errorf("parseValue() received unknown token %v", ct)
}

// parseArray starts parsing AFTER the opening square bracket has already been consumed
func parseArray(l lex.Lexer) ([]interface{}, error) {
	out := make([]interface{}, 0)
	var seenValue bool

	// array tokens are read in pairs after the 1st value
	// 1st token should be a comma followed by a value EXCEPT for the 1st
	// value in the array
	for t := l.ReadToken(); t.TokenType != tok.ClosingSquareBracket; t = l.ReadToken() {
		if seenValue {
			if t.TokenType != tok.Comma {
				return nil, fmt.Errorf("expected comma(%q), found: %q", ",", t.Literal)
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

func parseInteger(lit string) (int, error) {
	v, err := strconv.Atoi(lit)
	if err != nil {
		return 0, fmt.Errorf("Invalid value found: %q", lit)
	}
	return v, nil
}

func parseFloatingPoint(lit string) (float64, error) {
	v, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid value found: %q", lit)
	}
	return v, nil
}

func parseBool(lit string) (bool, error) {
	v, err := strconv.ParseBool(lit)
	if err != nil {
		return false, fmt.Errorf("Invalid value found: %q", lit)
	}
	return v, nil
}
