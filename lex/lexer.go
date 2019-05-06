package lex

import (
	"bufio"
	"io"

	"github.com/vyevs/json/tok"
)

// Lexer reads bytes from r into Tokens
type Lexer struct {
	r *bufio.Reader
}

// New returns a Lexer that will tokenize the input from r
func New(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

// ReadToken reads a single Token from the Lexer
func (l Lexer) ReadToken() tok.Token {
	if !consumeWhiteSpace(l.r) {
		return tok.EOFToken
	}
	return readTokenNoWhitespace(l.r)
}

// consumes all whitespace characters as defined by isWhiteSpace()
// returns whether there are any more characters to be read
// from the reader
func consumeWhiteSpace(r *bufio.Reader) bool {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return false
		}
		if !isWhitespace(b) {
			_ = r.UnreadByte()
			return true
		}
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t'
}

// reads a single token that begins with the next byte read from r
func readTokenNoWhitespace(r *bufio.Reader) tok.Token {
	b, _ := r.ReadByte()

	return readTokenBeginningWithByte(r, b)
}

func readTokenBeginningWithByte(r *bufio.Reader, b byte) tok.Token {
	tt := tok.ByteToTokenType(b)

	if tt == tok.Invalid {
		return tok.Token{TokenType: tok.Invalid, Literal: string(b)}
	}

	return readTokenOfType(r, tt)
}

func readTokenOfType(r *bufio.Reader, tt tok.TokenType) tok.Token {
	if tok, ok := tok.TokenTypeToPredefinedToken(tt); ok {
		return tok
	}

	switch tt {

	case tok.String:
		return readStringToken(r)

	case tok.Null:
		_ = r.UnreadByte()
		return readNullToken(r)

	case tok.Boolean:
		_ = r.UnreadByte()
		return readBoolToken(r)

	case tok.Integer:
		_ = r.UnreadByte()
		return readNumericToken(r)

	default:
		// not possible, indicates some internal bug
		panic("readTokenOfType() received unknown tokenType, INTERNAL BUG")
	}
}

// utility function used to read bool & null literals
func readNByteStr(r *bufio.Reader, n int) (string, error) {
	bytes := make([]byte, n)
	n, err := io.ReadFull(r, bytes)
	bytes = bytes[:n]
	return string(bytes), err
}
