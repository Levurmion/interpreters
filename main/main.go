package main

import (
	"fmt"
	"interpreters/internal/grammar"
)

func main() {
	Grammar, err := grammar.NewGrammarFromJsonConfig("./grammar-config.json")
	if (err != nil) {
		fmt.Println(err.Error())
	}

	fmt.Println(Grammar.GetProductionsDerivingSymbol("str"))
}