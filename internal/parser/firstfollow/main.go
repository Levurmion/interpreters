package firstfollow

import (
	"interpreters/internal/grammar"
	"interpreters/utilities/sets"
)

func ComputeFIRSTandFOLLOW(grammar *grammar.Grammar) (map[string]sets.Set[string], map[string]sets.Set[string]) {
	FIRSTSets := ComputeFIRSTSets(grammar)
	FOLLOWSets := ComputeFOLLOWSets(grammar, FIRSTSets)
	return FIRSTSets, FOLLOWSets
}