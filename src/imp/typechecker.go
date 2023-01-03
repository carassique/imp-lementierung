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
	if conditionType != TyBool {
		t.error(ite, "Condition type is not TyBool: "+ite.cond.pretty())
	}
	if !thenStatementTypeCheckResult {
		t.error(ite, "\"then\" branch did not pass type checking")
	}
	if !elseStatementTypeCheckResult {
		t.error(ite, "\"else\" branch did not pass type checking")
	}
	return conditionType == TyBool && thenStatementTypeCheckResult && elseStatementTypeCheckResult
}

func (stmt Seq) check(t Closure[Type]) bool {
	if !stmt[0].check(t) {
		t.error(stmt, "First statement of the sequence did not pass type checking")
		return false
	}
	if !stmt[1].check(t) {
		t.error(stmt, "Second statement of the sequence did not pass type checking")
		return false
	}
	return true
}

func (decl Decl) check(t Closure[Type]) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		t.error(decl, "Right hand side of the declaration statement is ill typed: '"+decl.rhs.pretty()+"'")
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

func tryAssertTypeOperator[O Exp, C Exp](t Closure[Type], root O, child C, actual Type, expected Type) {
	if actual != expected {
		t.error(root, "Expected "+mkType(expected)+" but received "+mkShard(actual, child)+" in '"+root.pretty()+"'")
	}
}

func mkType(t Type) string {
	return "[" + string(t) + "]"
}

func mkShard(t Type, e Exp) string {
	return mkType(t) + ": '" + e.pretty() + "'"
}

func (e LessThan) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt { //TODO: validate spec
		return TyBool
	}
	tryAssertTypeOperator(t, e, e[0], t1, TyInt)
	tryAssertTypeOperator(t, e, e[1], t2, TyInt)
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
	if t1 != t2 {
		t.error(e, "Operands have different types: ["+string(t1)+"] ["+string(t2)+"] in '"+e.pretty()+"'")
	}
	if t1 == TyIllTyped {
		t.error(e, "Left hand side is ill typed: "+mkShard(t1, e[0])+" in '"+e.pretty()+"'")
	}
	if t2 == TyIllTyped {
		t.error(e, "Right hand side is ill typed: "+mkShard(t2, e[1])+" in '"+e.pretty()+"'")
	}
	return TyIllTyped
}

func (e Not) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	if t1 == TyBool {
		return TyBool
	}
	t.error(e, "Is not a boolean type: "+mkShard(t1, e[0]))
	return TyIllTyped
}

func (x Var) infer(t Closure[Type]) Type {
	y := (string)(x)
	if t.has(y) {
		return t.get(y)
	} else {
		t.error(x, "Variable \""+string(x)+"\" is not declared for this context")
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
	tryAssertTypeOperator(t, e, e[0], t1, TyInt)
	tryAssertTypeOperator(t, e, e[1], t2, TyInt)
	return TyIllTyped
}

func (e Plus) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	tryAssertTypeOperator(t, e, e[0], t1, TyInt)
	tryAssertTypeOperator(t, e, e[1], t2, TyInt)
	return TyIllTyped
}

func (e And) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	tryAssertTypeOperator(t, e, e[0], t1, TyBool)
	tryAssertTypeOperator(t, e, e[1], t2, TyBool)
	return TyIllTyped
}

func (e Or) infer(t Closure[Type]) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	tryAssertTypeOperator(t, e, e[0], t1, TyBool)
	tryAssertTypeOperator(t, e, e[1], t2, TyBool)
	return TyIllTyped
}
