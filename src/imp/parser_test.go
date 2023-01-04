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

func TestSequenceStatement(t *testing.T) {
	assertTokensProduceStatement(t,
		sequenceStatement(
			declarationStatement("myvar", number(5)),
			sequenceStatement(
				assignmentStatement("myvar", number(-125)),
				printStatement(variableExpression("myvar")),
			),
		),

		variable("myvar"), terminal(DECLARATION), integer(5),
		terminal(STATEMENT_DELIMITER),
		variable("myvar"), terminal(ASSIGNMENT), integer(-125),
		terminal(STATEMENT_DELIMITER),
		terminal(PRINT), variable("myvar"),
	)

	assertTokensProduceProgram(t,
		sequenceStatement(
			declarationStatement("myvar", number(5)),
			sequenceStatement(
				assignmentStatement("myvar", number(-125)),
				printStatement(variableExpression("myvar")),
			),
		),

		variable("myvar"), terminal(DECLARATION), integer(5),
		terminal(STATEMENT_DELIMITER),
		variable("myvar"), terminal(ASSIGNMENT), integer(-125),
		terminal(STATEMENT_DELIMITER),
		terminal(PRINT), variable("myvar"),
	)
}

func TestDeclarationStatement(t *testing.T) {
	assertTokensProduceStatement(t,
		declarationStatement("myvar", number(5)),
		variable("myvar"), terminal(DECLARATION), integer(5),
	)
}

func TestAssignmentStatement(t *testing.T) {
	assertTokensProduceStatement(t,
		assignmentStatement("myvar", number(5)),
		variable("myvar"), terminal(ASSIGNMENT), integer(5),
	)
}

func TestExpressionParser(t *testing.T) {
	assertTokensProduceExpression(t, plus(number(1), number(2)),
		integer(1), terminal(ADD), integer(2))

	// TODO: check if default parsed order of operations matches requirements
	assertTokensProduceExpression(t, plus(number(1), plus(number(2), number(3))),
		integer(1), terminal(ADD), integer(2), terminal(ADD), integer(3))
	assertTokensProduceExpression(t, plus(plus(number(1), number(2)), number(3)),
		popen(), integer(1), terminal(ADD), integer(2), pclose(), terminal(ADD), integer(3))
}

func TestEqualsExpression(t *testing.T) {
	testBinaryExpression(t, op().equals)
}

func TestLessThanExpression(t *testing.T) {
	testBinaryExpression(t, op().lessThan)
}

func TestOrExpression(t *testing.T) {
	testBinaryExpression(t, op().or)
}

func TestAndExpression(t *testing.T) {
	testBinaryExpression(t, op().and)
}

func TestMultExpression(t *testing.T) {
	testBinaryExpression(t, op().mult)
}

func TestPlusExpression(t *testing.T) {
	testBinaryExpression(t, op().and)
}

func testBinaryExpression(t *testing.T, operator InfixOperator) {
	assertTokensProduceExpression(t,
		operator.make(number(123), number(456)),
		integer(123), terminal(operator.terminal), integer(456),
	)

	assertTokensProduceExpression(t,
		operator.make(number(123), variableExpression("test")),
		integer(123), terminal(operator.terminal), variable("test"),
	)

	// TODO: add more tests
}

func TestVariableExpression(t *testing.T) {
	assertTokensProduceExpression(t, variableExpression("test"), variable("test"))
	// more tests implemented in tokenizer_test
}

func TestBooleanExpression(t *testing.T) {
	assertTokensProduceExpression(t, boolean(true), booleanToken(true))
	assertTokensProduceExpression(t, boolean(false), booleanToken(false))
}

func TestIntegerExpression(t *testing.T) {
	assertTokensProduceExpression(t, number(0), integer(0))
	assertTokensProduceExpression(t, number(1), integer(1))
	assertTokensProduceExpression(t, number(-1), integer(-1))
	assertTokensProduceExpression(t, number(0), Token{
		tokenType:    IntegerValue,
		token:        "-0",
		integerValue: 0,
	})
	assertTokensProduceExpression(t, number(1234), integer(1234))
	assertTokensProduceExpression(t, number(-1234), integer(-1234))
	// Invalid number input is checked by the tokenizer, see tokenizer_test
}

func TestNegationExpression(t *testing.T) {
	assertTokensProduceExpression(t,
		not(boolean(true)),
		terminal(NOT), booleanToken(true),
	)
	assertTokensProduceExpression(t,
		not(not(variableExpression("test"))),
		terminal(NOT), terminal(NOT), variable("test"),
	)
}

func TestExpressionGroupingHalfOpenParenthesis(t *testing.T) {
	// Half-closed parenthesis
	assertTokensProduceError(t, variable("test"), pclose())

}

func TestExpressionGroupingParser(t *testing.T) {
	assertTokensProduceExpression(t,
		variableExpression("test"),
		popen(), variable("test"), pclose(),
	)

	// Empty parenthesis are not part of language syntax -> error
	assertTokensProduceError(t, popen(), pclose())

	// Half-open parenthesis
	assertTokensProduceError(t, popen(), variable("test"))

	// Mismatched nested parenthesis
	assertTokensProduceError(t, popen(), variable("test"), popen(), pclose())
}

func TestBinaryOperatorExpressions(t *testing.T) {
	assertTokensProduceExpression(t,
		and(
			plus(
				variableExpression("a"), variableExpression("b"),
			),
			boolean(false),
		),

		variable("a"), terminal(ADD), variable("b"),
		terminal(AND), booleanToken(false),
	)
}

func TestIfThenElseStatement(t *testing.T) {
	ast := IfThenElse{
		cond:     boolean(true),
		thenStmt: printStatement(boolean(true)),
		elseStmt: printStatement(boolean(false)),
	}
	tokens := []Token{
		terminal(IF), booleanToken(true),
		terminal(OPEN_BLOCK_GROUPING),
		terminal(PRINT),
		booleanToken(true),
		terminal(CLOSE_BLOCK_GROUPING),
		terminal(ELSE),
		terminal(OPEN_BLOCK_GROUPING),
		terminal(PRINT),
		booleanToken(false),
		terminal(CLOSE_BLOCK_GROUPING),
	}
	assertTokensProduceStatement(t, ast, tokens...)
}

func TestWhileStatement(t *testing.T) {
	ast := whileStatement(
		lessThan(
			variableExpression("myvar"),
			number(100),
		),
		sequenceStatement(
			printStatement(variableExpression("myvar")),
			assignmentStatement("myvar",
				plus(variableExpression("myvar"), number(-1)))),
	)

	tokens := []Token{
		terminal(WHILE),
		variable("myvar"), terminal(LESS_THAN), integer(100),
		terminal(OPEN_BLOCK_GROUPING),
		terminal(PRINT), variable("myvar"), terminal(STATEMENT_DELIMITER),
		variable("myvar"), terminal(ASSIGNMENT),
		variable("myvar"), terminal(ADD), integer(-1),
		terminal(CLOSE_BLOCK_GROUPING)}

	assertTokensProduceStatement(t,
		ast,
		tokens...,
	)
	assertTokensProduceProgram(t,
		ast,
		tokens...)
}

func TestPrintStatement(t *testing.T) {
	// Notice: different token sequences can result in identical AST
	assertTokensProduceProgram(t, printStatement(variableExpression("myvar")),
		terminal(PRINT), variable("myvar"))

	assertTokensProduceProgram(t, printStatement(number(5)),
		terminal(PRINT), integer(5))
	assertTokensProduceProgram(t, printStatement(number(5)),
		terminal(PRINT), popen(), integer(5), pclose())

	additionPrintStatement := printStatement(plus(
		number(10),
		number(-20),
	))
	assertTokensProduceProgram(t,
		additionPrintStatement,
		terminal(PRINT), integer(10), terminal(ADD), integer(-20),
	)
	assertTokensProduceProgram(t,
		additionPrintStatement,
		terminal(PRINT),
		terminal(OPEN_EXPRESSION_GROUPING),
		integer(10),
		terminal(CLOSE_EXPRESSION_GROUPING),
		terminal(ADD),
		integer(-20))
	assertTokensProduceProgram(t,
		additionPrintStatement,
		terminal(PRINT),
		terminal(OPEN_EXPRESSION_GROUPING),
		integer(10),
		terminal(ADD),
		integer(-20),
		terminal(CLOSE_EXPRESSION_GROUPING),
	)
}

func TestBlock(t *testing.T) {
	// Empty block is not part of the language
	ast, err := parseBlock(makeTokenizerStream(
		terminal(OPEN_BLOCK_GROUPING),

		terminal(CLOSE_BLOCK_GROUPING),
	))
	assert.Error(t, err)
	assert.Nil(t, ast)
}

func assertTokensProduceError(t *testing.T, tokenList ...Token) {
	ast, err := parseExpressionFromTokensDefault(t, tokenList...)
	assert.Error(t, err)
	assert.Nil(t, ast)
}

func makeTokenizerStream(tokenList ...Token) TokenizerStream {
	tokenizerResult := (TokenizerResultData)(tokenList)
	tokenizerStream := TokenizerStream{
		tokenList: &tokenizerResult,
	}
	return tokenizerStream
}

func parseExpressionFromTokensDefault(t *testing.T, tokenList ...Token) (Exp, error) {
	return parseIsolatedExpression(makeTokenizerStream(tokenList...))
}

func assertTokensProduceExpression(t *testing.T, expectedAst Exp, tokenList ...Token) (Exp, error) {
	exp, err := parseExpressionFromTokensDefault(t, tokenList...)
	assert.NoError(t, err)
	assert.Equal(t, expectedAst, exp)
	return exp, err
}

func assertTokensProduceStatement(t *testing.T, expectedAst Stmt, tokenList ...Token) (Stmt, error) {
	tokenizerResult := (TokenizerResultData)(tokenList)
	tokenizerStream := TokenizerStream{
		tokenList: &tokenizerResult,
	}
	stmt, err := parseStatement(tokenizerStream)
	assert.NoError(t, err)
	assert.Equal(t, expectedAst, stmt)
	return stmt, err
}

func assertTokensProduceProgram(t *testing.T, expectedAst Stmt, tokenList ...Token) (Stmt, ExecutionContext, error) {
	context := makeDefaultContext()
	wrappedTokenList := surroundWithBlock(tokenList...)
	ast, error := parseFromTokens(wrappedTokenList)
	assert.NoError(t, error)
	assert.Equal(t, expectedAst, ast)
	return ast, context, error
}

func testSource(t *testing.T, source string) {
	// TODO: move to evaluator test
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel, 0),
	}
	tokens, err := tokenize(source)
	assert.NoError(t, err)
	t.Log("Tokens: [", tokens, "]")
	program, error := parseFromTokens(tokens)
	assert.NoError(t, error)
	//closure := makeRootTypeClosure()
	//assert.NoError(t, error)
	//assert.True(t, program.check(closure))
	//t.Log(closure.errorStackToString())
	t.Log("\n\n" + program.pretty())
	execClosure := makeRootValueClosure(context)
	go func() {
		program.eval(execClosure)
		close(context.out)
		if len(execClosure.getErrorStack()) == 0 {
			context.signal <- true
		} else {
			t.Log(execClosure.errorStackToString())
			context.signal <- false
		}
	}()

	for {
		line, more := <-context.out
		if more == false {
			break
		} else {
			t.Log(line)
		}
	}
	for {
		<-context.signal
		break
	}
	//close(context.out)
	//context.signal <- true
	//hasFinishedExecuting := false

	// for {
	// 	line, more := <-context.out
	// 	if more == false {
	// 		break
	// 	} else {
	// 		t.Log(line) // TODO: check no-output-programs
	// 	}
	// }

	t.Log("Test finished")
}

func TestTokenizer(t *testing.T) {
	t.Log("Tokenizer test")
	tokenList, err := tokenize("print 123 -11 ham jam true { } = == =")
	assert.NoError(t, err)
	t.Logf("%v", tokenList)
}
