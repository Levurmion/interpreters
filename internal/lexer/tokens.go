package lexer

import (
	"regexp"
)

type Token struct {
	Type string
	Value string
	Line uint
	Col uint
}

type TokenConfigJson struct {
	Type string `json:"type"`
	Pattern string `json:"pattern"`
}

// implementations for sort.Interface
type TokenConfigJsonArr []TokenConfigJson
func (arr TokenConfigJsonArr) Len() int 				{ return len(arr) }
func (arr TokenConfigJsonArr) Swap(i int, j int) 		{ arr[i], arr[j] = arr[j], arr[i] }
func (arr TokenConfigJsonArr) Less(i int, j int) bool 	{ return len(arr[i].Pattern) < len(arr[j].Pattern) }

func (json *TokenConfigJson) CreateTokenConfig() *TokenConfig {
	if (len(json.Pattern) <= 0) {
		return nil
	} else if (json.Pattern[0] == '^') {
		regex := regexp.MustCompile(json.Pattern)
		return &TokenConfig{
			json.Type,
			regex,
		}
	} else {
		matchStartPattern := "^" + json.Pattern
		regex := regexp.MustCompile(matchStartPattern)
		return &TokenConfig{
			json.Type,
			regex,
		}
	}
}

type TokenConfig struct {
	Type string
	Pattern *regexp.Regexp
}

/*
Matches a `TokenConfig.Pattern` to the start of an input string. Returns `nil` if
not match was found.
*/
func (tokenConfig *TokenConfig) Match(input string) *Token {
	match := tokenConfig.Pattern.FindStringSubmatch(input)
	if (match == nil) {
		return nil
	} else {
		return &Token{
			tokenConfig.Type,
			match[0],
			0,
			0,
		}
	}
}