package lexer

import (
	"regexp"
)

type Token struct {
	Type string
	Value string
	Line int
	Col int
}

type TokenConfigJson struct {
	Type string `json:"type"`
	Pattern string `json:"pattern"`
}

type TokenConfig struct {
	Type string
	Pattern *regexp.Regexp
}

func (json *TokenConfigJson) CreateTokenConfig () *TokenConfig {
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

/*
Matches a `TokenConfig.Pattern` to the start of an input string. Returns `nil` if
not match was found.
*/
func (tokenConfig *TokenConfig) Match (input string) *Token {
	match := tokenConfig.Pattern.FindStringSubmatch(input)
	if (match == nil) {
		return nil
	} else {
		return &Token{
			tokenConfig.Type,
			match[0],
			-1,
			-1,
		}
	}
}