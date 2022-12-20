package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	program, error := parse("{ print 1337 }")
	typeMap := make(map[string]Type)
	assert.NoError(t, error)
	assert.True(t, program.check(typeMap))
	stateMap := make(map[string]Val)
	program.eval(stateMap)
}

func TestTokenizer(t *testing.T) {
	tokenList := tokenize("print 123 -11 ham Jam true { } = == =", terminalTokens)
	t.Logf("%v", tokenList)
}
