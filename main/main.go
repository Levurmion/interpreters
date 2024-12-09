package main

import (
	"fmt"
	"interpreters/internal/parser/lr1grammar"
	"interpreters/internal/parser/lr1item"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	Grammar, err := lr1grammar.NewAugmentedGrammarFromJsonConfig("./grammar-config.json")
	if (err != nil) {
		fmt.Println(err.Error())
	}
	
	fmt.Println(spew.Sdump(Grammar.ProductionRules))

	lr1Item, err := lr1item.NewLR1Item(
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