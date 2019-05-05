package lex

import (
	"bufio"
	"io"
	"vlad/json/token"
)

type Lexer struct {
	r *bufio.Reader
}

func New(r io.Reader) Lexer {
	return Lexer{r: bufio.NewReader(r)}
}

func (lexer Lexer) ReadToken() token.Token {
	if consumeWhiteSpace(lexer.r) != nil {
		return token.EOFToken
	}
	return readTokenNoWhitespace(lexer.r)
}

func consumeWhiteSpace(r *bufio.Reader) error {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		if !isWhitespace(b) {
			_ = r.UnreadByte()
			return nil
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
	tokenType := token.ByteToTokenType(b)

	if tokenType == token.Invalid {
		return token.Token{TokenType: token.Invalid, Literal: string(b)}
	}

	return readTokenOfType(r, tokenType)
}

func readTokenOfType(r *bufio.Reader, tokenType token.TokenType) token.Token {
	if tok, ok := token.TokenTypeToPredefinedToken(tokenType); ok {
		return tok
	}

	switch tokenType {

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
