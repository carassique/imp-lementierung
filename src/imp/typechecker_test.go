package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintStatement_NonExistantVariable_TypeCheckerFalse(t *testing.T) {
	closure := makeRootTypeClosure()
	printStatement := Print{
		exp: (Var)("nonExistantVariable"),
	}
	typeCheckResult := printStatement.check(closure)
	assert.False(t, typeCheckResult)
}
