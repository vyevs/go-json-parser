package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vyevs/gojson/lex"
)

func TestParseInteger(t *testing.T) {
	tests := []struct {
		literal string
		want    int
		wantErr bool
	}{
		{literal: "1", want: 1},
		{literal: "-512", want: -512},
		{literal: "123455623123123213213213", wantErr: true},
	}

	for _, test := range tests {
		got, err := parseInteger(test.literal)

		gotErr := err != nil
		if gotErr != test.wantErr || got != test.want {
			t.Errorf("literal: %q, got: %d, want: %d, err: %v, wantErr: %v",
				test.literal, got, test.want, err, test.wantErr)
		}
	}
}

func TestParseFloatingPoint(t *testing.T) {
	tests := []struct {
		literal string
		want    float64
		wantErr bool
	}{
		{literal: "1.123", want: 1.123},
		{literal: "-512.98231", want: -512.98231},
	}

	for _, test := range tests {
		got, err := parseFloatingPoint(test.literal)

		gotErr := err != nil
		if gotErr != test.wantErr || got != test.want {
			t.Errorf("literal: %q, got: %f, want: %f, err: %v, wantErr: %v",
				test.literal, got, test.want, err, test.wantErr)
		}
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		literal string
		want    bool
		wantErr bool
	}{
		{literal: "true", want: true},
		{literal: "false", want: false},
		{literal: "ball please", wantErr: true},
	}

	for _, test := range tests {
		got, err := parseBool(test.literal)

		gotErr := err != nil
		if gotErr != test.wantErr || got != test.want {
			t.Errorf("literal: %q, got: %v, want: %v, err: %v, wantErr: %v",
				test.literal, got, test.want, err, test.wantErr)
		}
	}
}

// TestParseObject & TestParseArray are intertwined, in that if one fails
// the other should fail too
func TestParseObject(t *testing.T) {
	tests := []struct {
		literal string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			literal: `"a": "potato", "b": 1.2345, "grass": null}`,
			want: map[string]interface{}{
				"a":     "potato",
				"b":     1.2345,
				"grass": nil,
			},
		},
		{
			literal: `"obj": {"array": ["val", "grass knoll"]}}`,
			want: map[string]interface{}{
				"obj": map[string]interface{}{
					"array": []interface{}{
						"val",
						"grass knoll",
					},
				},
			},
		},
		{literal: `"a": "b", "a": "b"}`, wantErr: true},
		{literal: `"b": "c" null}`, wantErr: true},
		{literal: `"a" null}`, wantErr: true},
		{literal: `a: "b"}`, wantErr: true},
		{literal: `"a": fals}`, wantErr: true},
	}

	for _, test := range tests {
		l := lex.New(strings.NewReader(test.literal))
		got, err := parseObject(l)

		gotErr := err != nil

		if gotErr != test.wantErr || !mapsEqual(got, test.want) {
			t.Errorf("literal: %q, got: %v, want: %v, err: %v, wantErr: %v",
				test.literal, got, test.want, err, test.wantErr)
		}
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		str     string
		want    []interface{}
		wantErr bool
	}{
		{str: "]", want: []interface{}{}},
		{str: `"a"]`, want: []interface{}{"a"}},
		{str: `"a", "b", "c"]`, want: []interface{}{"a", "b", "c"}},
		{str: "5213]", want: []interface{}{5213}},
		{str: `1, 2, 3, -3, -2, -1, 0]`, want: []interface{}{1, 2, 3, -3, -2, -1, 0}},
		{str: "-102.123]", want: []interface{}{-102.123}},
		{str: `0.1, 1.01, 10.001, -1.0, -0.534]`, want: []interface{}{0.1, 1.01, 10.001, -1.0, -0.534}},
		{str: `null]`, want: []interface{}{nil}},
		{str: `null, null, null]`, want: []interface{}{nil, nil, nil}},
		{str: `true]`, want: []interface{}{true}},
		{str: `false]`, want: []interface{}{false}},
		{str: `true, false, true, false]`, want: []interface{}{true, false, true, false}},
		{
			str:  `1, 1.0, false, 2323, "grass", null, 32]`,
			want: []interface{}{1, 1.0, false, 2323, "grass", nil, 32},
		},
		{
			str: `{"a": "b", "b": 1, "c": null}, {"a": false, "b": true}]`,
			want: []interface{}{
				map[string]interface{}{
					"a": "b",
					"b": 1,
					"c": nil,
				},
				map[string]interface{}{
					"a": false,
					"b": true,
				},
			},
		},
		{
			str: `{}, [{}, 53.2, null, {"a": 55}]]`,
			want: []interface{}{
				map[string]interface{}{},
				[]interface{}{
					map[string]interface{}{},
					53.2,
					nil,
					map[string]interface{}{
						"a": 55,
					},
				},
			},
		},
		{str: ``, want: nil, wantErr: true},
		{str: `1, 2, 3 4]`, want: nil, wantErr: true},
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

func TestParse(t *testing.T) {
	tests := []struct {
		literal string
		want    interface{}
		wantErr bool
	}{
		{literal: `"abc"`, want: "abc"},
		{literal: `123`, want: 123},
		{literal: `-124.231`, want: -124.231},
		{literal: `null`, want: nil},
		{literal: `true`, want: true},
		{
			literal: `[1, 2.0, null, false, true, "p", [5]]`,
			want:    []interface{}{1, 2.0, nil, false, true, "p", []interface{}{5}},
		},
		{literal: `nullp`, wantErr: true},
		{literal: `potato`, wantErr: true},
		{literal: ``, wantErr: true},
		{
			literal: `{
      "Actors": [
        {
          "name": "Tom Cruise",
          "age": 56,
          "Born At": "Syracuse, NY",
          "Birthdate": "July 3, 1962",
          "photo": "https://jsonformatter.org/img/tom-cruise.jpg",
          "wife": null,
          "weight": 67.5,
          "hasChildren": true,
          "hasGreyHair": false,
          "children": [
            "Suri",
            "Isabella Jane",
            "Connor"
          ]
        },
        {
          "name": "Robert Downey Jr.",
          "age": 53,
          "Born At": "New York City, NY",
          "Birthdate": "April 4, 1965",
          "photo": "https://jsonformatter.org/img/Robert-Downey-Jr.jpg",
          "wife": "Susan Downey",
          "weight": 77.1,
          "hasChildren": true,
          "hasGreyHair": false,
          "children": [
            "Indio Falconer",
            "Avri Roel",
            "Exton Elias"
          ]
        }
      ]
    }`,
			want: map[string]interface{}{
				"Actors": []interface{}{
					map[string]interface{}{
						"name":        "Tom Cruise",
						"age":         56,
						"Born At":     "Syracuse, NY",
						"Birthdate":   "July 3, 1962",
						"photo":       "https://jsonformatter.org/img/tom-cruise.jpg",
						"wife":        nil,
						"weight":      67.5,
						"hasChildren": true,
						"hasGreyHair": false,
						"children": []interface{}{
							"Suri",
							"Isabella Jane",
							"Connor",
						},
					},
					map[string]interface{}{
						"name":        "Robert Downey Jr.",
						"age":         53,
						"Born At":     "New York City, NY",
						"Birthdate":   "April 4, 1965",
						"photo":       "https://jsonformatter.org/img/Robert-Downey-Jr.jpg",
						"wife":        "Susan Downey",
						"weight":      77.1,
						"hasChildren": true,
						"hasGreyHair": false,
						"children": []interface{}{
							"Indio Falconer",
							"Avri Roel",
							"Exton Elias",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := Parse(strings.NewReader(test.literal))

		gotErr := err != nil

		if gotErr != test.wantErr || !equal(got, test.want) {
			t.Errorf("literal: %q, got: %v, want: %v, err: %v, wantErr: %v",
				test.literal, got, test.want, err, test.wantErr)
		}
	}
}

func slicesEqual(s1, s2 []interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if !equal(v, s2[i]) {
			return false
		}
	}
	return true
}

func mapsEqual(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		v2, ok := m2[k]
		if !ok || !equal(v, v2) {
			return false
		}
	}
	return true
}

func equal(v1, v2 interface{}) bool {
	switch actualV1 := v1.(type) {
	case []interface{}:
		actualV2, ok := v2.([]interface{})
		if !ok {
			return false
		}
		return slicesEqual(actualV1, actualV2)
	case map[string]interface{}:
		actualV2, ok := v2.(map[string]interface{})
		if !ok {
			return false
		}
		return mapsEqual(actualV1, actualV2)
	default:
		return v1 == v2
	}
}

func BenchmarkParse(b *testing.B) {
	paths, err := getTestFilePaths()
	if err != nil {
		b.Fatalf("getTestFilePaths(): %v", err)
	}

	for _, path := range paths {
		b.Run(path, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f, err := os.Open(path)
				if err != nil {
					b.Fatal(err)
				}
				_, err = Parse(f)
				if err != nil {
					b.Fatalf("Unexpected Parse() failure: %v", err)
				}
				f.Close()
			}
		})
		b.Run(fmt.Sprintf("%s%s", path, "STDLIB"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f, err := os.Open(path)
				if err != nil {
					b.Fatal(err)
				}
				d := json.NewDecoder(f)
				var v interface{}
				err = d.Decode(&v)
				if err != nil {
					b.Fatalf("Unexpected Decode() failure: %v", err)
				}
				f.Close()
			}
		})
	}
}

func getTestFilePaths() ([]string, error) {
	out := make([]string, 0)
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		out = append(out, path)
		return nil
	})
	return out, err
}
