package lexer_test

import (
	"interpreters/internal/lexer"
	"testing"

	"github.com/go-test/deep"
)

func TestLexer(t *testing.T) {
	Lexer, err := lexer.CreateLexerFromJsonConfig("./token-config.json")
	if (err != nil) {
		t.Error("Failed to initialize lexer: ", err.Error())
	}

	var testCases = []struct{
		name string
		input string
		output []*lexer.Token
	}{
		{
			"Lexer can tokenize a simple object.",
			`{ "prop_a": true }`, 
			[]*lexer.Token{
				{"{", "{", 1, 1}, 
				{"str_lit", `"prop_a"`, 1, 3},
				{":", ":", 1, 11},
				{"true", "true", 1, 13},
				{"}", "}", 1, 18},
			},
		},
		{
			"Lexer can tokenize a multiline object.",
			`{ 
				"prop_a": true
			 }`, 
			[]*lexer.Token{
				{"{", "{", 1, 1}, 
				{"str_lit", `"prop_a"`, 2, 5},
				{":", ":", 2, 13},
				{"true", "true", 2, 15},
				{"}", "}", 3, 5},
			},
		},
		{
			"Lexer can tokenize a single value",
			`true`, 
			[]*lexer.Token{
				{"true", "true", 1, 1},
			},
		},
		{
			"Lexer can tokenize arrays",
			`[ true, false, "string", -2.45 ]`, 
			[]*lexer.Token{
				{"[", "[", 1, 1},
				{"true", "true", 1, 3},
				{",", ",", 1, 7},
				{"false", "false", 1, 9},
				{",", ",", 1, 14},
				{"str_lit", `"string"`, 1, 16},
				{",", ",", 1, 24},
				{"num_lit", "-2.45", 1, 26},
				{"]", "]", 1, 32},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := Lexer.Tokenize(tc.input)
			if diff := deep.Equal(*output, tc.output); diff != nil {
				t.Error(diff)
			}
		})
	}
}