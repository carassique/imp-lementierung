package imp

/*
	Statements type checking
*/

func (stmt Print) check(t ClosureState[Type]) bool {
	parameterType := stmt.exp.infer(t)
	return parameterType != TyIllTyped

}

func (while While) check(t ClosureState[Type]) bool {
	conditionType := while.cond.infer(t)
	statementTypeCheckResult := while.stmt.check(t.makeChild())

	return conditionType == TyBool && statementTypeCheckResult
}

func (ite IfThenElse) check(t ClosureState[Type]) bool {
	conditionType := ite.cond.infer(t)
	thenStatementTypeCheckResult := ite.thenStmt.check(t.makeChild())
	elseStatementTypeCheckResult := ite.elseStmt.check(t.makeChild())
	return conditionType == TyBool && thenStatementTypeCheckResult && elseStatementTypeCheckResult
}

func (stmt Seq) check(t ClosureState[Type]) bool {
	if !stmt[0].check(t) {
		return false
	}
	return stmt[1].check(t)
}

func (decl Decl) check(t ClosureState[Type]) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}
	x := (string)(decl.lhs)
	t.declare(x, ty)
	return true //TODO: check redeclaration
}

func (a Assign) check(t ClosureState[Type]) bool {
	x := (string)(a.lhs)

	return t.has(x) && t.get(x) == a.rhs.infer(t)
}

/*
	Expression type inference
*/

func (e LessThan) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt { //TODO: validate spec
		return TyBool
	}
	return TyIllTyped
}

func (e Equals) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t) // TODO: check exists
	t2 := e[1].infer(t)

	if t1 == TyBool && t2 == TyBool { //TODO: validate spec
		return TyBool
	}
	if t1 == TyInt && t2 == TyInt {
		return TyBool
	}
	return TyIllTyped
}

func (e Not) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	if t1 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (x Var) infer(t ClosureState[Type]) Type {
	y := (string)(x)
	if t.has(y) {
		return t.get(y)
	} else {
		return TyIllTyped // variable does not exist yields illtyped
	}

}

func (x Bool) infer(t ClosureState[Type]) Type {
	return TyBool
}

func (x Num) infer(t ClosureState[Type]) Type {
	return TyInt
}

func (e Mult) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e Plus) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e And) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Or) infer(t ClosureState[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}
