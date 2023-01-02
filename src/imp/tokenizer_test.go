package imp

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutomatedAll(t *testing.T) {
	t.Log("Test all terminal tokens")
	for _, key := range terminalTokensPriority {
		assertTokensMatch(t, key, terminal(key))
		assertTokensMatch(t, " "+key+" ", terminal(key))
		testedSource := key + " " + key
		t.Log("Tested source: ", testedSource)
		assertTokensMatch(t, testedSource, terminal(key), terminal(key))
	}

	// TODO: larger token pool, more breadth
	t.Log("Randomize and test two terminals in a row")
	detRand := rand.New(rand.NewSource(123))
	for i := 0; i < 100; i++ {
		tokenClone := make([]string, len(terminalTokensPriority))
		copy(tokenClone, terminalTokensPriority)
		detRand.Shuffle(len(tokenClone), func(i, j int) {
			tokenClone[i], tokenClone[j] = tokenClone[j], tokenClone[i]
		})
		for index, key := range terminalTokensPriority {
			testedSource := key + " " + tokenClone[index]
			t.Log("Tested source: " + testedSource)
			assertTokensMatch(t, testedSource,
				terminal(key), terminal(tokenClone[index]))
		}
	}
}

func TestMisc(t *testing.T) {
	assertTokensMatch(t, "= =", terminal(ASSIGNMENT), terminal(ASSIGNMENT))
	assertTokensMatch(t, "== ==", terminal(EQUALS), terminal(EQUALS))
}

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

	assertTokensEmpty(t, "")
	assertTokensMatch(t, "whi", variable("whi"))

	assertTokensMatch(t, "whilewhi=", variable("whilewhi"), terminal(ASSIGNMENT))
}

func TestBooleans(t *testing.T) {
	t.Log("Test boolean literals")
	assertTokensMatch(t, "true", booleanToken(true))
	assertTokensMatch(t, "false", booleanToken(false))
	assertTokensMatch(t, "tru", variable("tru"))
}

func TestNumbers(t *testing.T) {
	t.Log("Test number literal tokenization")
	assertTokensMatch(t, "123", integer(123))
	assertTokensMatch(t, "-123", integer(-123))
	assertProducesErrorAtToken(t, "- 123")
	assertProducesErrorAtToken(t, "12,3", integer(12))
}

func TestNumberExpressions(t *testing.T) {
	t.Log("Test numbers within expressions")
	assertTokensMatch(t, "-12+-5", integer(-12), terminal(ADD), integer(-5))
	assertTokensMatch(t, "-12-5", integer(-12), integer(-5))
	assertTokensMatch(t, "(-12 + (-5))",
		terminal(OPEN_EXPRESSION_GROUPING), integer(-12), terminal(ADD),
		terminal(OPEN_EXPRESSION_GROUPING), integer(-5),
		terminal(CLOSE_EXPRESSION_GROUPING), terminal(CLOSE_EXPRESSION_GROUPING))
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

	assertTokensMatch(t, "trueThat", variable("trueThat"))

	//assertTokensEmpty(t, "123name")
	//TODO: test illegal variable names here or in parser?
	//go's own parser seems to roll with valid Integer,VariableName tokens,
	//and invalidates them during parsing
	//so maybe leave it as is
	//if tokens without delimiters are permitted (e.g. "(123)" instead of "( 123 )"),
	//then recognizing fused tokens in lexer is necessary and inevitable?

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

	t.Log("Parenthesis and literals")
	assertTokensMatch(t, "(test=(-123+-0)&&",
		terminal(OPEN_EXPRESSION_GROUPING),
		variable("test"), terminal(ASSIGNMENT),
		terminal(OPEN_EXPRESSION_GROUPING),
		integer(-123),
		terminal(ADD),
		Token{
			tokenType:    IntegerValue,
			token:        "-0",
			integerValue: 0,
		},
		terminal(CLOSE_EXPRESSION_GROUPING),
		terminal(AND))
}

func assertProducesErrorAtToken(t *testing.T, sourceCode string, expectedTokens ...Token) {
	tokens, err := tokenize(sourceCode)
	assert.Error(t, err)
	assert.Equal(t, (TokenizerResultData)(expectedTokens), tokens)
}

func assertTokensMatch(t *testing.T, sourceCode string, expectedTokens ...Token) {
	assertTokensResultMatch(t, sourceCode, tokens(expectedTokens...))
}

func assertTokensEmpty(t *testing.T, sourceCode string) {
	assertTokensResultMatch(t, sourceCode, TokenizerResultData{})
}

func assertTokensResultMatch(t *testing.T, sourceCode string, expectedTokens TokenizerResultData) {
	tokenList, err := tokenize(sourceCode)
	assert.NoError(t, err)
	assert.Equal(t, expectedTokens, tokenList)
}
