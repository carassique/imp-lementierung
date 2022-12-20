package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	// parse("{ print 1337 }")
	//program, error := parse("{ print 1337 }")
	t.Log("Test started")
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel, 100),
	}
	program, error := parse("{ print 7 + 2 + 1000 }", context)
	assert.NoError(t, error)
	typeMap := make(map[string]Type)
	assert.NoError(t, error)
	assert.True(t, program.check(typeMap))
	stateMap := make(map[string]Val)
	program.eval(stateMap)
	context.signal <- true
	hasFinishedExecuting := false
	for !hasFinishedExecuting {
		select {
		case line := <-context.out:
			t.Log(line)
		case <-context.signal:
			hasFinishedExecuting = true
		}
	}
	t.Log("Test finished")
}

func TestTokenizer(t *testing.T) {
	t.Log("Tokenizer test")
	tokenList := tokenize("print 123 -11 ham Jam true { } = == =", terminalTokens)
	t.Logf("%v", tokenList)
}
