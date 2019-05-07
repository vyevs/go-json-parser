package lex

import (
	"bufio"
	"io"

	"github.com/vyevs/gojson/tok"
)

// Lexer reads bytes from r into Tokens
type Lexer struct {
	r              *bufio.Reader
	lastTokenBytes []byte
}

// New returns a Lexer that will tokenize the input from r
func New(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

// ReadToken reads a single Token from the Lexer
func (l Lexer) ReadToken() tok.TokenType {
	if !consumeWhiteSpace(l.r) {
		return tok.EOFToken
	}
	return l.readTokenNoWhitespace()
}

func (l Lexer) GetTokenLiteral() []byte {
	return l.lastTokenBytes
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
		if !(b == ' ' || b == '\n' || b == '\t') {
			_ = r.UnreadByte()
			return true
		}
	}
}

// reads a single token that begins with the next byte read from r
func (l Lexer) readTokenNoWhitespace() tok.TokenType {
	b, _ := l.r.ReadByte()

	tt := tok.ByteToTokenType(b)

	var ok bool
	if tt == tok.BooleanTrue {
		ok = consumeBytes(l.r, []byte{'r', 'u', 'e'})
	} else if tt == tok.BooleanFalse {
		ok = consumeBytes(l.r, []byte{'a', 'l', 's', 'e'})
	} else if tt == tok.Null {
		ok = consumeBytes(l.r, []byte{'u', 'l', 'l'})
	} else if tt == tok.String {
		l.lastTokenBytes, ok = readStringBytes(l.r)
	} else if tt == tok.Integer {
		_ = l.r.UnreadByte()
		l.lastTokenBytes, tt = readNumericBytes(l.r)
	}

	if !ok {
		return tok.Invalid
	}

	return tt
}

func consumeBytes(r *bufio.Reader, target []byte) bool {
	for _, b := range target {
		curB, err := r.ReadByte()
		if err != nil || curB != b {
			return false
		}
	}
	return true
}

// utility function used to read bool & null literals
func readNByteStr(r *bufio.Reader, n int) (string, error) {
	bytes := make([]byte, n)
	n, err := io.ReadFull(r, bytes)
	bytes = bytes[:n]
	return string(bytes), err
}
