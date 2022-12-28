package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminals(t *testing.T) {
	t.Log("Test terminal tokens")
	assertTokensMatch(t, "= ", terminal(ASSIGNMENT))
	assertTokensMatch(t, "== ", terminal(EQUALS))
	assertTokensMatch(t, "while ", terminal(WHILE))
	assertTokensMatch(t, "while = ==", terminal(WHILE), terminal(ASSIGNMENT), terminal(EQUALS))
	assertTokensMatch(t, "===", terminal(EQUALS), terminal(ASSIGNMENT)) //TODO: consider
	assertTokensMatch(t, "while:=print==if<=", terminal(WHILE), terminal(DECLARATION),
		terminal(PRINT), terminal(EQUALS), terminal(IF),
		terminal(LESS_THAN), terminal(ASSIGNMENT))

	assertTokensResultMatch(t, "", TokenizerResultData{})
	assertTokensResultMatch(t, "whi", TokenizerResultData{})

	assertTokensMatch(t, "whilewhi=", variable("whilewhi"), terminal(ASSIGNMENT))
}

func TestNumbers(t *testing.T) {
	t.Log("Test number literal tokenization")
	assertTokensMatch(t, "123", integer(123))
	assertTokensMatch(t, "-123", integer(-123))
	assertTokensMatch(t, "- 123", integer(123)) //TODO: consider error

}

func TestStrangeVariableNames(t *testing.T) {
	t.Log("Test strange variable names")
	assertTokensMatch(t, "something123", variable("something123"))
	assertTokensMatch(t, "somethingwhile", variable("somethingwhile"))
	assertTokensMatch(t, "whileSomething", variable("whileSomething"))
	assertTokensMatch(t, "while something", terminal(WHILE), variable("something"))
	assertTokensMatch(t, "while\nsomething", terminal(WHILE), variable("something"))
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
