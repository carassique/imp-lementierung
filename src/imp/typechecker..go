package imp

/*
	Statements type checking
*/

func (stmt Print) check(t TyState) bool {
	parameterType := stmt.exp.infer(t)
	return parameterType != TyIllTyped
}

func (while While) check(t TyState) bool {
	conditionType := while.cond.infer(t) //TODO: change as below
	statementTypeCheckResult := while.stmt.check(t)

	return conditionType == TyBool && statementTypeCheckResult
}

func (ite IfThenElse) check(t TyState) bool {
	conditionType := ite.cond.infer(t) // TODO: change to remove simplificating assumption
	thenStatementTypeCheckResult := ite.thenStmt.check(t)
	elseStatementTypeCheckResult := ite.elseStmt.check(t)
	return conditionType == TyBool && thenStatementTypeCheckResult && elseStatementTypeCheckResult
}

func (stmt Seq) check(t TyState) bool {
	if !stmt[0].check(t) {
		return false
	}
	return stmt[1].check(t)
}

func (decl Decl) check(t TyState) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}

	x := (string)(decl.lhs)
	t[x] = ty
	return true
}

func (a Assign) check(t TyState) bool {
	x := (string)(a.lhs)
	return t[x] == a.rhs.infer(t)
}

/*
	Expression type inference
*/

func (e LessThan) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt { //TODO: validate spec
		return TyBool
	}
	return TyIllTyped
}

func (e Equals) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)

	if t1 == TyBool && t2 == TyBool { //TODO: validate spec
		return TyBool
	}
	if t1 == TyInt && t2 == TyInt {
		return TyBool
	}
	return TyIllTyped
}

func (x Var) infer(t TyState) Type {
	y := (string)(x)
	ty, ok := t[y]
	if ok {
		return ty
	} else {
		return TyIllTyped // variable does not exist yields illtyped
	}

}

func (x Bool) infer(t TyState) Type {
	return TyBool
}

func (x Num) infer(t TyState) Type {
	return TyInt
}

func (e Mult) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e Plus) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e And) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Or) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}
