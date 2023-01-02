package imp

import "errors"

type TokenType string

const (
	Terminal     TokenType = "Terminal"
	IntegerValue TokenType = "IntegerValue"
	BooleanValue TokenType = "BooleanValue"
	VariableName TokenType = "VariableName"
	Error        TokenType = "Error"
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

func (tokenList *TokenizerStream) expectTokenType(tokenType TokenType) (Token, error) {
	if tokenList.isEmpty() {
		return nothing(), errors.New("Expected " + string(tokenType) + ", received nothing")
	}
	tokenFromList, _ := tokenList.pop()
	if tokenFromList.tokenType == tokenType {
		return tokenFromList, nil
	}
	return tokenFromList, errors.New("Expected " + string(tokenType) + ", received " + string(tokenFromList.tokenType))
}

func (tokenList *TokenizerStream) expectTerminal(token string) (Token, error) {
	if tokenList.isEmpty() {
		return nothing(), errors.New("Expected terminal " + token + ", received nothing")
	}
	tokenFromList, err := tokenList.expectTokenType(Terminal)
	if err == nil && tokenFromList.token == token {
		return tokenFromList, nil
	}
	return tokenFromList, errors.New("Expected terminal " + token + ", received " + tokenFromList.token)
}

func (tokenList *TokenizerStream) isEmpty() bool {
	return len(*tokenList.tokenList) == 0
}

func (tokenList *TokenizerStream) peek() (Token, error) {
	if tokenList.isEmpty() {
		return nothing(), errors.New("No more tokens left")
	}
	return (*tokenList.tokenList)[0], nil
}

func (tokenList *TokenizerStream) pop() (Token, error) {
	if tokenList.isEmpty() {
		return nothing(), errors.New("No more tokens left")
	}
	var value Token
	deref := *tokenList.tokenList
	value = deref[0]
	*tokenList.tokenList = deref[1:]
	return value, nil
}
