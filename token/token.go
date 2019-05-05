package token

import "fmt"

type TokenType int

const (
	OpeningCurlyBrace TokenType = iota
	ClosingCurlyBrace

	OpeningSquareBracket
	ClosingSquareBracket

	Colon

	Comma

	String

	Integer

	FloatingPoint

	Boolean

	Null

	EOF

	Invalid
)

type Token struct {
	TokenType TokenType
	Literal   string
}

func (t Token) String() string {
	return fmt.Sprintf("{%s, %q}", t.TokenType, t.Literal)
}

// predefined tokens whose literals are always the same
var (
	OpeningCurlyBraceToken    = Token{TokenType: OpeningCurlyBrace, Literal: "{"}
	ClosingCurlyBraceToken    = Token{TokenType: ClosingCurlyBrace, Literal: "}"}
	OpeningSquareBracketToken = Token{TokenType: OpeningSquareBracket, Literal: "["}
	ClosingSquareBracketToken = Token{TokenType: ClosingSquareBracket, Literal: "]"}
	ColonToken                = Token{TokenType: Colon, Literal: ":"}
	CommaToken                = Token{TokenType: Comma, Literal: ","}
	EOFToken                  = Token{TokenType: EOF, Literal: "EOF"}
	NullToken                 = Token{TokenType: Null, Literal: "null"}
	BooleanTrueToken          = Token{TokenType: Boolean, Literal: "true"}
	BooleanFalseToken         = Token{TokenType: Boolean, Literal: "false"}
)

var tokenTypeToPredefinedToken = map[TokenType]Token{
	OpeningCurlyBrace:    OpeningCurlyBraceToken,
	ClosingCurlyBrace:    ClosingCurlyBraceToken,
	OpeningSquareBracket: OpeningSquareBracketToken,
	ClosingSquareBracket: ClosingSquareBracketToken,
	Colon:                ColonToken,
	Comma:                CommaToken,
}

// for use with single character tokens
func TokenTypeToPredefinedToken(tokenType TokenType) (Token, bool) {
	tok, ok := tokenTypeToPredefinedToken[tokenType]
	return tok, ok
}

var tokenTypeToString = map[TokenType]string{
	OpeningCurlyBrace:    "OpeningCurlyBrace",
	ClosingCurlyBrace:    "ClosingCurlyBrace",
	OpeningSquareBracket: "OpeningSquareBrace",
	ClosingSquareBracket: "ClosingSquareBrace",
	Colon:                "Colon",
	Comma:                "Comma",

	String:        "StringLiteral",
	Integer:       "Integer",
	FloatingPoint: "FloatingPoint",
	Boolean:       "Boolean",
	Null:          "Null",
	EOF:           "EOF",
	Invalid:       "Invalid",
}

func (tokType TokenType) String() string {
	str, ok := tokenTypeToString[tokType]
	if ok {
		return str
	}
	return "Unknown TokenType"
}

var byteToTokenType = map[byte]TokenType{
	'{': OpeningCurlyBrace,
	'}': ClosingCurlyBrace,
	'[': OpeningSquareBracket,
	']': ClosingSquareBracket,
	'"': String,
	':': Colon,
	',': Comma,
	'n': Null,
	't': Boolean,
	'f': Boolean,
	'-': Integer,
}

// ByteToTokenType returns what the next TokenType will be
// if the byte b is encountered
func ByteToTokenType(b byte) TokenType {
	if tokType, ok := byteToTokenType[b]; ok {
		return tokType
	}
	if b >= '0' && b <= '9' {
		return Integer
	}
	return Invalid
}
