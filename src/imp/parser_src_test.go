package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSources(t *testing.T) {
	assertSourceMatches(t,
		whileStatement(
			and(
				lessThan(
					number(0),
					variableExpression("iterator"),
				),
				lessThan(
					number(5),
					variableExpression("iterator"),
				),
			),
			printStatement(variableExpression("iterator")),
		),
		// (0 < ((iterator && 5) < iterator)
		// (0 < iterator) && (5 < iterator)
		"while 0 < iterator && 5 < iterator { print iterator }",
	)

	// a + b * c
	assertSourceMatches(t,
		printStatement(plus(
			variableExpression("a"),
			mult(variableExpression("b"), variableExpression("c")),
		)),
		"print a + b * c",
	)
}

func assertSourceMatches(t *testing.T, expectedAst Stmt, source string) {
	tokens, err := tokenize(source)
	assert.NoError(t, err)
	assertTokensProduceProgram(t, expectedAst, ([]Token)(tokens)...)
}

func parseSourceDefault(t *testing.T, source string) (Stmt, error) {
	tokens, err := tokenize(source)
	assert.NoError(t, err)
	stream := makeTokenizerStream(tokens...)
	return parseProgram(stream)
}

func assertSourceProducesExpression(t *testing.T, expectedAst Exp, source string) {
	tokens, err := tokenize(source)
	assert.NoError(t, err)
	assertTokensProduceExpression(t, expectedAst, ([]Token)(tokens)...)
}
func TestExpressionSimple(t *testing.T) {
	assertSourceProducesExpression(t,
		lessThan(plus(number(1), number(1)), number(2)),
		"1 + 1 < 2")
}
