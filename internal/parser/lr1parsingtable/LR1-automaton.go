package lr1parsingtable

import (
	"errors"
	"interpreters/internal/parser/firstfollow"
	"interpreters/internal/parser/lr1closureset"
	"interpreters/internal/parser/lr1grammar"
	"interpreters/internal/parser/lr1item"
	"interpreters/internal/symbols"
	"interpreters/utilities/arrays"
	"interpreters/utilities/sets"
)

type LR1Automaton struct {
	grammar *lr1grammar.Grammar
	FIRSTSets map[string]sets.Set[string]
	States map[int]ParserState
}

func NewLR1Automaton(grammar *lr1grammar.Grammar) (*LR1Automaton, error) {
	FIRSTSets := firstfollow.ComputeFIRSTSets(grammar)
	states := make(map[int]ParserState)
	automaton := LR1Automaton{grammar, FIRSTSets, states}

	// verify that grammar is properly augmented
	augmentedProduction := grammar.GetProductionsOfNonTerminal(symbols.AugmentedStart)
	if (len(augmentedProduction) != 1) {
		return nil, errors.New("LR1Automaton requires an augmented grammar")
	}

	// initialize I_0 with the augmented start production rule
	originalStartSymbol := augmentedProduction[0].Production[0]
	firstItem, err := lr1item.NewLR1Item(
		symbols.AugmentedStart, 
		[]string{originalStartSymbol},
		0,
		sets.NewSet(symbols.EOF),
	)
	if err != nil {
		return nil, err
	}

	I0ClosureSet := lr1closureset.NewLR1ClosureSet(firstItem)
	I0NextStates := make(map[string]int)
	I0 := ParserState{
		I0ClosureSet,
		I0NextStates,
	}
	automaton.States[0] = I0

	return &automaton, nil
}

func (automaton *LR1Automaton) Lookahead(context []string, currLookahead sets.Set[string]) sets.Set[string] {
	if len(context) == 0 {
		return currLookahead.Clone()
	} else if context[0] == symbols.Epsilon {
		return currLookahead.Clone()
	} else {
		lookaheadSet := sets.NewEmptySet[string]()
		for _, symbol := range context {
			if automaton.grammar.Terminals.Has(symbol) {
				// symbol is a terminal: this is the only possible lookahead
				lookaheadSet.Add(symbol)
				break
			} else {
				// symbol is a non-terminal: FIRST(symbol) is in lookaheadSet
				symbolFIRSTSet := automaton.FIRSTSets[symbol]
				lookaheadSet = lookaheadSet.Union(symbolFIRSTSet)
				if !automaton.grammar.DerivesEpsilon(symbol) {
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

func (automaton *LR1Automaton) CLOSURE(item *lr1item.LR1Item) *lr1closureset.LR1ClosureSet {
	result := lr1closureset.NewEmptyLR1ClosureSet()
	result.Add(item)

	changed := true
	for changed {
		nextSymbols := result.GetTransitionSymbols()
		nextNonTerminalSymbols := arrays.Filter(nextSymbols.GetItems(), func (symbol string) bool {
			return automaton.grammar.NonTerminals.Has(symbol)
		})
	}
}