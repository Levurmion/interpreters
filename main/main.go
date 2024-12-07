package main

import (
	"fmt"
	"interpreters/internal/grammar"
	"interpreters/internal/parser/firstfollow"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	Grammar, err := grammar.NewGrammarFromJsonConfig("./grammar-config.json")
	if (err != nil) {
		fmt.Println(err.Error())
	}

	FIRSTSets := firstfollow.ComputeFIRSTSets(Grammar)
	pretty := spew.Sdump(FIRSTSets)
	fmt.Println(pretty)
}