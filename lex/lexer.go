package lex

import (
	"bufio"
	"io"
	"strings"

	"github.com/vyevs/gojson/tok"
)

// Lexer reads bytes from r into Tokens
type Lexer struct {
	r              *bufio.Reader
	lastTokenBytes []byte
	lastString     strings.Builder
}

// New returns a Lexer that will tokenize the input from r
func New(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r),
		lastTokenBytes: make([]byte, 0, 1024)}
}

// ReadToken reads a single Token from the Lexer
func (l *Lexer) ReadToken() tok.Type {
	if !consumeWhiteSpace(l.r) {
		return tok.EOF
	}
	return l.readTokenNoWhitespace()
}

func (l *Lexer) GetTokenBytes() []byte {
	out := l.lastTokenBytes
	l.lastTokenBytes = l.lastTokenBytes[:0]
	return out
}

func (l *Lexer) GetLastString() string {
	out := l.lastString.String()
	l.lastString.Reset()
	return out
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
func (l *Lexer) readTokenNoWhitespace() tok.Type {
	b, _ := l.r.ReadByte()

	tt := tok.ByteToType(b)

	ok := true
	if tt == tok.True {
		ok = consumeBytes(l.r, []byte{'r', 'u', 'e'})
	} else if tt == tok.False {
		ok = consumeBytes(l.r, []byte{'a', 'l', 's', 'e'})
	} else if tt == tok.Null {
		ok = consumeBytes(l.r, []byte{'u', 'l', 'l'})
	} else if tt == tok.String {
		ok = l.readStringBytes()
	} else if tt == tok.Integer {
		_ = l.r.UnreadByte()
		tt = l.readNumericBytes()
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
