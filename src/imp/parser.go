package imp

import (
	"errors"
)

func expect(token string, tokens TokenizerResult) {

}

func accept() {

}

func parseVariable(tokens TokenizerStream) (Var, error) {
	variable, err := tokens.expectTokenType(VariableName)
	if err == nil {
		// TODO: additional variable name format checks?
		return Var(variable.token), nil
	}
	return Var(""), err
}

func parseExpressionVariable(tokens TokenizerStream) (Exp, error) {
	return parseVariable(tokens)
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
	_, err := tokens.expectTerminal(OPEN_EXPRESSION_GROUPING)
	if err == nil {
		exp, err := parseExpression(tokens)
		if err == nil {
			_, err = tokens.expectTerminal(CLOSE_EXPRESSION_GROUPING)
			if err == nil {
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
	_, err := tokens.expectTerminal(PRINT)
	if err != nil {
		return nil, err
	}
	exp, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}
	return (Stmt)(Print{
		exp: exp,
		out: tokens.context.out,
	}), err
}

func parseStatementIfThenElse(tokens TokenizerStream) (Stmt, error) {
	_, err := tokens.expectTerminal(IF)
	if err != nil {
		return nil, err
	}
	condition, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}
	thenBlock, err := parseBlock(tokens)
	if err != nil {
		return nil, err
	}
	_, err = tokens.expectTerminal(ELSE)
	if err != nil {
		return nil, err
	}
	elseBlock, err := parseBlock(tokens)
	if err != nil {
		return nil, err
	}
	return IfThenElse{
		cond:     condition,
		thenStmt: thenBlock,
		elseStmt: elseBlock,
	}, nil
}

func parseStatementWhile() {

}

func parseStatementVariableDeclarationOrAssignment(tokens TokenizerStream) (Stmt, error) {
	// TODO: implement fork assignment/declaration
	variable, err := parseVariable(tokens)
	if err != nil {
		return nil, err
	}
	operand, err := tokens.expectTokenType(Terminal)
	if err != nil {
		return nil, err
	}
	if operand.token != ASSIGNMENT && operand.token != DECLARATION {
		return nil, errors.New("Expected variable declaration or assignment, received " + operand.token)
	}
	exp, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}
	switch operand.token {
	case DECLARATION:
		return declarationStatement(string(variable), exp), nil
	case ASSIGNMENT:
		return assignmentStatement(string(variable), exp), nil
	}
	return nil, errors.New("Parsing variable declaration or assignment failed")
}

func parseStatementConcrete(tokens TokenizerStream) (Stmt, error) {
	token, err := tokens.peek()
	if err != nil {
		return nil, errors.New("Expected statement, received nothing")
	}
	switch token.tokenType {
	case VariableName:
		return parseStatementVariableDeclarationOrAssignment(tokens)
	}
	return parseStatementPrint(tokens)
}

func parseStatement(tokens TokenizerStream) (Stmt, error) {
	stmt1, err1 := parseStatementConcrete(tokens)
	if err1 != nil {
		return nil, err1
	}
	token, err := tokens.peek()
	if err == nil {
		if token.tokenType == Terminal && token.token == STATEMENT_DELIMITER {
			tokens.pop()
			stmt2, err2 := parseStatement(tokens)
			if err2 != nil {
				return nil, err2
			}
			return sequenceStatement(stmt1, stmt2), err2
		}
		// else case can be a block, so leave it
	}
	return stmt1, err1
}

func parseBlock(tokens TokenizerStream) (Stmt, error) {
	_, err := tokens.expectTerminal(OPEN_BLOCK_GROUPING)
	if err == nil {
		// err is redeclared along with stmt... so multiple error returns are necessary
		stmt, err := parseStatement(tokens)
		if err == nil {
			_, err := tokens.expectTerminal(CLOSE_BLOCK_GROUPING)
			if err == nil {
				return stmt, nil
			}
			return nil, err
		}
		return nil, err
	}
	return nil, err
}

func parseProgram(tokens TokenizerStream) (Stmt, error) {
	return parseBlock(tokens)
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
