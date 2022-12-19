package imp

import (
	"errors"
	"regexp"
	"strconv"
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

func isInteger(tokenCandidate string) bool {
	// TODO: consider returning parsed value as a tuple or struct
	if _, err := strconv.Atoi(tokenCandidate); err == nil {
		return true
	}
	return false
}

func isBoolean(tokenCandidate string) bool {
	return tokenCandidate == BOOLEAN_TRUE || tokenCandidate == BOOLEAN_FALSE
}

func isVariableName(tokenCandidate string) bool {
	// TODO: implement variable format
	match, _ := regexp.MatchString("^[a-z]([A-Za-z]|[0-9])*$", tokenCandidate)
	return match
}

func isValue(tokenCandidate string) bool {
	return isInteger(tokenCandidate) || isVariableName(tokenCandidate) || isBoolean(tokenCandidate)
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
	// TODO: simplify code, use scanner or generic lexer
	tokenList := make([]string, 0)

	currentToken := ""
	tokenCandidate := ""
	for _, character := range sourceCode {
		if unicode.IsSpace(character) {
			if len(currentToken) > 0 {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, currentToken)
				} else {
					//TODO: is non-terminal?
					if isValue(currentToken) {
						tokenList = append(tokenList, currentToken)
					}

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
	// Anything remains after the last character, it should be matched
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
