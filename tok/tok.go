package tok

// TokenType represents a sequence of characters in a JSON document
type TokenType int

// the TokenTypes that a valid json doc will contain (not including Invalid)
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

	BooleanTrue
	BooleanFalse

	Null

	EOF

	Invalid
)

// ByteToTokenType returns what the next TokenType will be
// once byte b is encountered
func ByteToTokenType(b byte) TokenType {
	if b == '{' {
		return OpeningCurlyBrace
	} else if b == '}' {
		return ClosingCurlyBrace
	} else if b == '[' {
		return OpeningSquareBracket
	} else if b == ']' {
		return ClosingSquareBracket
	} else if b == ':' {
		return Colon
	} else if b == ',' {
		return Comma
	} else if b == 'n' {
		return Null
	} else if b == 't' {
		return BooleanTrue
	} else if b == 'f' {
		return BooleanFalse
	} else if b == '"' {
		return String
	} else if b == '-' || b >= '0' && b <= '9' {
		return Integer
	} else {
		return Invalid
	}
}
