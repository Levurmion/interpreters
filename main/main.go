package main

import (
	"fmt"
	"interpreters/internal/parser/lr1grammar"
	"interpreters/internal/parser/lr1parsingtable"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	Grammar, err := lr1grammar.NewAugmentedGrammarFromJsonConfig("./grammar-config.json")
	if (err != nil) {
		fmt.Println(err.Error())
	}
	
	automaton, err := lr1parsingtable.NewLR1Automaton(Grammar)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(spew.Sdump(automaton))
	}
}