package lexer

type LexerConfigJson struct {
	KeywordTokens []TokenConfigJson `json:"keywordTokens"`
	SymbolTokens []TokenConfigJson `json:"symbolTokens"`
	GenericTokens []TokenConfigJson `json:"genericTokens"`
}

type Lexer struct {
	keywordTokens []*TokenConfig
	symbolTokens []*TokenConfig
	genericTokens []*TokenConfig

	line uint
	col uint
}

func CreateLexer (config LexerConfigJson) *Lexer {
	

}