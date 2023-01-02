package imp

import "strconv"

func tokens(tokens ...Token) TokenizerResultData {
	return (TokenizerResultData)(tokens)
}

func noMatch(word string) Token {
	return Token{
		tokenType: Error,
		token:     word,
	}
}

func terminal(value string) Token {
	return Token{
		tokenType: Terminal,
		token:     value,
	}
}

func booleanToken(value bool) Token {
	token := BOOLEAN_FALSE
	if value {
		token = BOOLEAN_TRUE
	}
	return Token{
		tokenType:    BooleanValue,
		token:        token,
		booleanValue: value,
	}
}

func integer(value int) Token {
	return Token{
		tokenType:    IntegerValue,
		token:        strconv.FormatInt((int64)(value), 10),
		integerValue: value,
	}
}

func variable(name string) Token {
	return Token{
		tokenType: VariableName,
		token:     name,
	}
}

func popen() Token {
	return openExpressionGrouping()
}

func pclose() Token {
	return closeExpressionGrouping()
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
