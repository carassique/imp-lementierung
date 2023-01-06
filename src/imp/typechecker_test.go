package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertStatementPasses(t *testing.T, ast Stmt) {
	closure := makeRootTypeClosure()
	result := ast.check(closure)
	assert.True(t, result)
	if len(closure.getErrorStack()) > 0 {
		t.Log(closure.errorStackToString())
		t.Error("Statement produced errors on error stack")
	}
}

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
	t.Log(closure.errorStackToString())
	assert.True(t, foundMatchingExpression)
}

func assertFailsOnStatement(t *testing.T, ast Stmt, stmt Stmt) {
	closure := makeRootTypeClosure()
	result := ast.check(closure)
	assert.False(t, result)
	errorStack := closure.getErrorStack()
	var foundMatchingStatement = false
	for _, err := range errorStack {
		if err.offenderType == Statement {
			if *err.offendingStatement == stmt {
				foundMatchingStatement = true
				break
			}
		}
	}
	t.Log(closure.errorStackToString())
	assert.True(t, foundMatchingStatement)
}

func TestAssignmentTypecheck(t *testing.T) {
	// Assignment before declaration
	ast := assignmentStatement("something", number(5))
	assertFailsOnStatement(t, ast, ast)

	// Incompatible types
	declaration := declarationStatement("var", number(5))
	assignment := assignmentStatement("var", boolean(true))
	ast = sequenceStatement(declaration, assignment)
	assertFailsOnStatement(t, ast, assignment)

	// Reassignment
	assignment = assignmentStatement("var", number(10))
	assertStatementPasses(t, sequenceStatement(declaration, assignment))

	// Redeclaration in the same scope
	declaration2 := declarationStatement("var", number(123))
	assertStatementPasses(t, sequenceStatement(declaration, declaration2))

	// Redeclaration in the same scope with different type
	declaration2 = declarationStatement("var", boolean(false))
	assertStatementPasses(t, sequenceStatement(declaration, declaration2))
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
