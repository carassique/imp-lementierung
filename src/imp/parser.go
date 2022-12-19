package imp

import "errors"

func expect(token string) {

}

func accept() {

}

func parseExpressionVariable() {

}

func parseExpressionBoolean() {

}

func parseExpressionInteger() {

}

func parseExpression() {

}

func parseStatementPrint() {
	expect("print")

}

func parseStatementIfThenElse() {

}

func parseStatementWhile() {

}

func parseStatementVariableAssignment() {

}

func parseStatementVariableDeclaration() {

}

func parseStatementSequence() {

}

func parseStatement() {

}

func parseBlock() {
	expect("{")
	parseBlock()
	expect("}")
}

func parseProgram() {
	parseBlock()
}

func isValue() {

}

func isTerminal(){

}

func advanceToken(){

}

func tokenize(sourceCode string) string[] {

}


type StringSet map[string]struct{}

func toSet(tokens []string) StringSet {
	tokenSet := make(map[string]struct{})
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	return tokenSet
}

const OPEN_BLOCK_GROUPING = "{"
const CLOSE_BLOCK_GROUPING = "}"
const PRINT = "print"
const WHILE = "while"
const IF = "if"
const ELSE = "else"
const STATEMENT_DELIMITER = ";"
const DECLARATION = ":="
const ASSIGNMENT = "="
const ADD = "+"
const MULTIPLY = "*"
const OR = "||"
const AND = "&&"
const NOT = "!"
const EQUALS = "=="
const LESS_THAN = "<"
const OPEN_EXPRESSION_GROUPING = "("
const CLOSE_EXPRESSION_GROUPING = ")"


func parse(sourceCode string) (Stmt, error) {
	terminalTokens := toSet([...]string{
		OPEN_BLOCK_GROUPING, CLOSE_BLOCK_GROUPING,
		PRINT, WHILE, IF, ELSE, STATEMENT_DELIMITER, 
		DECLARATION, ASSIGNMENT, ADD, MULTIPLY,
		OR, AND, NOT, EQUALS, LESS_THAN,
		OPEN_EXPRESSION_GROUPING, CLOSE_EXPRESSION_GROUPING
	})
	
	return nil, errors.New("Not implemented")
}
