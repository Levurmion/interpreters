package lr1_items

import (
	"interpreters/internal/grammar"
	"interpreters/internal/parser/firstfollow"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"
)

// A factory object responsible for creating `LR1Items` based on the given `Grammar`.
// `LR1ItemFactory` will automatically compute the lookahead sets for each produced
// `LR1Item`.
type LR1ItemFactory struct {
	Grammar *grammar.Grammar
	FIRSTSets map[string]sets.Set[string]
}

func NewLR1ItemFactory(grammar *grammar.Grammar) *LR1ItemFactory {
	FIRSTSets := firstfollow.ComputeFIRSTSets(grammar)
	return &LR1ItemFactory{
		grammar,
		FIRSTSets,
	}
}

// Computes the lookahead set for a given context: sequence of symbols following a
// non-terminal to the RHS of the parsing progress (the dot).
func (factory *LR1ItemFactory) Lookahead(context []string, currLookahead sets.Set[string]) sets.Set[string] {
	if len(context) == 0 {
		return currLookahead.Clone()
	} else if context[0] == symbols.Epsilon {
		return currLookahead.Clone()
	} else {
		lookaheadSet := sets.NewEmptySet[string]()
		for _, symbol := range context {
			if factory.Grammar.Terminals.Has(symbol) {
				// symbol is a terminal: this is the only possible lookahead
				lookaheadSet.Add(symbol)
				break
			} else {
				// symbol is a non-terminal: FIRST(symbol) is in lookaheadSet
				symbolFIRSTSet := factory.FIRSTSets[symbol]
				lookaheadSet = lookaheadSet.Union(symbolFIRSTSet)
				if !factory.Grammar.DerivesEpsilon(symbol) {
					// non-terminal cannot derive Epsilon: no other possible lookaheads
					return lookaheadSet
				}
			}
		}

		// reached the end of `context` and all previous symbols could derive Epsilon
		// so we also include the parent item's lookahead set
		return lookaheadSet.Union(currLookahead)
	}
}



