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

func makeDefaultContext() ExecutionContext {
	return ExecutionContext{}
}

func surroundWithBlock(token ...Token) []Token {
	wrappedTokenList := []Token{terminal(OPEN_BLOCK_GROUPING)}
	wrappedTokenList = append(wrappedTokenList, token...)
	wrappedTokenList = append(wrappedTokenList, terminal(CLOSE_BLOCK_GROUPING))
	return wrappedTokenList
}

func surroundWithParenthesis(token ...Token) []Token {
	wrappedTokenList := []Token{terminal(OPEN_EXPRESSION_GROUPING)}
	wrappedTokenList = append(wrappedTokenList, token...)
	wrappedTokenList = append(wrappedTokenList, terminal(CLOSE_EXPRESSION_GROUPING))
	return wrappedTokenList
}

func TestPrintStatement(t *testing.T) {
	// Notice: different token sequences can result in identical AST
	testTokens(t, printStatement(variableExpression("myvar")),
		terminal(PRINT), variable("myvar"))

	testTokens(t, printStatement(number(5)),
		terminal(PRINT), integer(5))
	testTokens(t, printStatement(number(5)),
		terminal(PRINT), popen(), integer(5), pclose())

	additionPrintStatement := printStatement(plus(
		number(10),
		number(-20),
	))
	testTokens(t,
		additionPrintStatement,
		terminal(PRINT), integer(10), terminal(ADD), integer(-20),
	)
	testTokens(t,
		additionPrintStatement,
		terminal(PRINT),
		terminal(OPEN_EXPRESSION_GROUPING),
		integer(10),
		terminal(CLOSE_EXPRESSION_GROUPING),
		terminal(ADD),
		integer(-20))
	testTokens(t,
		additionPrintStatement,
		terminal(PRINT),
		terminal(OPEN_EXPRESSION_GROUPING),
		integer(10),
		terminal(ADD),
		integer(-20),
		terminal(CLOSE_EXPRESSION_GROUPING),
	)

}

func testTokens(t *testing.T, expectedAst Stmt, tokenList ...Token) (Stmt, ExecutionContext, error) {
	context := makeDefaultContext()
	wrappedTokenList := surroundWithBlock(tokenList...)
	ast, error := parseFromTokens(wrappedTokenList, context)
	assert.NoError(t, error)
	assert.Equal(t, expectedAst, ast)
	return ast, context, error
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
