package lr1parsingtable

import (
	"interpreters/internal/parser/firstfollow"
	"interpreters/internal/parser/lr1grammar"
	"interpreters/internal/parser/lr1item"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"
)

type ParsingTable struct {
	Grammar *lr1grammar.Grammar
	FIRSTSets map[string]sets.Set[string]
	table map[string]map[string]ParserAction
}

func NewLR1ParsingTable(grammar *lr1grammar.Grammar) *ParsingTable {
	FIRSTSets := firstfollow.ComputeFIRSTSets(grammar)
	table := make(map[string]map[string]ParserAction)

	parsingTable := ParsingTable{
		grammar,
		FIRSTSets,
		table,
	}

	return &parsingTable
}

func (pt *ParsingTable) CLOSURE(item lr1item.LR1Item) {}

func (pt *ParsingTable) GOTO(item lr1item.LR1Item) {}


// Computes the lookahead set for a given context: sequence of symbols following a
// non-terminal to the RHS of the parsing progress (the dot).
func (pt *ParsingTable) Lookahead(context []string, currLookahead sets.Set[string]) sets.Set[string] {
	if len(context) == 0 {
		return currLookahead.Clone()
	} else if context[0] == symbols.Epsilon {
		return currLookahead.Clone()
	} else {
		lookaheadSet := sets.NewEmptySet[string]()
		for _, symbol := range context {
			if pt.Grammar.Terminals.Has(symbol) {
				// symbol is a terminal: this is the only possible lookahead
				lookaheadSet.Add(symbol)
				break
			} else {
				// symbol is a non-terminal: FIRST(symbol) is in lookaheadSet
				symbolFIRSTSet := pt.FIRSTSets[symbol]
				lookaheadSet = lookaheadSet.Union(symbolFIRSTSet)
				if !pt.Grammar.DerivesEpsilon(symbol) {
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