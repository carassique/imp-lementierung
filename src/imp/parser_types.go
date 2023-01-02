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

func (tokenList *TokenizerStream) expectTerminal(token string) bool {
	tokenFromList, err := tokenList.pop()
	if err == nil && tokenFromList.tokenType == Terminal && tokenFromList.token == token {
		return true
	}
	return false
}

func (tokenList *TokenizerStream) isEmpty() bool {
	return len(*tokenList.tokenList) == 0
}

func (tokenList *TokenizerStream) peek() (Token, error) {
	if tokenList.isEmpty() {
		return Token{}, errors.New("No more tokens left")
	}
	return (*tokenList.tokenList)[0], nil
}

func (tokenList *TokenizerStream) pop() (Token, error) {
	if tokenList.isEmpty() {
		return Token{}, errors.New("No more tokens left")
	}
	var value Token
	deref := *tokenList.tokenList
	value = deref[0]
	*tokenList.tokenList = deref[1:]
	return value, nil
}
