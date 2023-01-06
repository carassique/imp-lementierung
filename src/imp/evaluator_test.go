package imp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEvaluator(t *testing.T) {
	condition := (LessThan)([2]Exp{number(0),
		(Var)("iterator")})

	wh := While{
		cond: condition,
		stmt: Seq{
			Assign{
				lhs: "iterator",
				rhs: plus((Var)("iterator"), number(-1)),
			},
			Print{
				exp: (Var)("iterator"),
			},
		},
	}

	seq := Seq{Assign{
		lhs: "iterator",
		rhs: number(10),
	}, wh}

	runStatement(seq)
	t.Error("Evaluator error")
}

func stringToVal(value string) (Val, bool) {
	ok, token := isBoolean(value)
	if ok {
		return mkBool(token.booleanValue), true
	}
	ok, token = isInteger(value)
	if ok {
		return mkInt(token.integerValue), true
	}
	return mkUndefined(), false
}

func stringsToVal(t *testing.T, values ...string) []Val {
	vals := []Val{}
	for _, val := range values {
		parsed, ok := stringToVal(val)
		assert.True(t, ok)
		vals = append(vals, parsed)
	}
	return vals
}

func assertOutputEqualsSource(t *testing.T, programSource string, values ...string) {
	t.Log("----------------------------------------------")
	t.Log("Program: " + programSource)
	vals := stringsToVal(t, values...)
	tokens, err := tokenize(programSource)
	assert.NoError(t, err)
	ast, err := parseFromTokens(tokens)
	assert.NoError(t, err)
	typeClosure := makeRootTypeClosure()
	checked := ast.check(typeClosure)
	assert.True(t, checked)                       //Typecheck successful
	assert.Len(t, typeClosure.getErrorStack(), 0) //No errors occured
	assertOutputEquals(t, ast, vals...)
}

func assertOutputEquals(t *testing.T, program Stmt, values ...Val) {
	stack := makeStack(values...)
	failed := false
	consumer := func(value Val) {
		if stack.isEmpty() {
			failed = true
			t.Log("Received more output than expected: " + valToString(value))
			t.FailNow()
		}
		expectedValue := stack.pop()
		t.Log("Output: " + valToString(value) + " Expected: " + valToString(expectedValue))
		assert.Equal(t, expectedValue, value)
		if value != expectedValue {
			t.FailNow()
			failed = true
		}
	}
	closure := executeAst(program, consumer)
	assert.False(t, failed)                   //No unexpected result appeared
	assert.True(t, stack.isEmpty())           //Every expected value matched output value
	assert.Len(t, closure.getErrorStack(), 0) //No errors occured
}

func TestEvalPrint(t *testing.T) {
	assertOutputEqualsSource(t, "{ print -1 }", "-1")
	assertOutputEqualsSource(t, "{ print 0 }", "0")
	assertOutputEqualsSource(t, "{ print 1 }", "1")
	assertOutputEqualsSource(t, "{ print true }", "true")
	assertOutputEqualsSource(t, "{ print false }", "false")
	assertOutputEqualsSource(t, "{ val:=123; print val }", "123")

}

func TestWhile(t *testing.T) {

	// Infinite loop
	// counter := 2
	// for counter < 10 {
	// 	counter := counter + 1 // Receives counter from the outer scope
	// 	t.Log(counter)
	// 	counter = 10
	// 	t.Log(counter)
	// }
}
