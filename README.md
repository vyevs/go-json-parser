[![Go Report Card](https://goreportcard.com/badge/github.com/vyevs/go-json-parser)](https://goreportcard.com/report/github.com/vyevs/go-json-parser)

# go-json-parser
go-json-parser

Simple streaming JSON parser.

Parses JSON documents into `map[string]interface{}`, or for single value docs (array, string, numeric, bool, null) into either 
`[]interface{}`, `string`, `int` or `float64`, `bool`, or `nil`, respectively

Currently slightly slower than the encoding/json Decoder

The `Parse(io.Reader)` and `ParseStr(string)` functions are the interface provided for JSON parsing

Run tests:

`GOPATH/github.com/vyevs/json> test ./...`

Run benchmark that's setup against encoding/json:

`GOPATH/github.com/vyevs/json> test bench=. ./...`

NOTE: Benchmarks may take some time to run, as they run against the test files in parse/testdata
Any files added to testdata will be tested against
