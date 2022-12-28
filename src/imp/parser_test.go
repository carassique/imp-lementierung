package imp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSourceCodeFilesDirectory = "test_source"

func readAvailableTestSourceFiles() []string {
	entries, _ := os.ReadDir("./" + testSourceCodeFilesDirectory)
	fileNames := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames
}

func readSourceCodeFile(filename string) string {
	// TODO: handle error?
	data, _ := os.ReadFile("./" + testSourceCodeFilesDirectory + "/" + filename)
	return string(data)
}

func TestAllSourceFiles(t *testing.T) {
	filenames := readAvailableTestSourceFiles()
	for _, filename := range filenames {
		testSourceFile(t, filename)
	}
}

func testSourceFile(t *testing.T, filename string) {
	t.Log("Test started for " + filename)
	testSource(t, readSourceCodeFile(filename))
}

func testSource(t *testing.T, source string) {
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel, 100),
	}
	tokens := tokenize(source)
	t.Log("Tokens: [", tokens, "]")
	program, error := parseFromTokens(tokens, context)
	assert.NoError(t, error)
	typeMap := make(map[string]Type)
	assert.NoError(t, error)
	assert.True(t, program.check(typeMap))
	stateMap := make(map[string]Val)
	program.eval(stateMap)
	close(context.out)
	context.signal <- true
	//hasFinishedExecuting := false

	for {
		line, more := <-context.out
		if more == false {
			break
		} else {
			t.Log(line) // TODO: check no-output-programs
		}
	}
	t.Log("Test finished")
}

func TestTokenizer(t *testing.T) {
	t.Log("Tokenizer test")
	tokenList := tokenize("print 123 -11 ham Jam true { } = == =")
	t.Logf("%v", tokenList)
}
