package lex

import (
	"bufio"
	"strings"
	"testing"

	"github.com/vyevs/json/tok"
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
		want tok.Token
	}{
		{
			str:  "true",
			want: tok.Token{TokenType: tok.Boolean, Literal: "true"},
		},
		{
			str:  "false",
			want: tok.Token{TokenType: tok.Boolean, Literal: "false"},
		},
		{
			str:  "s",
			want: tok.Token{TokenType: tok.Invalid, Literal: "s"},
		},
		{
			str:  "abcdef",
			want: tok.Token{TokenType: tok.Invalid, Literal: "abcde"},
		},
		{
			str:  " true",
			want: tok.Token{TokenType: tok.Invalid, Literal: " true"},
		},
		{
			str:  "truefalse",
			want: tok.Token{TokenType: tok.Boolean, Literal: "true"},
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
		want tok.Token
	}{
		{
			str:  "null",
			want: tok.Token{TokenType: tok.Null, Literal: "null"},
		},
		{
			str:  "s",
			want: tok.Token{TokenType: tok.Invalid, Literal: "s"},
		},
		{
			str:  " null",
			want: tok.Token{TokenType: tok.Invalid, Literal: " nul"},
		},
		{
			str:  "nul l",
			want: tok.Token{TokenType: tok.Invalid, Literal: "nul "},
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
		{str: "", want: "", wantOk: false},
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
		want tok.Token
	}{
		{
			str:  "1",
			want: tok.Token{TokenType: tok.Integer, Literal: "1"},
		},
		{
			str:  "1234531231",
			want: tok.Token{TokenType: tok.Integer, Literal: "1234531231"},
		},
		{
			str:  "1.1",
			want: tok.Token{TokenType: tok.FloatingPoint, Literal: "1.1"},
		},
		{
			str:  "0.123455",
			want: tok.Token{TokenType: tok.FloatingPoint, Literal: "0.123455"},
		},
		{
			str:  "000010.21323",
			want: tok.Token{TokenType: tok.Invalid, Literal: "000010.21323"},
		},
		{
			str:  "abc",
			want: tok.Token{TokenType: tok.Invalid, Literal: "a"},
		},
		{
			str:  "1222.23123.1",
			want: tok.Token{TokenType: tok.Invalid, Literal: "1222.23123."},
		},
		{
			str:  "123abc",
			want: tok.Token{TokenType: tok.Integer, Literal: "123"},
		},
		{
			str:  "123.",
			want: tok.Token{TokenType: tok.Invalid, Literal: "123."},
		},
		{
			str:  "-1",
			want: tok.Token{TokenType: tok.Integer, Literal: "-1"},
		},
		{
			str:  "-1235432132",
			want: tok.Token{TokenType: tok.Integer, Literal: "-1235432132"},
		},
		{
			str:  "-123.123",
			want: tok.Token{TokenType: tok.FloatingPoint, Literal: "-123.123"},
		},
		{
			str:  "-123.123.",
			want: tok.Token{TokenType: tok.Invalid, Literal: "-123.123."},
		},
		{
			str:  "-0.123",
			want: tok.Token{TokenType: tok.FloatingPoint, Literal: "-0.123"},
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
		want tok.Token
	}{
		{
			str:  `potato"`,
			want: tok.Token{TokenType: tok.String, Literal: "potato"},
		},
		{
			str:  `123abc123.123abc."`,
			want: tok.Token{TokenType: tok.String, Literal: "123abc123.123abc."},
		},
		{
			str:  `abc123"abc123`,
			want: tok.Token{TokenType: tok.String, Literal: "abc123"},
		},
		{
			str:  `   abc123  abc123   "`,
			want: tok.Token{TokenType: tok.String, Literal: "   abc123  abc123   "},
		},
		{
			str:  `"`,
			want: tok.Token{TokenType: tok.String, Literal: ""},
		},
		{
			str:  ``,
			want: tok.Token{TokenType: tok.Invalid, Literal: ""},
		},
		{
			str:  `"""`,
			want: tok.Token{TokenType: tok.String, Literal: ""},
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
		want []tok.Token
	}{
		{
			str: `{"a":"b"}`,
			want: []tok.Token{
				tok.OpeningCurlyBraceToken,
				{TokenType: tok.String, Literal: "a"},
				tok.ColonToken,
				{TokenType: tok.String, Literal: "b"},
				tok.ClosingCurlyBraceToken,
				tok.EOFToken,
			},
		},
		{
			str: "   ",
			want: []tok.Token{
				tok.EOFToken,
			},
		},
		{
			str: `   "a"   123 true`,
			want: []tok.Token{
				{TokenType: tok.String, Literal: "a"},
				{TokenType: tok.Integer, Literal: "123"},
				tok.BooleanTrueToken,
				tok.EOFToken,
			},
		},
		{
			str: "{a}",
			want: []tok.Token{
				tok.OpeningCurlyBraceToken,
				{TokenType: tok.Invalid, Literal: "a"},
				tok.ClosingCurlyBraceToken,
				tok.EOFToken,
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
			want: []tok.Token{
				tok.OpeningCurlyBraceToken,
				{TokenType: tok.String, Literal: "a"},
				tok.ColonToken,
				{TokenType: tok.String, Literal: "1"},
				tok.CommaToken,
				{TokenType: tok.String, Literal: "b"},
				tok.ColonToken,
				{TokenType: tok.Integer, Literal: "1"},
				tok.CommaToken,
				{TokenType: tok.String, Literal: "c"},
				tok.ColonToken,
				{TokenType: tok.FloatingPoint, Literal: "1.0"},
				tok.CommaToken,
				{TokenType: tok.String, Literal: "d"},
				tok.ColonToken,
				tok.OpeningCurlyBraceToken,
				tok.ClosingCurlyBraceToken,
				tok.CommaToken,
				{TokenType: tok.String, Literal: "e"},
				tok.ColonToken,
				tok.OpeningSquareBracketToken,
				tok.OpeningCurlyBraceToken,
				{TokenType: tok.String, Literal: "z"},
				tok.ColonToken,
				tok.OpeningSquareBracketToken,
				{TokenType: tok.String, Literal: "a"},
				tok.CommaToken,
				{TokenType: tok.FloatingPoint, Literal: "0.15"},
				tok.CommaToken,
				tok.BooleanFalseToken,
				tok.CommaToken,
				tok.NullToken,
				tok.CommaToken,
				tok.BooleanTrueToken,
				tok.CommaToken,
				{TokenType: tok.Integer, Literal: "-12"},
				tok.CommaToken,
				{TokenType: tok.FloatingPoint, Literal: "-0.123"},
				tok.ClosingSquareBracketToken,
				tok.ClosingCurlyBraceToken,
				tok.ClosingSquareBracketToken,
				tok.CommaToken,
				{TokenType: tok.String, Literal: "f"},
				tok.ColonToken,
				{TokenType: tok.Null, Literal: "null"},
				tok.CommaToken,
				{TokenType: tok.String, Literal: "g"},
				tok.ColonToken,
				tok.BooleanTrueToken,
				tok.ClosingCurlyBraceToken,
				tok.EOFToken,
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

func readAllTokens(lexer Lexer) []tok.Token {
	toks := make([]tok.Token, 0)
	for {
		t := lexer.ReadToken()
		toks = append(toks, t)
		if t == tok.EOFToken {
			return toks
		}
	}
}

func equalTokenSlices(t1, t2 []tok.Token) bool {
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
