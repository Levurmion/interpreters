package main

import (
	"fmt"
	"interpreters/internal/grammar"
	"interpreters/internal/parser/lr1_items"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	Grammar, err := grammar.NewAugmentedGrammarFromJsonConfig("./grammar-config.json")
	if (err != nil) {
		fmt.Println(err.Error())
	}
	
	fmt.Println(spew.Sdump(Grammar.ProductionRules))

	lr1Item, err := lr1_items.NewLR1Item(
		"OBJECT", 
		[]string{"{", "ENTRIES?", "}"}, 
		0, 
		sets.NewSet[string](symbols.EOF),
	)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(lr1Item.GetName())
		fmt.Println(lr1Item.GetContextForNextSymbol())
		nextItem, err := lr1Item.AdvanceDot()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(nextItem.GetName())
			fmt.Println(nextItem.LookaheadSet.GetItems())
		}
	}
}