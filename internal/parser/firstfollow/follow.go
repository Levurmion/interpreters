package firstfollow

import (
	"interpreters/internal/parser/lr1grammar"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"
)

func ComputeFOLLOWSets(grammar *lr1grammar.Grammar, FIRSTSets map[string]sets.Set[string]) map[string]sets.Set[string] {
	FOLLOWSets := make(map[string]sets.Set[string])
	grammarNonTerminals := grammar.NonTerminals.GetItems()

	// initialize FOLLOWSets
	for _, nonTerminal := range grammarNonTerminals {
		if nonTerminal == grammar.StartSymbol {
			FOLLOWSets[nonTerminal] = sets.NewSet(symbols.EOF)
		} else {
			FOLLOWSets[nonTerminal] = sets.NewEmptySet[string]()
		}
	}

	FOLLOW := func (symbol string) bool {
		symbolFOLLOWSet := FOLLOWSets[symbol]
		newSymbolFOLLOWSet := symbolFOLLOWSet.Clone()
		productionRulesDerivingSymbol := grammar.GetProductionsDerivingSymbol(symbol)

		for _, productionRule := range productionRulesDerivingSymbol {
			productionNonTerminal := productionRule.NonTerminal
			production := productionRule.Production
			nonTerminalFOLLOWSet := FOLLOWSets[productionNonTerminal]

			// Iterate over all occurrences of 'symbol' in the production
			for idx := 0; idx < len(production); idx++ {
				if production[idx] != symbol {
					continue
				}

				// If 'symbol' is at the end of the production, add FOLLOW(LHS) to FOLLOW(symbol)
				if idx == len(production)-1 {
					newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(nonTerminalFOLLOWSet)
				} else {
					// The symbol that follows 'symbol'
					nextSymbol := production[idx+1]
					nextSymbolFIRSTSet := FIRSTSets[nextSymbol]

					// Add FIRST(nextSymbol) to FOLLOW(symbol)
					newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(nextSymbolFIRSTSet)

					// If nextSymbol can derive epsilon (and is a NonTerminal),
					// then FOLLOW(LHS) also goes into FOLLOW(symbol)
					if grammar.NonTerminals.Has(nextSymbol) && grammar.DerivesEpsilon(nextSymbol) {
						newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(nonTerminalFOLLOWSet)
					}
				}
			}
		}

		// Epsilon should never appear in FOLLOW sets
		newSymbolFOLLOWSet.Delete(symbols.Epsilon)
		FOLLOWSets[symbol] = newSymbolFOLLOWSet
		return newSymbolFOLLOWSet.Size() > symbolFOLLOWSet.Size()
	}
	
	changed := true
	for changed {
		changed = false
		for _, nonTerminal := range grammarNonTerminals {
			nonTerminalFOLLOWSetChanged := FOLLOW(nonTerminal)
			changed = changed || nonTerminalFOLLOWSetChanged
		}
	}

	return FOLLOWSets
}