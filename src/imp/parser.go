package imp

import (
	"errors"
	"unicode"
)

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

func isValue(tokenCandidate string) bool {
	return false
}

func isAmbiguous(tokenCandidate string) bool {
	return tokenCandidate == ASSIGNMENT
}

func isTerminal(tokenCandidate string, terminalTokens StringSet) bool {
	if _, ok := terminalTokens[tokenCandidate]; ok {
		return ok
	}
	return false
}

func advanceToken() {

}

func tokenize(sourceCode string, terminalTokens StringSet) []string {
	// TODO: simplify code
	tokenList := make([]string, 0)

	currentToken := ""
	tokenCandidate := ""
	for _, character := range sourceCode {
		if unicode.IsSpace(character) {
			if len(currentToken) > 0 {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, currentToken)
				} else {
					//TODO: error - no token recognized!
				}
				tokenCandidate = ""
				currentToken = ""
			}
		} else {
			// Ignore spaces between tokens
			currentToken += (string)(character)
			if isTerminal(currentToken, terminalTokens) {
				tokenCandidate = currentToken
				if !isAmbiguous(currentToken) {
					tokenList = append(tokenList, currentToken)
					currentToken = ""
					tokenCandidate = ""
				}
			} else {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, tokenCandidate)
					currentToken = (string)(character)
				}
				if isTerminal(currentToken, terminalTokens) {
					tokenCandidate = currentToken
					if !isAmbiguous(currentToken) {
						tokenList = append(tokenList, currentToken)
						currentToken = ""
						tokenCandidate = ""
					}
				}
			}
		}

	}
	if len(tokenCandidate) > 0 {
		tokenList = append(tokenList, tokenCandidate)
	}

	return tokenList
}

type StringSet map[string]struct{}

func toSet(tokens []string) StringSet {
	tokenSet := make(map[string]struct{})
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	return tokenSet
}

var terminalTokens = toSet([]string{
	OPEN_BLOCK_GROUPING, CLOSE_BLOCK_GROUPING,
	PRINT, WHILE, IF, ELSE, STATEMENT_DELIMITER,
	DECLARATION, ASSIGNMENT, ADD, MULTIPLY,
	OR, AND, NOT, EQUALS, LESS_THAN,
	OPEN_EXPRESSION_GROUPING, CLOSE_EXPRESSION_GROUPING,
})

func parse(sourceCode string) (Stmt, error) {

	//tokensArray := tokenize(sourceCode, terminalTokens)

	return nil, errors.New("Not implemented")
}
