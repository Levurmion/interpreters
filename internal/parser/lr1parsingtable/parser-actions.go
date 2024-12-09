package lr1parsingtable

import "fmt"

type ParserActionVerbs string

const (
	SHIFT 	ParserActionVerbs = "shift"
	REDUCE 	ParserActionVerbs = "reduce"
	ACCEPT 	ParserActionVerbs = "accept"
	GOTO 	ParserActionVerbs = "goto"
	ERROR 	ParserActionVerbs = "error"
)

type ParserAction interface {
	ActionVerb() 	ParserActionVerbs
	NextState()		int
	ReduceByRule()	int
	Message()		string
}


// ----- SHIFT ACTION -----
type ShiftAction struct {
	nextState int
	symbol string
}
func (a *ShiftAction) ActionVerb() ParserActionVerbs { return SHIFT }
func (a *ShiftAction) NextState() int { return a.nextState }
func (a *ShiftAction) ReduceByRule() int { return -1 }
func (a *ShiftAction) Message() string { return fmt.Sprintf("shift symbol: %s", a.symbol) }

// ----- REDUCE ACTION ------
type ReduceAction struct {
	ruleId int
}
func (a *ReduceAction) ActionVerb() ParserActionVerbs { return REDUCE }
func (a *ReduceAction) NextState() int { return -1 }
func (a *ReduceAction) ReduceByRule() int { return a.ruleId }
func (a *ReduceAction) Message() string { return fmt.Sprintf("reduce by rule id: %d", a.ruleId) }

// ----- GOTO ACTION -----
type GotoAction struct {
	nextState int
}
func (a *GotoAction) ActionVerb() ParserActionVerbs { return GOTO }
func (a *GotoAction) NextState() int { return a.nextState }
func (a *GotoAction) ReduceByRule() int { return -1 }
func (a *GotoAction) Message() string { return fmt.Sprintf("go to state: %d", a.nextState) }

// ----- ACCEPT ACTION -----
type AcceptAction struct {}
func (a *AcceptAction) ActionVerb() ParserActionVerbs { return ACCEPT }
func (a *AcceptAction) NextState() int { return -1 }
func (a *AcceptAction) ReduceByRule() int { return -1 }
func (a *AcceptAction) Message() string { return "accept" }

// ----- ERROR ACTION -----
type ErrorAction struct {
	errorMessage string
}
func (a *ErrorAction) ActionVerb() ParserActionVerbs { return ERROR }
func (a *ErrorAction) NextState() int { return -1 }
func (a *ErrorAction) ReduceByRule() int { return -1 }
func (a *ErrorAction) Message() string { return a.errorMessage }
