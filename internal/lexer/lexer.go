package lexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"interpreters/internal/symbols"
	"interpreters/utilities/arrays"
	"io"
	"os"
	"regexp"
	"sort"
)

// ----- DEFAULT PATTERNS -----
var whitespacePattern regexp.Regexp = *regexp.MustCompile(`^\s`)
var newlinePattern regexp.Regexp = *regexp.MustCompile(`^(\n|\r)`)

type LexerConfigJson struct {
	KeywordTokens TokenConfigJsonArr `json:"keywordTokens"`
	SymbolTokens TokenConfigJsonArr `json:"symbolTokens"`
	GenericTokens TokenConfigJsonArr `json:"genericTokens"`
}

type Lexer struct {
	keywordTokens []*TokenConfig
	symbolTokens []*TokenConfig
	genericTokens []*TokenConfig
	line uint
	col uint
}

// callback function to retrieve a `TokenConfig` from `LexerConfigJson` fields.
func getTokenConfig(tokenConfigJson TokenConfigJson) *TokenConfig {
	return tokenConfigJson.CreateTokenConfig()
}

func CreateLexer(config LexerConfigJson) *Lexer {
	// sort keyword and symbol patterns by decreasing length to ensure maximal 
	// length tokens are matched first
	sort.Sort(sort.Reverse(config.KeywordTokens))
	sort.Sort(sort.Reverse(config.SymbolTokens))

	keywords := arrays.Map(config.KeywordTokens, getTokenConfig)
	symbols := arrays.Map(config.SymbolTokens, getTokenConfig)
	generics := arrays.Map(config.GenericTokens, getTokenConfig)

	return &Lexer{
		keywords,
		symbols,
		generics,
		0,
		0,
	}
}

func CreateLexerFromJsonConfig(path string) (*Lexer, error) {
	configFile, err := os.Open(path)
    if err != nil {
        return nil, errors.New(`Error opening config file: ` + err.Error())
    }
    defer configFile.Close()

	bytes, err := io.ReadAll(configFile)
    if err != nil {
        return nil, errors.New(`Error reading config file: ` + err.Error())
    }

	var data LexerConfigJson
	err = json.Unmarshal(bytes, &data)
    if err != nil {
		return nil, errors.New(`Error unmarshalling config file: ` + err.Error())
    }

	return CreateLexer(data), nil
}

func (lex *Lexer) matchTokenGroup(tokenConfigs []*TokenConfig, inputStream string) *Token {
	for _, tokenConfig := range tokenConfigs {
		token := tokenConfig.Match(inputStream)
		if token != nil {
			token.Col = lex.col + 1
			token.Line = lex.line + 1
			lex.col += uint(len(token.Value))
			return token
		}
	}

	return nil
}

func (lex *Lexer) Tokenize(input string) *[]*Token {
	result := make([]*Token, 0)
	processed := 0
	lex.line = 0
	lex.col = 0
	tokenGroups := [][]*TokenConfig{
		lex.symbolTokens,
		lex.keywordTokens,
		lex.genericTokens,
	}

	for processed < len(input) {
		currInput := input[processed:]

		// check for newlines
		if (newlinePattern.MatchString(currInput)) {
			processed++
			lex.line++
			lex.col = 0
			continue
		}
		// check for whitespaces
		if (whitespacePattern.MatchString(currInput)) {
			processed++
			lex.col++
			continue
		}

		// match symbols, keywords, and then generics - in that order
		var token *Token
		for _, tokenGroup := range tokenGroups {
			token = lex.matchTokenGroup(tokenGroup, currInput)
			if (token != nil) {
				break
			}
		}

		if (token == nil) {
			message := fmt.Sprintf("Unrecognized symbol at %d:%d", lex.line, lex.col)
			panic(message)
		} 

		result = append(result, token)
		processed += len(token.Value)
	}

	result = append(result, &Token{
		symbols.Epsilon,
		symbols.Epsilon,
		0,
		0,
	})

	return &result
}


