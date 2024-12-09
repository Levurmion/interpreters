package firstfollow

import (
	"interpreters/internal/parser/lr1grammar"
	"interpreters/utilities/sets"
)

func ComputeFIRSTSets(grammar *lr1grammar.Grammar) map[string]sets.Set[string] {
	FIRSTSets := make(map[string]sets.Set[string])

	// intialize FIRSTSets
	for _, nonTerminal := range grammar.NonTerminals.GetItems() {
		FIRSTSets[nonTerminal] = sets.NewEmptySet[string]()
	}
	for _, terminal := range grammar.Terminals.GetItems() {
		FIRSTSets[terminal] = sets.NewSet(terminal)
	}

	// Computes the `FIRST(symbol)` and returns whether `FIRST(symbol)`
	// was modified in this iteration
	FIRST := func (symbol string) bool {
		productionRules := grammar.GetProductionsOfNonTerminal(symbol)
		symbolFIRSTSet := FIRSTSets[symbol]
		newSymbolFIRSTSet := symbolFIRSTSet.Clone()

		for _, productionRule := range productionRules {
			for _, ruleSymbol := range productionRule.Production {
				// ruleSymbol is the current leading symbol: FIRST(ruleSymbol) in FIRST(symbol)
				ruleSymbolFIRSTSet := FIRSTSets[ruleSymbol]
				newSymbolFIRSTSet = newSymbolFIRSTSet.Union(ruleSymbolFIRSTSet)

				if (!grammar.DerivesEpsilon(ruleSymbol) || grammar.Terminals.Has(ruleSymbol)) {
					// leading symbol does not derive Epsilon or is a terminal: no other possible 
					// leading symbols for this `productionRule`
					break
				}
			}
		}

		FIRSTSets[symbol] = newSymbolFIRSTSet

		return newSymbolFIRSTSet.Size() > symbolFIRSTSet.Size()
	}

	// only stop refreshing FIRSTs when no more changes are made
	changed := true
	for changed {
		changed = false
		for _, nonTerminal := range grammar.NonTerminals.GetItems() {
			nonTerminalSetChanged := FIRST(nonTerminal)
			changed = changed || nonTerminalSetChanged
		}
	}

	return FIRSTSets
}