package lr1parsingtable

import "interpreters/internal/parser/lr1closureset"

type ParserState struct {
	CLOSURESet *lr1closureset.LR1ClosureSet
	NextStates map[string]int
}

