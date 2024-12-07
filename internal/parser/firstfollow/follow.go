package firstfollow

import (
	"interpreters/internal/grammar"
	"interpreters/internal/symbols"
	"interpreters/utilities/arrays"
	"interpreters/utilities/sets"
)

func ComputeFOLLOWSets(grammar *grammar.Grammar, FIRSTSets map[string]sets.Set[string]) map[string]sets.Set[string] {
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

			// find the position of this `symbol` in the production rule that derives it
			symbolIdx := arrays.FindFirstIdx(production, func (ruleSymbol string) bool {
				return ruleSymbol == symbol
			})

			if symbolIdx < 0 {
				continue
			}

			// `symbol` occurs last in the `production`
			if symbolIdx == len(production) - 1 {
				newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(nonTerminalFOLLOWSet)
				continue
			}

			// The `FOLLOWSymbol` is the symbol directly to the right of `symbol`
			FOLLOWSymbol := production[symbolIdx + 1]
			FOLLOWSymbolFIRSTSet := FIRSTSets[FOLLOWSymbol]
			
			// `FIRST(FOLLOWSymbol)` is always in `FOLLOW(symbol)`
			newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(FOLLOWSymbolFIRSTSet)
			
			if grammar.DerivesEpsilon(FOLLOWSymbol) && grammar.NonTerminals.Has(FOLLOWSymbol) {
				// `FOLLOWSymbol` can derive Epsilon: the `FOLLOW(FOLLOWSymbol)` is in `FOLLOW(symbol)`
				newSymbolFOLLOWSet = newSymbolFOLLOWSet.Union(nonTerminalFOLLOWSet)
			}
		}

		// always exclude Epsilon from FOLLOWSets
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