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

func parseExpressionVariable(tokens TokenizerStream) (Exp, error) {
	token := tokens.pop()
	if token.tokenType == VariableName {
		return Var(token.token), nil
	}
	return nil, errors.New("Could not parse variable name")
}

func parseExpressionBoolean(tokens TokenizerStream) (Exp, error) {
	return nil, nil
}

func parseExpressionInteger(tokens TokenizerStream) (Exp, error) {
	token := tokens.pop()
	if token.tokenType == IntegerValue {
		return number(token.integerValue), nil
	}
	return nil, errors.New("Could not parse integer")
}

// func parseExpressionGrouping(tokens TokenizerStream) (Exp, error) {
// 	tokens.expectTerminal(OPEN_BLOCK_GROUPING)

// }

// func parseExpressionNot(tokens TokenizerStream) (Exp, error) {

// }

// func parseExpressionMultRhs(tokens TokenizerStream) (Exp, error) {

//}

func parseExpressionMult(tokens TokenizerStream) (Exp, error) {
	return parseExpressionValue(tokens)
}

func parseExpressionPlusRhs(tokens TokenizerStream) (Exp, error) {
	token := tokens.peek()
	if token.tokenType == Terminal && token.token == ADD {
		tokens.pop()
		return parseExpressionPlus(tokens)
	}
	// else skip
	//TODO: remove error
	return nil, errors.New("Could not parse PlusRhs")
}

func parseExpressionPlus(tokens TokenizerStream) (Exp, error) {
	lhs, lerr := parseExpressionMult(tokens)
	rhs, rerr := parseExpressionPlusRhs(tokens)
	if rerr == nil {
		return Plus{
			lhs,
			rhs,
		}, nil
	} else {
		return lhs, lerr
	}
}

func parseExpression(tokens TokenizerStream) (Exp, error) {
	return parseExpressionPlus(tokens)
}

func parseExpressionValue(tokens TokenizerStream) (Exp, error) {
	firstToken := tokens.peek()
	switch firstToken.tokenType {
	case IntegerValue:
		return parseExpressionInteger(tokens)
	case BooleanValue:
		return parseExpressionBoolean(tokens)
	case VariableName:
		return parseExpressionVariable(tokens)
	}
	return nil, errors.New("Could not parse value")
}

// func parseLeafExpression(tokens TokenizerStream) (Exp, error) {

// 		switch firstToken.token {
// 		case OPEN_BLOCK_GROUPING:
// 			return parseExpressionGrouping(tokens)
// 		case NOT:
// 			return parseExpressionNot(tokens)
// 		}
// 	}
// }

func parseStatementPrint(tokens TokenizerStream) (Stmt, error) {
	tokens.expectTerminal(PRINT)
	exp, error := parseExpression(tokens)
	//TODO: handle error
	return (Stmt)(Print{
		exp: exp,
		out: tokens.context.out,
	}), error
	//expect("print")

}

func parseStatementIfThenElse(tokens TokenizerStream) (Stmt, error) {
	tokens.expectTerminal(IF)
	condition, error := parseExpression(tokens)
	thenBlock, error := parseBlock(tokens)
	tokens.expectTerminal(ELSE)
	elseBlock, error := parseBlock(tokens)
	return IfThenElse{
		cond:     condition,
		thenStmt: thenBlock,
		elseStmt: elseBlock,
	}, error
}

func parseStatementWhile() {

}

func parseStatementVariableAssignment() {

}

func parseStatementVariableDeclaration() {

}

func parseStatementSequence() {

}

func parseStatement(tokens TokenizerStream) (Stmt, error) {
	return parseStatementPrint(tokens)
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

type TokenizerStream struct {
	tokenList *TokenizerResultData
	context   ExecutionContext
}

type ExecutionContext struct {
	out    PrintChannel
	signal SignalChannel
}

type TokenizerResult interface {
	pop() Token
	peek() Token
	expectTerminal(token string) bool
}

func (tokenList *TokenizerStream) expectTerminal(token string) bool {
	tokenFromList := tokenList.pop()
	if tokenFromList.tokenType == Terminal && tokenFromList.token == token {
		return true
	}
	return false
}

func (tokenList *TokenizerStream) peek() Token {
	return (*tokenList.tokenList)[0]
}

func (tokenList *TokenizerStream) pop() Token {
	var value Token
	deref := *tokenList.tokenList
	value = deref[0]
	*tokenList.tokenList = deref[1:]
	return value
}

func parseBlock(tokens TokenizerStream) (Stmt, error) {
	//token := tokens.pop()
	// if token.tokenType == Terminal && token.token == OPEN_BLOCK_GROUPING {
	// 	//start creating struct for block??
	// }
	//expect(OPEN_BLOCK_GROUPING)
	isValid := tokens.expectTerminal(OPEN_BLOCK_GROUPING)
	stmt, error := parseStatement(tokens)
	isValid = tokens.expectTerminal(CLOSE_BLOCK_GROUPING)
	println(isValid)
	return stmt, error
	// types: ...
	//expect(CLOSE_BLOCK_GROUPING)
}

func parseProgram(tokens TokenizerStream) (Stmt, error) {
	return parseBlock(tokens)
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

func parse(sourceCode string, context ExecutionContext) (Stmt, error) {

	tokensArray := tokenize(sourceCode, terminalTokens)
	ast, error := parseProgram(TokenizerStream{
		tokenList: &tokensArray,
		context:   context,
	})
	return ast, error
}
