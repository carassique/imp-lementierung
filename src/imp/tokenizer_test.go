package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func tokens(tokens ...Token) TokenizerResultData {
	return (TokenizerResultData)(tokens)
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

func TestVariable(t *testing.T) {
	t.Log("Variable name.")
	variableToken := variable("test")
	assertTokensMatch(t, "test", variableToken)
	assertTokensMatch(t, " test", variableToken)
	assertTokensMatch(t, "test ", variableToken)
	assertTokensMatch(t, "testWithNumeric123", variable("testWithNumeric123"))

	//TODO: test illegal variable names here or in parser?
	//split variable names "vari ablename"
	//illegal format "123name"
}

func TestExpressionGrouping(t *testing.T) {
	t.Log("Empty parenthesis with different spacing.")
	emptyParenthesisTokens := TokenizerResultData{
		openExpressionGrouping(),
		closeExpressionGrouping(),
	}
	assertTokensResultMatch(t, "()", emptyParenthesisTokens)
	assertTokensResultMatch(t, "( )", emptyParenthesisTokens)
	assertTokensResultMatch(t, " ()", emptyParenthesisTokens)
	assertTokensResultMatch(t, "() ", emptyParenthesisTokens)
	assertTokensResultMatch(t, " () ", emptyParenthesisTokens)
	assertTokensResultMatch(t, " (  ) ", emptyParenthesisTokens)

	t.Log("Parenthesis with variable names.")

	variableNameInParenthesisTokens := TokenizerResultData{
		openExpressionGrouping(),
		variable("test"),
		closeExpressionGrouping(),
	}

	assertTokensResultMatch(t, "(test)", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, "( test)", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, "(test )", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, "( test )", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, " (test)", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, "(test) ", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, " ( test ) ", variableNameInParenthesisTokens)
	assertTokensResultMatch(t, "(test  )", variableNameInParenthesisTokens)
}

func assertTokensMatch(t *testing.T, sourceCode string, expectedTokens ...Token) {
	assertTokensResultMatch(t, sourceCode, tokens(expectedTokens...))
}

func assertTokensResultMatch(t *testing.T, sourceCode string, expectedTokens TokenizerResultData) {
	tokenList := tokenize(sourceCode, terminalTokens)
	assert.Equal(t, expectedTokens, tokenList)
}
