package lex

import (
	"strings"
	"testing"

	"github.com/vyevs/gojson/tok"
)

func TestReadNumericBytes(t *testing.T) {
	tests := []struct {
		literal  string
		wantStr  string
		wantType tok.Type
	}{
		{
			literal: "1", wantStr: "1", wantType: tok.Integer,
		},
		{
			literal: "1232452142", wantStr: "1232452142", wantType: tok.Integer,
		},
		{
			literal: "-1234567", wantStr: "-1234567", wantType: tok.Integer,
		},
		{
			// turns out this is valid
			literal: "-0", wantStr: "-0", wantType: tok.Integer,
		},
		{
			literal: "123.123", wantStr: "123.123", wantType: tok.Float,
		},
		{
			literal: "-5091.231230", wantStr: "-5091.231230", wantType: tok.Float,
		},
		{
			literal: "-0.0000", wantStr: "-0.0000", wantType: tok.Float,
		},
		{
			literal: "-0.12345", wantStr: "-0.12345", wantType: tok.Float,
		},
		{
			literal: "01", wantType: tok.Invalid,
		},
		{
			literal: "-01", wantType: tok.Invalid,
		},
		{
			literal: "", wantType: tok.Invalid,
		},
	}

	for _, test := range tests {
		l := New(strings.NewReader(test.literal))

		wantBytes := []byte(test.wantStr)

		gotType := l.readNumericBytes()

		gotBytes := l.GetTokenBytes()

		if gotType != test.wantType || test.wantStr != "" && !equalSlices(gotBytes, wantBytes) {
			t.Errorf("literal: %q, got: %q, want: %q, gotType: %s, wantType: %s",
				test.literal, gotBytes, test.wantStr, gotType, test.wantType)
		}
	}
}

func equalSlices(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, b := range s1 {
		if b != s2[i] {
			return false
		}
	}
	return true
}

func TestReadToken(t *testing.T) {
	type want struct {
		typ tok.Type
		str string
	}
	tests := []struct {
		literal string
		want    []want
	}{
		{
			literal: `{
  "firstName": "John",
  "lastName": "Smith",
  "male": true,
  "from123": false,
  "n": null,
  "age": 25,
  "money": -569.231,
  "address": {
    "streetAddress": "21 2nd Street",
    "city": "New York",
    "state": "NY",
    "postalCode": 10021
  }}
`,
			want: []want{
				want{typ: tok.OpeningCurlyBrace},
				want{typ: tok.String, str: "firstName"},
				want{typ: tok.Colon},
				want{typ: tok.String, str: "John"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "lastName"},
				want{typ: tok.Colon},
				want{typ: tok.String, str: "Smith"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "male"},
				want{typ: tok.Colon},
				want{typ: tok.True},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "from123"},
				want{typ: tok.Colon},
				want{typ: tok.False},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "n"},
				want{typ: tok.Colon},
				want{typ: tok.Null},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "age"},
				want{typ: tok.Colon},
				want{typ: tok.Integer, str: "25"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "money"},
				want{typ: tok.Colon},
				want{typ: tok.Float, str: "-569.231"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "address"},
				want{typ: tok.Colon},
				want{typ: tok.OpeningCurlyBrace},
				want{typ: tok.String, str: "streetAddress"},
				want{typ: tok.Colon},
				want{typ: tok.String, str: "21 2nd Street"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "city"},
				want{typ: tok.Colon},
				want{typ: tok.String, str: "New York"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "state"},
				want{typ: tok.Colon},
				want{typ: tok.String, str: "NY"},
				want{typ: tok.Comma},
				want{typ: tok.String, str: "postalCode"},
				want{typ: tok.Colon},
				want{typ: tok.Integer, str: "10021"},
				want{typ: tok.ClosingCurlyBrace},
				want{typ: tok.ClosingCurlyBrace},
				want{typ: tok.EOF},
			},
		},
	}

	for _, test := range tests {
		l := New(strings.NewReader(test.literal))

		for _, want := range test.want {
			got := l.ReadToken()

			gotBytes := l.GetTokenBytes()

			if got != want.typ ||
				(want.str != "" && !equalSlices(gotBytes, []byte(want.str))) {
				t.Fatalf("got: %s, want: %s, wantStr: %q, gotStr: %q",
					got, want.typ, want.str, string(gotBytes))
			}
		}
	}
}
