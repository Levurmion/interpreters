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
	FOLLOWSets := firstfollow.ComputeFOLLOWSets(Grammar, FIRSTSets)
	// FIRSTPretty := spew.Sdump(FIRSTSets)
	FOLLOWPretty := spew.Sdump(FOLLOWSets)
	// fmt.Println(FIRSTPretty)
	fmt.Println(FOLLOWPretty)
}