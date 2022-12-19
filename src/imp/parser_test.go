package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	program, error := parse("{ print 1 }")
	typeMap := make(map[string]Type)
	assert.NoError(t, error)
	assert.True(t, program.check(typeMap))
}

func TestTokenizer(t *testing.T) {
	tokenList := tokenize("print { } = == =", terminalTokens)
	t.Logf("%v", tokenList)
}
