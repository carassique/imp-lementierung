package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertFailsOnExpression(t *testing.T, ast Exp, exp Exp) {
	closure := makeRootTypeClosure()
	result := ast.infer(closure)
	assert.Equal(t, TyIllTyped, result)
	errorStack := closure.getErrorStack()
	var foundMatchingExpression = false
	for _, err := range errorStack {
		if err.offenderType == Expression {
			if *err.offendingExpression == exp {
				foundMatchingExpression = true
				break
			}
		}
	}
	assert.True(t, foundMatchingExpression)
}

func assertFailsOnStatement(t *testing.T, ast Stmt, stmt Stmt) {

}

func TestIncompatibleBinaryExpressions(t *testing.T) {
	num := number(1)
	orAst := or(num, boolean(true))
	ast := and(boolean(true), orAst)

	assertFailsOnExpression(t, ast, orAst)

}

func TestPrintStatement_NonExistantVariable_TypeCheckerFalse(t *testing.T) {
	closure := makeRootTypeClosure()
	printStatement := Print{
		exp: (Var)("nonExistantVariable"),
	}
	typeCheckResult := printStatement.check(closure)
	assert.False(t, typeCheckResult)
	t.Log(closure.errorStackToString())

}
