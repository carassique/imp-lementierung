package imp

/*
	Statements type checking
*/

func (stmt Print) check(t Closure[Type]) bool {
	parameterType := stmt.exp.infer(t)
	if parameterType == TyIllTyped {
		t.error(stmt, "Ill typed parameter for print")
	}
	return parameterType != TyIllTyped

}

func (while While) check(t Closure[Type]) bool {
	conditionType := while.cond.infer(t)
	statementTypeCheckResult := while.stmt.check(t.makeChild())
	if conditionType != TyBool {
		t.error(while, "Unsupported condition type: "+(string)(conditionType)+", expected TyBool")
	}
	if !statementTypeCheckResult {
		t.error(while, "Body of the while statement did not pass type checking")
	}
	return conditionType == TyBool && statementTypeCheckResult
}

func (ite IfThenElse) check(t Closure[Type]) bool {
	conditionType := ite.cond.infer(t)
	thenStatementTypeCheckResult := ite.thenStmt.check(t.makeChild())
	elseStatementTypeCheckResult := ite.elseStmt.check(t.makeChild())
	return conditionType == TyBool && thenStatementTypeCheckResult && elseStatementTypeCheckResult
}

func (stmt Seq) check(t Closure[Type]) bool {
	if !stmt[0].check(t) {
		return false
	}
	return stmt[1].check(t)
}

func (decl Decl) check(t Closure[Type]) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}
	x := (string)(decl.lhs)
	t.declare(x, ty)
	return true //TODO: check redeclaration
}

func (a Assign) check(t Closure[Type]) bool {
	x := (string)(a.lhs)
	if t.has(x) {
		t1 := t.get(x)
		t2 := a.rhs.infer(t)
		if t1 == t2 {
			return true
		}
		t.error(a, "Trying to assign value of type \""+string(t2)+"\" to variable \""+x+"\" of type \""+string(t1)+"\"")
	} else {
		t.error(a, "Variable \""+x+"\" does not exist in this scope")
	}
	return false
}

/*
	Expression type inference
*/

func (e LessThan) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt { //TODO: validate spec
		return TyBool
	}
	return TyIllTyped
}

func (e Equals) infer(t Closure[Type]) Type {
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

func (e Not) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	if t1 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (x Var) infer(t Closure[Type]) Type {
	y := (string)(x)
	if t.has(y) {
		return t.get(y)
	} else {
		return TyIllTyped // variable does not exist yields illtyped
	}

}

func (x Bool) infer(t Closure[Type]) Type {
	return TyBool
}

func (x Num) infer(t Closure[Type]) Type {
	return TyInt
}

func (e Mult) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e Plus) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e And) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Or) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}
