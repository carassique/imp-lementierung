package imp

import (
	"errors"
)

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
	booleanToken, err := tokens.expectTokenType(BooleanValue)
	if err == nil {
		return boolean(booleanToken.booleanValue), err
	}
	return nil, err
}

func parseExpressionInteger(tokens TokenizerStream) (Exp, error) {
	token, err := tokens.expectTokenType(IntegerValue)
	if err == nil {
		return number(token.integerValue), err
	}
	return nil, err
}

func parseExpressionGenericRhs(tokens TokenizerStream, operator InfixOperator) (Exp, error) {
	_, err := tokens.peekTerminal(operator.terminal)
	// If it does not match terminal, it could mean the following things:
	// (a) expect rhs for different operator * + etc
	// (b) there is no rhs, it's a single value
	if err == nil {
		tokens.pop()
		return parseExpressionGeneric(tokens, operator)
	}
	// else skip
	//TODO: remove error
	return nil, errors.New("Could not parse Rhs for " + operator.terminal)
}

func parseExpressionGeneric(tokens TokenizerStream, operator InfixOperator) (Exp, error) {
	var lhs, lerr = (Exp)(nil), (error)(nil)
	if operator.higherPriority != nil {
		// Try different lhs
		lhs, lerr = parseExpressionGeneric(tokens, *operator.higherPriority)
	} else {
		// Lhs can only be a value
		lhs, lerr = parseExpressionValue(tokens)
	}
	rhs, rerr := parseExpressionGenericRhs(tokens, operator)
	if rerr == nil {
		return operator.make(lhs, rhs), nil
	} else {
		return lhs, lerr
	}
}

func parseExpressionBinaryOperator(tokens TokenizerStream) (Exp, error) {
	multiply := InfixOperator{
		make: func(lhs, rhs Exp) Exp {
			return Mult{
				lhs,
				rhs,
			}
		},
		terminal: MULTIPLY,
	}
	plus := InfixOperator{
		make: func(lhs, rhs Exp) Exp {
			return Plus{
				lhs, rhs,
			}
		},
		terminal:       ADD,
		higherPriority: &multiply,
	}
	and := InfixOperator{
		make: func(lhs, rhs Exp) Exp {
			return And{
				lhs, rhs,
			}
		},
		terminal:       AND,
		higherPriority: &plus,
	}
	or := InfixOperator{
		make: func(lhs, rhs Exp) Exp {
			return Or{
				lhs, rhs,
			}
		},
		terminal:       OR,
		higherPriority: &and,
	}
	lessThan := InfixOperator{
		make: func(lhs, rhs Exp) Exp {
			return LessThan{
				lhs, rhs,
			}
		},
		terminal:       LESS_THAN,
		higherPriority: &or,
	}

	return parseExpressionGeneric(tokens, lessThan)
}

func parseIsolatedExpression(tokens TokenizerStream) (Exp, error) {
	exp, err := parseExpression(tokens)
	if err == nil {
		if !tokens.isEmpty() {
			return nil, errors.New("some tokens were not consumed")
		}
	}
	return exp, err
}

func parseExpression(tokens TokenizerStream) (Exp, error) {
	return parseExpressionBinaryOperator(tokens)
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
		return nil, err
	}
	return nil, errors.New("Could not parse expression grouping")
}

func parseExpressionNegation(tokens TokenizerStream) (Exp, error) {
	_, err := tokens.expectTerminal(NOT)
	if err == nil {
		exp, err := parseExpression(tokens)
		if err == nil {
			return not(exp), err
		}
		return nil, err
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
	token, _ := tokens.pop()
	return nil, errors.New("Could not parse value, token: " + token.token + " of type " + string(token.tokenType))
}

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

func parseStatementWhile(tokens TokenizerStream) (Stmt, error) {
	_, err := tokens.expectTerminal(WHILE)
	if err != nil {
		return nil, err
	}
	exp, err := parseExpression(tokens)
	if err != nil {
		return nil, err
	}
	block, err := parseBlock(tokens)
	if err != nil {
		return nil, err
	}
	return (Stmt)(While{
		cond: exp,
		stmt: block,
	}), err
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
	case Terminal:
		switch token.token {
		case PRINT:
			return parseStatementPrint(tokens)
		case WHILE:
			return parseStatementWhile(tokens)
		case IF:
			return parseStatementIfThenElse(tokens)
		}
	}
	return nil, errors.New("Could not parse statement, received " + token.token + " of type " + string(token.tokenType))
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
	ast, err := parseBlock(tokens)
	if err == nil {
		if !tokens.isEmpty() {
			return nil, errors.New("some tokens were not consumed")
		}
	}
	return ast, err
}

func parseFromTokens(tokens TokenizerResultData, context ExecutionContext) (Stmt, error) {
	ast, error := parseProgram(TokenizerStream{
		tokenList: &tokens,
		context:   context,
	})
	return ast, error
}

func parse(sourceCode string, context ExecutionContext) (Stmt, error) {
	tokensArray, err := tokenize(sourceCode)
	if err != nil {
		return nil, err
	}
	return parseFromTokens(tokensArray, context)
}
