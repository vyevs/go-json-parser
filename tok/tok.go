package tok

// TokenType represents a sequence of characters in a JSON document
type Type int

// the TokenTypes that a valid json doc will contain (not including Invalid)
const (
	OpeningCurlyBrace Type = iota
	ClosingCurlyBrace

	OpeningSquareBracket
	ClosingSquareBracket

	Colon

	Comma

	String

	Integer

	Float

	True
	False

	Null

	EOF

	Invalid
)

func (t Type) String() string {
	if t == OpeningCurlyBrace {
		return "OpeningCurlyBrace"
	} else if t == ClosingCurlyBrace {
		return "ClosingCurlyBrace"
	} else if t == OpeningSquareBracket {
		return "OpeningSquareBracket"
	} else if t == Colon {
		return "Colon"
	} else if t == Comma {
		return "Comma"
	} else if t == String {
		return "String"
	} else if t == Integer {
		return "Integer"
	} else if t == Float {
		return "FloatingPoint"
	} else if t == True {
		return "BooleanTrue"
	} else if t == False {
		return "BooleanFalse"
	} else if t == Null {
		return "Null"
	} else if t == EOF {
		return "EOF"
	} else if t == Invalid {
		return "Invalid"
	}
	return "TOK TYPE STRING() RECEIVED UNKNOWN TYPE"
}

// ByteToTokenType returns what the next TokenType will be
// once byte b is encountered
func ByteToType(b byte) Type {
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
		return True
	} else if b == 'f' {
		return False
	} else if b == '"' {
		return String
	} else if b == '-' || b >= '0' && b <= '9' {
		return Integer
	} else {
		return Invalid
	}
}
