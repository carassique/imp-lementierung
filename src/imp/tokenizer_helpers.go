package imp

func tokens(tokens ...Token) TokenizerResultData {
	return (TokenizerResultData)(tokens)
}

func noMatch() Token {
	return Token{}
}

func terminal(value string) Token {
	return Token{
		tokenType: Terminal,
		token:     value,
	}
}

func variable(name string) Token {
	return Token{
		tokenType: VariableName,
		token:     name,
	}
}

func openExpressionGrouping() Token {
	return Token{
		tokenType: Terminal,
		token:     OPEN_EXPRESSION_GROUPING,
	}
}

func closeExpressionGrouping() Token {
	return Token{
		tokenType: Terminal,
		token:     CLOSE_EXPRESSION_GROUPING,
	}
}
