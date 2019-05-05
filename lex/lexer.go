package lex

import (
	"bufio"
	"io"

	"github.com/vyevs/json/token"
)

type Lexer struct {
	r *bufio.Reader
}

func New(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

func (l Lexer) ReadToken() token.Token {
	if !consumeWhiteSpace(l.r) {
		return token.EOFToken
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

func readTokenNoWhitespace(r *bufio.Reader) token.Token {
	b, _ := r.ReadByte()

	return readTokenBeginningWithByte(r, b)
}

func readTokenBeginningWithByte(r *bufio.Reader, b byte) token.Token {
	tt := token.ByteToTokenType(b)

	if tt == token.Invalid {
		return token.Token{TokenType: token.Invalid, Literal: string(b)}
	}

	return readTokenOfType(r, tt)
}

func readTokenOfType(r *bufio.Reader, tt token.TokenType) token.Token {
	if tok, ok := token.TokenTypeToPredefinedToken(tt); ok {
		return tok
	}

	switch tt {

	case token.String:
		return readStringToken(r)

	case token.Null:
		_ = r.UnreadByte()
		return readNullToken(r)

	case token.Boolean:
		_ = r.UnreadByte()
		return readBooleanToken(r)

	case token.Integer:
		_ = r.UnreadByte()
		return readNumericToken(r)
	}

	// not possible, but anyways
	panic("readTokenOfType() received unknown tokenType")
}

func readNByteString(r *bufio.Reader, n int) (string, error) {
	bytes := make([]byte, n)
	n, err := io.ReadFull(r, bytes)
	bytes = bytes[:n]
	return string(bytes), err
}
