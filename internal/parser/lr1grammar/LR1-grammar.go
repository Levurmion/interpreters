package lr1grammar

import (
	"encoding/json"
	"errors"
	"fmt"
	"interpreters/internal/lexer"
	"interpreters/internal/symbols"
	"interpreters/utilities/arrays"
	"interpreters/utilities/files"
	"interpreters/utilities/sets"
)

type GrammarConfigJson struct {
	Terminals 		lexer.LexerConfigJson	`json:"terminals"`
	NonTerminals 	map[string][][]string	`json:"nonTerminals"`
	StartSymbol		string					`json:"startSymbol"`
}

type ProductionRule struct {
	NonTerminal 	string
	Production 		[]string
}

type Grammar struct {
	Terminals 					sets.Set[string]
	NonTerminals 				sets.Set[string]
	AllSymbols					sets.Set[string]
	StartSymbol					string
	ProductionRules				map[uint]ProductionRule

	productionRulesIdx 			map[string]*[]uint
	productionRulesInvertedIdx 	map[string]*[]uint
}

func NewAugmentedGrammar(config GrammarConfigJson) *Grammar {
	// add symbols.AugmentedStart
	config.NonTerminals[symbols.AugmentedStart] = [][]string{{config.StartSymbol}}
	config.StartSymbol = symbols.AugmentedStart
	return NewGrammar(config)
}

func NewGrammar(config GrammarConfigJson) *Grammar {
	terminals := sets.NewEmptySet[string]()
	nonTerminals := sets.NewEmptySet[string]()
	enumeratedProductionRules := make(map[uint]ProductionRule)

	// maps non-terminal -> production rule of non-terminal
	enumeratedProductionRulesIdx := make(map[string]*[]uint)
	
	// maps any symbol -> production rule containing symbol
	enumeratedProductionRulesInvertedIdx := make(map[string]*[]uint)

	// load all token types uinto terminals set
	for _, token := range config.Terminals.SymbolTokens {
		terminals.Add(token.Type)
	}
	for _, token := range config.Terminals.KeywordTokens {
		terminals.Add(token.Type)
	}
	for _, token := range config.Terminals.GenericTokens {
		terminals.Add(token.Type)
	}

	// verify that no symbols are reserved keywords
	for _, terminal := range terminals.GetItems() {
		if (terminal == symbols.EOF || terminal == symbols.Dot) {
			panic(fmt.Sprintf("Grammar cannot use the reserved symbol: %s", terminal))
		}
	}

	// also add Epsilon and EOF as possible terminals
	terminals.Add(symbols.Epsilon)
	terminals.Add(symbols.EOF)

	// load all non-terminals uinto nonTerminals set
	for nonTerminal := range config.NonTerminals {
		nonTerminals.Add(nonTerminal)
	}

	// enumerate production rules and create forward index
	var i uint
	for nonTerminal, productionRules := range config.NonTerminals {
		index := make([]uint, len(productionRules))
		enumeratedProductionRulesIdx[nonTerminal] = &index
		for j, productionRule := range productionRules {
			enumeratedProductionRules[i] = ProductionRule{
				nonTerminal,
				productionRule,
			}
			index[j] = i
			i++
		}
	}

	// create inverted index
	for ruleId, productionRule := range enumeratedProductionRules {
		for _, symbol := range productionRule.Production {
			invertedIndex, exists := enumeratedProductionRulesInvertedIdx[symbol]
			if exists {
				*invertedIndex = append(*invertedIndex, ruleId)
			} else {
				enumeratedProductionRulesInvertedIdx[symbol] = &[]uint{ruleId}
			}
		}
	}

	return &Grammar{
		terminals,
		nonTerminals,
		terminals.Union(nonTerminals),
		config.StartSymbol,
		enumeratedProductionRules,
		enumeratedProductionRulesIdx,
		enumeratedProductionRulesInvertedIdx,
	}
}

func NewAugmentedGrammarFromJsonConfig(path string) (*Grammar, error) {
	bytes, err := files.OpenFileToByteStream(path)
	if err != nil {
		return nil, err
	}

	var data GrammarConfigJson
	err = json.Unmarshal(bytes, &data)
    if err != nil {
		return nil, errors.New(`Error unmarshalling config file: ` + err.Error())
    }

	return NewAugmentedGrammar(data), nil
}

func NewGrammarFromJsonConfig(path string) (*Grammar, error) {
	bytes, err := files.OpenFileToByteStream(path)
	if err != nil {
		return nil, err
	}

	var data GrammarConfigJson
	err = json.Unmarshal(bytes, &data)
    if err != nil {
		return nil, errors.New(`Error unmarshalling config file: ` + err.Error())
    }

	return NewGrammar(data), nil
}

// ----- GRAMMAR METHODS -----

func (g *Grammar) GetProductionsOfNonTerminal(nonTerminal string) []ProductionRule {
	pIndex, exists := g.productionRulesIdx[nonTerminal]
	if (exists) {
		index := *pIndex
		return arrays.Map(index, func (idx uint) ProductionRule {
			return g.ProductionRules[idx]
		})
	} else {
		return []ProductionRule{}
	}
}

func (g *Grammar) GetProductionsDerivingSymbol(symbol string) []ProductionRule {
	pIndex, exists := g.productionRulesInvertedIdx[symbol]
	if (exists) {
		index := *pIndex
		return arrays.Map(index, func (idx uint) ProductionRule {
			return g.ProductionRules[idx]
		})
	} else {
		return []ProductionRule{}
	}
}

// Find the ID of a production rule as registered in the `Grammar`. Returns -1 if
// the queried production rule does not exist.
func (g *Grammar) GetProductionId(LHS string, RHS []string) (int, error) {
	var possibleIds sets.Set[uint]

	// get all production IDs for the non-terminal
	pForwardIndex, exists := g.productionRulesIdx[LHS]
	if (!exists) {
		return -1, fmt.Errorf(`Non-terminal: %s does not exist in the specified grammar.`, LHS)
	} else {
		forwardIndex := *pForwardIndex
		possibleIds = sets.NewSet(forwardIndex...)
	}

	// get all possible production IDs for RHS symbols
	for _, symbol :=  range RHS {
		pInvertedIndex, exists := g.productionRulesInvertedIdx[symbol]
		if (!exists) {
			return -1, fmt.Errorf(`symbol: %s does not exist in the specified grammar`, symbol)
		} else {
			invertedIndex := *pInvertedIndex
			invertedIndexSet := sets.NewSet(invertedIndex...)
			possibleIds = possibleIds.Intersection(invertedIndexSet)
		}
	}

	// there should be no ambiguity in resolving the ID of a production
	if possibleIds.Size() > 1 {
		return -1, errors.New(`production rule maps to multiple ambiguous IDs`)
	} else if possibleIds.Size() == 0 {
		return -1, errors.New(`production rule not found in the specified grammar`)
	} else {
		return int(possibleIds.GetItems()[0]), nil
	}
}

func (g *Grammar) DerivesEpsilon(symbol string) bool {
	// If symbol is not a non-terminal, it never derives epsilon
	if !g.NonTerminals.Has(symbol) {
		return false
	}

	productionRules := g.GetProductionsOfNonTerminal(symbol)
	for _, productionRule := range productionRules {
		// Check if this production is a single epsilon
		if len(productionRule.Production) == 1 && productionRule.Production[0] == symbols.Epsilon {
			return true
		}
	}

	return false
}




