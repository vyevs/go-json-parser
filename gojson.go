package gojson

import (
	"bytes"
	"io"
	"strings"

	"github.com/vyevs/gojson/parse"
)

// Parse parses the bytes in the reader into a json object
//
// since null, a string, an array, a number by themselves
// are valid json docs, the return value can be either
// nil, string, []interface{}, int, float64
// or for most cases, a map[string]interface{}
// representing a json object
func Parse(r io.Reader) (interface{}, error) {
	return parse.Parse(r)
}

// ParseStr is a convenience method
// Does same thing as Parse but reading from a string
func ParseStr(str string) (interface{}, error) {
	r := strings.NewReader(str)

	return parse.Parse(r)
}

func ParseBytes(b []byte) (interface{}, error) {
	r := bytes.NewReader(b)

	return parse.Parse(r)
}
