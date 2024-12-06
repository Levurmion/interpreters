package grammar

import (
	"encoding/json"
	"errors"
	"fmt"
	"interpreters/internal/lexer"
	"interpreters/internal/symbols"
	"interpreters/utilities/arrays"
	"interpreters/utilities/sets"
	"io"
	"os"
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
	for _, nonTerminal := range nonTerminals.GetItems() {
		if (nonTerminal == symbols.AugmentedStart) {
			panic(fmt.Sprintf("Grammar cannot use the reserved symbol: %s", nonTerminal))
		}
	}

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

func NewGrammarFromJsonConfig(path string) (*Grammar, error) {
	configFile, err := os.Open(path)
    if err != nil {
        return nil, errors.New(`Error opening config file: ` + err.Error())
    }
    defer configFile.Close()

	bytes, err := io.ReadAll(configFile)
    if err != nil {
        return nil, errors.New(`Error reading config file: ` + err.Error())
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
	productionRules := g.GetProductionsDerivingSymbol(symbol)
	return arrays.FindFirstIdx(productionRules, func(productionRule ProductionRule) bool {
		return productionRule.Production[0] == symbols.Epsilon
	}) >= 0
}