package imp

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_SOURCE_ROOT = "test_source"
const SHOULD_FAIL = "should_fail"
const SHOULD_PASS = "should_pass"

type File struct {
	filename  string
	directory string
}

func readAvailableTestSourceFiles(directory string) []File {
	relativeDir := "./" + directory
	entries, _ := os.ReadDir(relativeDir)
	files := make([]File, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, File{
				filename:  entry.Name(),
				directory: relativeDir,
			})
		}
	}
	return files
}

func readSourceCodeFile(file File) string {
	data, _ := os.ReadFile(file.directory + "/" + file.filename)
	return string(data)
}

func testAllSourceFiles(t *testing.T, relativeDir string, expectFail bool) {
	files := readAvailableTestSourceFiles(relativeDir)
	for _, file := range files {
		if expectFail {
			assertSourceFileFails(t, file)
		} else {
			assertSourceFilePasses(t, file)
		}
	}
}

func assertSourceFileFails(t *testing.T, file File) {
	t.Log("Test started for " + file.filename)
	//testSource(t, readSourceCodeFile(file))
	programSource := readSourceCodeFile(file)
	t.Log("----------------------------------------------")
	t.Log("Program: " + programSource)
	tokens, tokenizerError := tokenize(programSource)
	ast, parserError := parseFromTokens(tokens)
	checked := false
	typeErrorStack := false
	if ast != nil {
		t.Log("Interpreted AST: \n{\n" + indent(ast.pretty()) + "}")
		typeClosure := makeRootTypeClosure()
		checked = ast.check(typeClosure)

		errorStackForTypecheck := typeClosure.getErrorStack()
		if len(errorStackForTypecheck) > 0 {
			typeErrorStack = true
			print(typeClosure.errorStackToString())
		}
		assertExecutesWithError(t, ast)
	}
	if tokenizerError == nil && parserError == nil && checked && !typeErrorStack {
		t.Error("Program produced no errors for parsing and typechecking stage, but some were expected")
	}
}

func assertSourceFilePasses(t *testing.T, file File) {
	t.Log("Test started for " + file.filename)
	//testSource(t, readSourceCodeFile(file))
	programSource := readSourceCodeFile(file)
	t.Log("----------------------------------------------")
	t.Log("Program: " + programSource)
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
		assertExecutesWithoutError(t, ast)

	}
}

func assertExecutesWithError(t *testing.T, program Stmt) {
	counter := 0
	consumer := func(value Val) {
		counter++
		t.Log("[" + strconv.Itoa(counter) + "] Output: " + valToString(value))
	}
	closure := executeAst(program, consumer)
	errorStack := closure.getErrorStack()
	if len(errorStack) == 0 {
		t.Error("Executed without errors")
	} else {
		t.Log(closure.errorStackToString())
	}
}

func assertExecutesWithoutError(t *testing.T, program Stmt) {
	counter := 0
	consumer := func(value Val) {
		counter++
		t.Log("[" + strconv.Itoa(counter) + "] Output: " + valToString(value))
	}
	closure := executeAst(program, consumer)
	errorStack := closure.getErrorStack()
	if len(errorStack) > 0 {
		t.Error(closure.errorStackToString())
	}
}

func TestAllSources(t *testing.T) {
	testAllSourceFiles(t, TEST_SOURCE_ROOT+"/"+SHOULD_FAIL, true)
	testAllSourceFiles(t, TEST_SOURCE_ROOT+"/"+SHOULD_PASS, false)
}
