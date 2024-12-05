package main

import (
	"encoding/json"
	"fmt"
	"interpreters/internal/lexer"
	"interpreters/utilities/arrays"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./token-config.json")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

	bytes, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

	var data lexer.LexerConfigJson
	err = json.Unmarshal(bytes, &data)
    if err != nil {
        fmt.Println("Error unmarshaling JSON:", err)
        return
    }

	Lexer := lexer.CreateLexer(data)
	input := `{ "prop_a": ["elem", true, false], "prop_b": null }`

	tokens := Lexer.Tokenize(input) 

    fmt.Println(arrays.Map(*tokens, func(token *lexer.Token) lexer.Token { 
		return *token 
	}))
}