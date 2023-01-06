package imp

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.NotNil(t, ast)
	if ast != nil {
		t.Log("Interpreted AST: \n{\n" + indent(ast.pretty()) + "}")
		typeClosure := makeRootTypeClosure()
		checked := ast.check(typeClosure)
		assert.True(t, checked) //Typecheck successful
		errorStackForTypecheck := typeClosure.getErrorStack()
		assert.Equal(t, len(errorStackForTypecheck), 0) //No errors occured
		if len(errorStackForTypecheck) > 0 {
			print(typeClosure.errorStackToString())
		}
		assertOutputEquals(t, ast, vals...)
	}
}

func assertOutputEquals(t *testing.T, program Stmt, values ...Val) {
	stack := makeStack(values...)
	failed := false
	counter := 0
	consumer := func(value Val) {
		if stack.isEmpty() {
			failed = true
			t.Log("Received more output than expected: " + valToString(value))
			t.FailNow()
		}
		expectedValue := stack.pop()
		counter++
		t.Log("[" + strconv.Itoa(counter) + "] Output: " + valToString(value) + " Expected: " + valToString(expectedValue))
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
	assertOutputEqualsSource(t, "{ print 1234; print -1234; print 5566; print 000 }", "1234", "-1234", "5566", "0")
	assertOutputEqualsSource(t, "{ print true }", "true")
	assertOutputEqualsSource(t, "{ print false }", "false")
	assertOutputEqualsSource(t, "{ val:=123; print val }", "123")
}

func TestEvalClosure(t *testing.T) {
	// Inner variables do not leak from context, also: redeclaration is allowed
	assertOutputEqualsSource(t, "{ var := 1; if var == 1 { print var; var := true; print var } else { print false }; print var }",
		"1", "true", "1")
}

func TestEvalWhile(t *testing.T) {
	assertOutputEqualsSource(t, "{ counter := 0; while counter < 10 { print counter; counter = counter + 1 }}",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9")

	assertOutputEqualsSource(t, `
	{ 
	  fibonacciPrevPrev := 0;
	  fibonacciPrev := 1;
	  continue := true;
	  count := 0; 
	  print 0;
	  print 1;
	  while continue == true { 
		count = count + 1;
		fibonacci := fibonacciPrev + fibonacciPrevPrev;
		fibonacciPrevPrev = fibonacciPrev;
		fibonacciPrev = fibonacci;
		if count == 10 {
			continue = false
		} else {
			print fibonacci
		}
	  }
	}`, "0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55")
}

func TestEvalIfThenElse(t *testing.T) {
	assertOutputEqualsSource(t, `
	{
		if true {
			if true == false {
				print 1
			} else {
				print 2
			}
		} else {
			print 3
		}
	}
	`, "2")
}

func TestEvalSimpleExpressions(t *testing.T) {
	// Boolean
	assertOutputEqualsSource(t, `
	{ print true; print false; val := false; print val }
	`, "true", "false", "false")

	// Integer
	assertOutputEqualsSource(t, `
	{ print 125; print -130; print 0; print -0; num := 578; print num }
	`, "125", "-130", "0", "0", "578")

	// Not
	assertOutputEqualsSource(t, `
	{ 
		print !true;
	  	print !false;
	  	bool:=false;
	  	print !bool;
	  	print !!bool;
	  	print !!!bool
	}
	`, "false", "true", "true", "false", "true")
	assertOutputEqualsSource(t, `
	{ 
		bool := true;
		bool =!bool;
		print bool
	}
	`, "false")

	// And
	assertOutputEqualsSource(t, `
	{
		print true && true;
		print true && false;
		print false && true;
		print false && false
	}
	`, "true", "false", "false", "false")

	// Or
	assertOutputEqualsSource(t, `
	{
		print true || true;
		print true || false;
		print false || true;
		print false || false
	}
	`, "true", "true", "true", "false")

	// Plus
	assertOutputEqualsSource(t, `
	{ 
		print 0 + 0;
		print 0 + 1;
		print 1 + 0;
		print 1 + 1 + 1;
		print 1 + -2
	}
	`, "0", "1", "1", "3", "-1")

	// Mult
	assertOutputEqualsSource(t, `
		{
			print 0 * 0;
			print 0 * 1;
			print 1 * 0;
			print 1 * 2;
			print 2 * 5;
			print -2 * 5;
			print 2 * 2 * 2
		}
		`, "0", "0", "0", "2", "10", "-10", "8")

	// Equals
	assertOutputEqualsSource(t, `
	{
		print true == true;
		print true == false;
		print false == true;
		print false == false
	}
	`, "true", "false", "false", "true")
	assertOutputEqualsSource(t, `
	{
		print 0 == 0;
		print 1 == 1;
		print 1 == 0;
		print 0 == 1;
		print 234 == 432;
		print 1 == -1
	}
	`, "true", "true", "false", "false", "false", "false")

	// Less Than
	assertOutputEqualsSource(t, `
	{
		print 0 < 10;
		print 10 < 20;
		print 20 < 10;
		print -10 < 10;
		print 10 < -10
	}
	`, BOOLEAN_TRUE, BOOLEAN_TRUE, BOOLEAN_FALSE, BOOLEAN_TRUE, BOOLEAN_FALSE)

	// // Parenthesis
	// assertOutputEqualsSource(t, `
	// {
	// }
	// `)

	// // Variable
	// assertOutputEqualsSource(t, `
	// {
	// }
	// `)
}

func TestEvalComplexExpressions(t *testing.T) {
	assertOutputEqualsSource(t, `
	{
		var := 1 < 2 && false == true;
		print var
	}
	`, "false")

	assertOutputEqualsSource(t, `
	{
		var := 2 < 1 || false;
		print var;
		print var == false;
		if var == false {
			number := 1234;
			print 1234 == number
		} else {
			print false
		}
	}
	`, "false", "true", "true")
}
