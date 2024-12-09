package firstfollow

import (
	"interpreters/internal/parser/lr1grammar"
	"interpreters/utilities/sets"
)

func ComputeFIRSTandFOLLOW(grammar *lr1grammar.Grammar) (map[string]sets.Set[string], map[string]sets.Set[string]) {
	FIRSTSets := ComputeFIRSTSets(grammar)
	FOLLOWSets := ComputeFOLLOWSets(grammar, FIRSTSets)
	return FIRSTSets, FOLLOWSets
}