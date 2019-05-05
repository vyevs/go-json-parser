package parse

import (
	"fmt"
	"io"
	"strconv"

	"github.com/vyevs/json/lex"
	"github.com/vyevs/json/token"
)

func Parse(r io.Reader) (interface{}, error) {
	l := lex.New(r)

	tok := l.ReadToken()

	switch tok.TokenType {
	case token.OpeningCurlyBrace:
		return parseValue(l, tok)
	case token.Invalid:
		return nil, fmt.Errorf("Found invalid token: %s", tok.Literal)
	default:
		return parseSingleValueDoc(l, tok)
	}
}

func parseSingleValueDoc(l lex.Lexer, ct token.Token) (interface{}, error) {
	v, err := parseValue(l, ct)
	if err != nil {
		return nil, err
	}
	eof := l.ReadToken()
	if eof.TokenType != token.EOF {
		return nil, fmt.Errorf("Expected end of document, found: %q", eof.Literal)
	}
	return v, nil
}

// parsing an object consists of reading a key followed by a colon
// followed by a value
func parseObject(l lex.Lexer) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	var seenValue bool
	for tok := l.ReadToken(); tok.TokenType != token.ClosingCurlyBrace; tok = l.ReadToken() {
		fmt.Println(tok)
		if seenValue {
			if tok.TokenType != token.Comma {
				return nil, fmt.Errorf("Expected comma(%q) got %q", ",", tok.Literal)
			}
			tok = l.ReadToken()
		}
		if tok.TokenType != token.String {
			return nil, fmt.Errorf("Expected key, got: %q", tok.Literal)
		}
		key := tok.Literal

		tok = l.ReadToken()
		if tok.TokenType != token.Colon {
			return nil, fmt.Errorf("Expected colon (%q): got: %q", ":", tok.Literal)
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
func parseValue(l lex.Lexer, ct token.Token) (interface{}, error) {
	switch ct.TokenType {
	case token.String:
		return ct.Literal, nil
	case token.Integer:
		return parseInteger(ct.Literal)
	case token.FloatingPoint:
		return parseFloatingPoint(ct.Literal)
	case token.OpeningCurlyBrace:
		return parseObject(l)
	case token.OpeningSquareBracket:
		return parseArray(l)
	case token.Boolean:
		return parseBool(ct.Literal)
	case token.Null:
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
	for tok := l.ReadToken(); tok.TokenType != token.ClosingSquareBracket; tok = l.ReadToken() {
		if seenValue {
			if tok.TokenType != token.Comma {
				return out, fmt.Errorf("expected comma(%q), found: %q", ",", tok.Literal)
			}
			tok = l.ReadToken()
		}
		seenValue = true
		v, err := parseValue(l, tok)
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
