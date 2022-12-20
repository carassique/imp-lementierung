package imp

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"
)

func expect(token string, tokens TokenizerResult) {

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
	//expect("print")

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

type TokenType int

const (
	Terminal     TokenType = 0
	IntegerValue TokenType = 1
	BooleanValue TokenType = 2
	VariableName TokenType = 3
)

type Token struct {
	tokenType    TokenType
	token        string
	integerValue int
	booleanValue bool
}

type TokenizerResult interface {
	pop() Token
	expectTerminal(token string) bool
}

func (tokenList TokenizerResultData) expectTerminal(token string) bool {
	//tokenList.pop()
	return false
}

func (tokenList TokenizerResultData) pop() Token {
	var value Token
	value, tokenList = tokenList[0], tokenList[1:]
	return value
}

func parseBlock(tokens TokenizerResult) {
	//token := tokens.pop()
	// if token.tokenType == Terminal && token.token == OPEN_BLOCK_GROUPING {
	// 	//start creating struct for block??
	// }
	//expect(OPEN_BLOCK_GROUPING)
	parseStatement()
	// types: ...
	//expect(CLOSE_BLOCK_GROUPING)
}

func parseProgram(tokens TokenizerResult) {
	parseBlock(tokens)
}

func isInteger(tokenCandidate string) (bool, Token) {
	// TODO: consider returning parsed value as a tuple or struct
	if value, err := strconv.Atoi(tokenCandidate); err == nil {
		return true, Token{token: tokenCandidate, tokenType: IntegerValue, integerValue: value}
	}
	return false, Token{}
}

func isBoolean(tokenCandidate string) (bool, Token) {
	if tokenCandidate == BOOLEAN_TRUE {
		return true, Token{token: tokenCandidate, tokenType: BooleanValue, booleanValue: true}
	}
	if tokenCandidate == BOOLEAN_FALSE {
		return true, Token{token: tokenCandidate, tokenType: BooleanValue, booleanValue: false}
	}
	return false, Token{}
}

func isVariableName(tokenCandidate string) (bool, Token) {
	// TODO: implement variable format
	match, _ := regexp.MatchString("^[a-z]([A-Za-z]|[0-9])*$", tokenCandidate)
	if match {
		return true, Token{token: tokenCandidate, tokenType: VariableName}
	}
	return false, Token{}
}

func isValue(tokenCandidate string) (bool, Token) {
	isInteger, integerToken := isInteger(tokenCandidate)
	if isInteger {
		return true, integerToken
	}
	isBoolean, booleanToken := isBoolean(tokenCandidate)
	if isBoolean {
		return true, booleanToken
	}
	// Order matters
	isVariableName, variableNameToken := isVariableName(tokenCandidate)
	if isVariableName {
		return true, variableNameToken
	}
	return false, Token{}
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

type TokenizerResultData []Token

func tokenize(sourceCode string, terminalTokens StringSet) TokenizerResultData {
	// TODO: simplify code, use scanner or generic lexer
	tokenList := make([]Token, 0)

	currentToken := ""
	tokenCandidate := ""
	for _, character := range sourceCode {
		if unicode.IsSpace(character) {
			if len(currentToken) > 0 {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     currentToken,
					})
				} else {
					//TODO: is non-terminal?
					if constitutesValue, value := isValue(currentToken); constitutesValue {
						tokenList = append(tokenList, value)
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
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     currentToken,
					})
					currentToken = ""
					tokenCandidate = ""
				}
			} else {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     tokenCandidate,
					})
					currentToken = (string)(character)
				}
				if isTerminal(currentToken, terminalTokens) {
					tokenCandidate = currentToken
					if !isAmbiguous(currentToken) {
						tokenList = append(tokenList, Token{
							tokenType: Terminal,
							token:     currentToken,
						})
						currentToken = ""
						tokenCandidate = ""
					}
				}
			}
		}

	}
	// Anything remains after the last character, it should be matched
	if len(tokenCandidate) > 0 {
		tokenList = append(tokenList, Token{
			tokenType: Terminal,
			token:     tokenCandidate,
		})
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

	tokensArray := tokenize(sourceCode, terminalTokens)
	parseProgram(tokensArray)
	return nil, errors.New("Not implemented")
}
