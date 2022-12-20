package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	// parse("{ print 1337 }")
	program, error := parse("{ print 1337 }")

	program, error = parse("{ print 7 + 2 }")
	assert.NoError(t, error)
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
