package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintStatement_NonExistantVariable_TypeCheckerFalse(t *testing.T) {
	typeMap := make(map[string]Type)
	printStatement := Print{
		exp: (Var)("nonExistantVariable"),
	}
	typeCheckResult := printStatement.check(typeMap)
	assert.False(t, typeCheckResult)
}
