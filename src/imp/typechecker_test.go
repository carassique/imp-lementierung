package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncompatibleBinaryExpressions(t *testing.T) {

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
