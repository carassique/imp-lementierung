package imp

import (
	"errors"
)

func expect(token string, tokens TokenizerResult) {

}

func accept() {

}

func parseExpressionVariable(tokens TokenizerStream) (Exp, error) {
	token, err := tokens.pop()
	if err == nil && token.tokenType == VariableName {
		return Var(token.token), nil
	}
	return nil, errors.New("Could not parse variable name")
}

func parseExpressionBoolean(tokens TokenizerStream) (Exp, error) {
	return nil, nil
}

func parseExpressionInteger(tokens TokenizerStream) (Exp, error) {
	token, err := tokens.pop()
	if err == nil && token.tokenType == IntegerValue {
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
	token, err := tokens.peek()
	if err == nil && token.tokenType == Terminal && token.token == ADD {
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

func parseExpressionGrouping(tokens TokenizerStream) (Exp, error) {
	ok := tokens.expectTerminal(OPEN_EXPRESSION_GROUPING)
	if ok {
		exp, err := parseExpression(tokens)
		if err == nil {
			ok = tokens.expectTerminal(CLOSE_EXPRESSION_GROUPING)
			if ok {
				return exp, err
			}
		}
	}
	return nil, errors.New("Could not parse expression grouping")
}

func parseExpressionNegation(tokens TokenizerStream) (Exp, error) {
	token, err := tokens.pop()
	if err == nil && token.tokenType == Terminal && token.token == NOT {
		return parseExpression(tokens)
	}
	return nil, errors.New("Could not parse logic negation")
}

func parseExpressionValue(tokens TokenizerStream) (Exp, error) {
	firstToken, err := tokens.peek()
	if err != nil {
		return nil, errors.New("Token list empty while expecting a token")
	}
	switch firstToken.tokenType {
	case IntegerValue:
		return parseExpressionInteger(tokens)
	case BooleanValue:
		return parseExpressionBoolean(tokens)
	case VariableName:
		return parseExpressionVariable(tokens)
	case Terminal:
		switch firstToken.token {
		case NOT:
			return parseExpressionNegation(tokens)
		case OPEN_EXPRESSION_GROUPING:
			return parseExpressionGrouping(tokens)
		}
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

func advanceToken() {

}

func parseFromTokens(tokens TokenizerResultData, context ExecutionContext) (Stmt, error) {
	ast, error := parseProgram(TokenizerStream{
		tokenList: &tokens,
		context:   context,
	})
	return ast, error
}

func parse(sourceCode string, context ExecutionContext) (Stmt, error) {
	tokensArray := tokenize(sourceCode)
	return parseFromTokens(tokensArray, context)
}
