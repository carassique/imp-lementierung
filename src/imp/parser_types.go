package imp

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
