package parse

import (
	"strings"
	"testing"

	"github.com/vyevs/json/lex"
)

func TestParseArray(t *testing.T) {
	tests := []struct {
		str     string
		want    []interface{}
		wantErr bool
	}{
		{str: "]", want: []interface{}{}, wantErr: false},
		{str: `"a"]`, want: []interface{}{"a"}, wantErr: false},
		{str: `"a", "b", "c"]`, want: []interface{}{"a", "b", "c"}, wantErr: false},
		{str: "5213]", want: []interface{}{5213}, wantErr: false},
		{str: `1, 2, 3, -3, -2, -1, 0]`, want: []interface{}{1, 2, 3, -3, -2, -1, 0}, wantErr: false},
		{str: "102.123]", want: []interface{}{102.123}, wantErr: false},
		{str: `0.1, 1.01, 10.001, -1.0, -0.534]`, want: []interface{}{0.1, 1.01, 10.001, -1.0, -0.534}, wantErr: false},
		{str: `null, null, null]`, want: []interface{}{nil, nil, nil}, wantErr: false},
	}

	for _, test := range tests {
		l := lex.New(strings.NewReader(test.str))
		got, err := parseArray(l)

		gotErr := err != nil

		if gotErr != test.wantErr || !slicesEqual(got, test.want) {
			t.Errorf("str: %q, got: %v, want: %v, err: %v, wantErr: %v",
				test.str, got, test.want, err, test.wantErr)
		}
	}
}

func slicesEqual(s1, s2 []interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}
