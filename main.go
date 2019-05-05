package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vyevs/json/parse"
)

func main() {
	f, _ := os.Open("sample.json")
	m, err := parse.Parse(f)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", m)
}
