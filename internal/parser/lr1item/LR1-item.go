package lr1item

import (
	"errors"
	"fmt"
	"interpreters/internal/symbols"
	"interpreters/utilities/sets"
	"strings"
)

type LR1Item struct {
	LHS string
	RHS []string
	LookaheadSet sets.Set[string]
	DotPosition int
	OriginalRHS []string
}


// Generates a new `LR1Item`.
func NewLR1Item(LHS string, RHS []string, dotPosition int, lookaheadSet sets.Set[string]) (*LR1Item, error) {
	RHSLength := len(RHS)

	// ensure dotPosition is valid
	if dotPosition > RHSLength || dotPosition < 0 {
		message := fmt.Sprintf(
			`dotPostion (%d) needs to be between 0 and %d (length of RHS).`, 
			dotPosition, 
			RHSLength,
		)
		return nil, errors.New(message)
	}

	// create new copies of each array for manipulation
	RHSLeft := RHS
	RHSRight := RHS
	RHSWithDot := []string{}

	RHSLeft = RHSLeft[0:dotPosition]
	RHSRight = RHSRight[dotPosition:RHSLength]
	RHSWithDot = append(RHSWithDot, RHSLeft...)
	RHSWithDot = append(RHSWithDot, symbols.Dot)
	RHSWithDot = append(RHSWithDot, RHSRight...)

	return &LR1Item{
		LHS,
		RHSWithDot,
		lookaheadSet,
		dotPosition,
		RHS,
	}, nil
}

// Get the symbol immediately right of the `dot`. Returns `Epsilon` if the RHS 
// is complete (`dot` at the end of the production).
func (item *LR1Item) GetNextSymbol() string {
	if item.DotPosition == len(item.RHS) - 1 {
		return symbols.Epsilon
	} else {
		return item.RHS[item.DotPosition + 1]
	}
}

func (item *LR1Item) ProductionIsComplete() bool {
	return item.GetNextSymbol() == symbols.Epsilon
}

// Retrieves the context sequence for the next symbol in the production rule.
func (item *LR1Item) GetContextForNextSymbol() []string {
	if item.GetNextSymbol() == symbols.Epsilon {
		return []string{}
	} else {
		return item.RHS[item.DotPosition + 2:]
	}
}

// Returns a new `LR1Item` with the `dot` advanced by one position in the RHS.
// Also returns the context for the symbol processed by the parser when advancing.
func (item *LR1Item) AdvanceDot() (*LR1Item, error) {
	if item.GetNextSymbol() == symbols.Epsilon {
		return nil, fmt.Errorf(`LR1Item production is complete: cannot advance dot`)
	} else {
		return NewLR1Item(
			item.LHS, 
			item.OriginalRHS, 
			item.DotPosition+1,
			item.LookaheadSet,
		)
	}
}

// Get the `string` representation of this item.
func (item *LR1Item) GetName() string {
	name := item.LHS + "->" + strings.Join(item.RHS, "")
	return name
}