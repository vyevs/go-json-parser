package lex

import (
	"bufio"
	"strings"
	"testing"

	"github.com/vyevs/json/token"
)

func TestReadBoolLiteral(t *testing.T) {
	tests := []struct {
		str    string
		want   string
		wantOk bool
	}{
		{str: "true", want: "true", wantOk: true},
		{str: "false", want: "false", wantOk: true},
		{str: "s", want: "s", wantOk: false},
		{str: "abcdef", want: "abcde", wantOk: false},
		{str: " true", want: " true", wantOk: false},
		{str: "truefalse", want: "true", wantOk: true},
		{str: "pota", want: "pota", wantOk: false},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got, ok := readBoolLiteral(r)

		if ok != test.wantOk || got != test.want {
			t.Errorf("str: %q, want: %q, got: %q, wantOk: %v, got ok: %v",
				test.str, test.want, got, test.wantOk, ok)
		}
	}

}

func TestReadBoolToken(t *testing.T) {
	tests := []struct {
		str  string
		want token.Token
	}{
		{
			str:  "true",
			want: token.Token{TokenType: token.Boolean, Literal: "true"},
		},
		{
			str:  "false",
			want: token.Token{TokenType: token.Boolean, Literal: "false"},
		},
		{
			str:  "s",
			want: token.Token{TokenType: token.Invalid, Literal: "s"},
		},
		{
			str:  "abcdef",
			want: token.Token{TokenType: token.Invalid, Literal: "abcde"},
		},
		{
			str:  " true",
			want: token.Token{TokenType: token.Invalid, Literal: " true"},
		},
		{
			str:  "truefalse",
			want: token.Token{TokenType: token.Boolean, Literal: "true"},
		},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got := readBoolToken(r)

		if got != test.want {
			t.Errorf("str: %q, want: %q, got: %q", test.str, test.want, got)
		}
	}
}

func TestReadNullLiteral(t *testing.T) {
	tests := []struct {
		str    string
		want   string
		wantOk bool
	}{
		{str: "null", want: "null", wantOk: true},
		{str: "s", want: "s", wantOk: false},
		{str: " null", want: " nul", wantOk: false},
		{str: "nul l", want: "nul ", wantOk: false},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got, ok := readNullLiteral(r)

		if ok != test.wantOk || got != test.want {
			t.Errorf("str: %q, want: %q, got: %q, wantOk: %v, got ok: %v",
				test.str, test.want, got, test.wantOk, ok)
		}
	}
}

func TestReadNullToken(t *testing.T) {
	tests := []struct {
		str  string
		want token.Token
	}{
		{
			str:  "null",
			want: token.Token{TokenType: token.Null, Literal: "null"},
		},
		{
			str:  "s",
			want: token.Token{TokenType: token.Invalid, Literal: "s"},
		},
		{
			str:  " null",
			want: token.Token{TokenType: token.Invalid, Literal: " nul"},
		},
		{
			str:  "nul l",
			want: token.Token{TokenType: token.Invalid, Literal: "nul "},
		},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got := readNullToken(r)

		if got != test.want {
			t.Errorf("str:%q, want: %q, got: %q", test.str, test.want, got)
		}
	}
}

func TestReadNumericLiteral(t *testing.T) {
	tests := []struct {
		str    string
		want   string
		wantOk bool
	}{
		{str: "1", want: "1", wantOk: true},
		{str: "1234531231", want: "1234531231", wantOk: true},
		{str: "1.1", want: "1.1", wantOk: true},
		{str: "0.123455", want: "0.123455", wantOk: true},
		{str: "000010.21323", want: "000010.21323", wantOk: true},
		{str: " ", want: " ", wantOk: false},
		{str: "1222.23123.1", want: "1222.23123.", wantOk: false},
		{str: "123abc", want: "123", wantOk: true},
		{str: "-1", want: "-1", wantOk: true},
		{str: "-1235432132", want: "-1235432132", wantOk: true},
		{str: "-123.123", want: "-123.123", wantOk: true},
		{str: "-123.123.", want: "-123.123.", wantOk: false},
		{str: "-0.123", want: "-0.123", wantOk: true},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got, ok := readNumericLiteral(r)

		if ok != test.wantOk || got != test.want {
			t.Errorf("str: %q, want: %q, got: %q, wantOk: %v, got ok: %v",
				test.str, test.want, got, test.wantOk, ok)
		}
	}
}

func TestReadNumericToken(t *testing.T) {
	tests := []struct {
		str  string
		want token.Token
	}{
		{
			str:  "1",
			want: token.Token{TokenType: token.Integer, Literal: "1"},
		},
		{
			str:  "1234531231",
			want: token.Token{TokenType: token.Integer, Literal: "1234531231"},
		},
		{
			str:  "1.1",
			want: token.Token{TokenType: token.FloatingPoint, Literal: "1.1"},
		},
		{
			str:  "0.123455",
			want: token.Token{TokenType: token.FloatingPoint, Literal: "0.123455"},
		},
		{
			str:  "000010.21323",
			want: token.Token{TokenType: token.Invalid, Literal: "000010.21323"},
		},
		{
			str:  "abc",
			want: token.Token{TokenType: token.Invalid, Literal: "a"},
		},
		{
			str:  "1222.23123.1",
			want: token.Token{TokenType: token.Invalid, Literal: "1222.23123."},
		},
		{
			str:  "123abc",
			want: token.Token{TokenType: token.Integer, Literal: "123"},
		},
		{
			str:  "123.",
			want: token.Token{TokenType: token.Invalid, Literal: "123."},
		},
		{
			str:  "-1",
			want: token.Token{TokenType: token.Integer, Literal: "-1"},
		},
		{
			str:  "-1235432132",
			want: token.Token{TokenType: token.Integer, Literal: "-1235432132"},
		},
		{
			str:  "-123.123",
			want: token.Token{TokenType: token.FloatingPoint, Literal: "-123.123"},
		},
		{
			str:  "-123.123.",
			want: token.Token{TokenType: token.Invalid, Literal: "-123.123."},
		},
		{
			str:  "-0.123",
			want: token.Token{TokenType: token.FloatingPoint, Literal: "-0.123"},
		},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got := readNumericToken(r)

		if got != test.want {
			t.Errorf("str: %q, want: %q, got: %q", test.str, test.want, got)
		}
	}
}

func TestReadStringLiteral(t *testing.T) {
	tests := []struct {
		str    string
		want   string
		wantOk bool
	}{
		{str: `potato"`, want: "potato", wantOk: true},
		{str: `123abc123.123abc."`, want: `123abc123.123abc.`, wantOk: true},
		{str: `abc123"abc123`, want: "abc123", wantOk: true},
		{str: `"`, want: "", wantOk: true},
		{str: ``, want: "", wantOk: false},
		{str: `"""`, want: "", wantOk: true},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got, ok := readStringLiteral(r)

		if ok != test.wantOk || got != test.want {
			t.Errorf("str: %q, want: %q, got: %q, wantOk: %v, got ok: %v",
				test.str, test.want, got, test.wantOk, ok)
		}
	}
}

func TestReadStringToken(t *testing.T) {
	tests := []struct {
		str  string
		want token.Token
	}{
		{
			str:  `potato"`,
			want: token.Token{TokenType: token.String, Literal: "potato"},
		},
		{
			str:  `123abc123.123abc."`,
			want: token.Token{TokenType: token.String, Literal: "123abc123.123abc."},
		},
		{
			str:  `abc123"abc123`,
			want: token.Token{TokenType: token.String, Literal: "abc123"},
		},
		{
			str:  `   abc123  abc123   "`,
			want: token.Token{TokenType: token.String, Literal: "   abc123  abc123   "},
		},
		{
			str:  `"`,
			want: token.Token{TokenType: token.String, Literal: ""},
		},
		{
			str:  ``,
			want: token.Token{TokenType: token.Invalid, Literal: ""},
		},
		{
			str:  `"""`,
			want: token.Token{TokenType: token.String, Literal: ""},
		},
	}

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.str))

		got := readStringToken(r)

		if got != test.want {
			t.Errorf("str: %q, want: %q, got: %q", test.str, test.want, got)
		}
	}
}

func TestReadToken(t *testing.T) {
	tests := []struct {
		str  string
		want []token.Token
	}{
		{
			str: `{"a":"b"}`,
			want: []token.Token{
				token.OpeningCurlyBraceToken,
				token.Token{TokenType: token.String, Literal: "a"},
				token.ColonToken,
				token.Token{TokenType: token.String, Literal: "b"},
				token.ClosingCurlyBraceToken,
				token.EOFToken,
			},
		},
		{
			str: "   ",
			want: []token.Token{
				token.EOFToken,
			},
		},
		{
			str: `   "a"   123 true`,
			want: []token.Token{
				token.Token{TokenType: token.String, Literal: "a"},
				token.Token{TokenType: token.Integer, Literal: "123"},
				token.BooleanTrueToken,
				token.EOFToken,
			},
		},
		{
			str: "{a}",
			want: []token.Token{
				token.OpeningCurlyBraceToken,
				token.Token{TokenType: token.Invalid, Literal: "a"},
				token.ClosingCurlyBraceToken,
				token.EOFToken,
			},
		},
		{
			str: `{
  			"a": "1",
  			"b": 1,
  			"c": 1.0,
  			"d": {},
  			"e": [
  				{
  					"z": ["a", 0.15, false, null, true, -12, -0.123]
  				}
  			],
  			"f": null,
  			"g": true
			}`,
			want: []token.Token{
				token.OpeningCurlyBraceToken,
				token.Token{TokenType: token.String, Literal: "a"},
				token.ColonToken,
				token.Token{TokenType: token.String, Literal: "1"},
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "b"},
				token.ColonToken,
				token.Token{TokenType: token.Integer, Literal: "1"},
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "c"},
				token.ColonToken,
				token.Token{TokenType: token.FloatingPoint, Literal: "1.0"},
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "d"},
				token.ColonToken,
				token.OpeningCurlyBraceToken,
				token.ClosingCurlyBraceToken,
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "e"},
				token.ColonToken,
				token.OpeningSquareBracketToken,
				token.OpeningCurlyBraceToken,
				token.Token{TokenType: token.String, Literal: "z"},
				token.ColonToken,
				token.OpeningSquareBracketToken,
				token.Token{TokenType: token.String, Literal: "a"},
				token.CommaToken,
				token.Token{TokenType: token.FloatingPoint, Literal: "0.15"},
				token.CommaToken,
				token.BooleanFalseToken,
				token.CommaToken,
				token.NullToken,
				token.CommaToken,
				token.BooleanTrueToken,
				token.CommaToken,
				token.Token{TokenType: token.Integer, Literal: "-12"},
				token.CommaToken,
				token.Token{TokenType: token.FloatingPoint, Literal: "-0.123"},
				token.ClosingSquareBracketToken,
				token.ClosingCurlyBraceToken,
				token.ClosingSquareBracketToken,
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "f"},
				token.ColonToken,
				token.Token{TokenType: token.Null, Literal: "null"},
				token.CommaToken,
				token.Token{TokenType: token.String, Literal: "g"},
				token.ColonToken,
				token.BooleanTrueToken,
				token.ClosingCurlyBraceToken,
				token.EOFToken,
			},
		},
	}

	for _, test := range tests {
		lexer := New(strings.NewReader(test.str))

		got := readAllTokens(lexer)

		if !equalTokenSlices(test.want, got) {
			t.Errorf("str: %q, got(%d): %v, want(%d): %v", test.str, len(got), got, len(test.want), test.want)
		}
	}
}

func readAllTokens(lexer Lexer) []token.Token {
	tokens := make([]token.Token, 0)
	for {
		tok := lexer.ReadToken()
		tokens = append(tokens, tok)
		if tok == token.EOFToken {
			return tokens
		}
	}
}

func equalTokenSlices(t1, t2 []token.Token) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i, t := range t1 {
		if t != t2[i] {
			return false
		}
	}
	return true
}
