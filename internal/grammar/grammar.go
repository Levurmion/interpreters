package grammar

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

	// load all token types into terminals set
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

	// load all non-terminals into nonTerminals set
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

