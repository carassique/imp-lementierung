package imp

// Evaluator

/////////////////////////
// Stmt instances

func (stmt Print) eval(s Closure[Val]) {
	s.getExecutionContext().out <- stmt.exp.eval(s)
}

func isValidValueType(val Val) bool {
	return val.flag == ValueBool || val.flag == ValueInt
}

func isRuntimeTypeCompatible(a Val, b Val) bool {
	// Assuming that allowing undefined and error type assignments is unhelpful
	return isValidValueType(a) && isValidValueType(b) && a.flag == b.flag
}

func (stmt Assign) eval(s Closure[Val]) {
	identifier := stmt.lhs
	value := stmt.rhs.eval(s)
	if s.has(identifier) {
		variable := s.get(identifier)
		if isRuntimeTypeCompatible(variable, value) {
			s.assign(identifier, value)
		} else {
			s.error(stmt, "Cannot assign "+expToString(stmt.rhs, value)+" to variable '"+identifier+"' of type "+valToString(variable))
		}
	} else {
		s.error(stmt, "Variable '"+identifier+"' has not been declared in this context")
	}
}

func valToString(val Val) string {
	return "[" + string(val.flag) + "]: '" + showVal(val) + "'"
}

// eval
func expToString(exp Exp, val Val) string {
	return valToString(val) + " in '" + exp.pretty() + "'"
}

func (stmt Seq) eval(s Closure[Val]) {
	stmt[0].eval(s)
	stmt[1].eval(s)
}

func (ite IfThenElse) eval(s Closure[Val]) {
	v := ite.cond.eval(s)
	if v.flag == ValueBool {
		switch {
		case v.valB:
			ite.thenStmt.eval(s.makeChild())
		case !v.valB:
			ite.elseStmt.eval(s.makeChild())
		}
	} else {
		s.error(ite, "Condition is not of boolean type: "+expToString(ite.cond, v))
	}
}

func (while While) eval(s Closure[Val]) {
	conditionHolds := true
	for conditionHolds && !s.isInterrupted() {
		v := while.cond.eval(s)
		if v.flag == ValueBool {
			if v.valB {
				while.stmt.eval(s.makeChild())
			} else {
				conditionHolds = false
			}
		} else {
			s.error(while, "Condition is not of boolean type: "+expToString(while.cond, v))
		}
	}
}

// Maps are represented via points.
// Hence, maps are passed by "reference" and the update is visible for the caller as well.
func (decl Decl) eval(s Closure[Val]) {
	value := decl.rhs.eval(s)
	identifier := (string)(decl.lhs)
	if value.flag == Undefined {
		s.error(decl, "Cannot declare variable '"+identifier+"' with value of undefined type: "+valToString(value))
	}
	s.declare(identifier, value)
}

func operatorToString(op Exp, lhs Val, rhs Val, text string) string {
	return text + valToString(lhs) + "; " + valToString(rhs) + " in " + op.pretty()
}

const (
	IncompatibleTypes = "Incompatible value types: "
)

func (exp Equals) eval(s Closure[Val]) Val {
	//TODO: verify spec for AND (short circuting allowed?)
	e1 := exp[0].eval(s)
	e2 := exp[1].eval(s)
	if e1.flag == e2.flag {
		switch e1.flag {
		case ValueBool:
			return mkBool(e1.valB == e2.valB)
		case ValueInt:
			return mkBool(e1.valI == e2.valI)
		}
		s.error(exp, operatorToString(exp, e1, e2, "Unsupported value types"))
		return mkUndefined()
	} else {
		// type mismatch! (should have been caught by the typechecker)
		// throw runtime error
		s.error(exp, operatorToString(exp, e1, e2, IncompatibleTypes))
		return mkUndefined()
	}
}

func (exp LessThan) eval(s Closure[Val]) Val {
	lhs := exp[0].eval(s)
	rhs := exp[1].eval(s)
	if isRuntimeTypeCompatible(lhs, rhs) && lhs.flag == ValueInt {
		return mkBool(lhs.valI < rhs.valI)
	}
	s.error(exp, operatorToString(exp, lhs, rhs, IncompatibleTypes))
	return mkUndefined()
}

func (exp Var) eval(s Closure[Val]) Val {
	identifier := (string)(exp)
	if s.has(identifier) {
		return s.get(identifier)
	} else {
		s.error(exp, "Variable '"+identifier+"' does not exist in this context.")
		return mkUndefined()
	}
}

func (exp Not) eval(s Closure[Val]) Val {
	val := exp[0].eval(s)
	if val.flag == ValueBool {
		return mkBool(!val.valB)
	}
	s.error(exp, "Not a boolean value: "+expToString(exp, val))
	return mkUndefined()
}

func (x Bool) eval(s Closure[Val]) Val {
	return mkBool((bool)(x))
}

func (x Num) eval(s Closure[Val]) Val {
	return mkInt((int)(x))
}

func (e Mult) eval(s Closure[Val]) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI * n2.valI)
	}
	s.error(e, operatorToString(e, n1, n2, IncompatibleTypes))
	return mkUndefined()
}

func (e Plus) eval(s Closure[Val]) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
	s.error(e, operatorToString(e, n1, n2, IncompatibleTypes))
	return mkUndefined()
}

func (e And) eval(s Closure[Val]) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == false:
		return mkBool(false)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB && b2.valB)
	}
	s.error(e, operatorToString(e, b1, b2, IncompatibleTypes))
	return mkUndefined()
}

func (e Or) eval(s Closure[Val]) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == true:
		return mkBool(true)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB || b2.valB)
	}
	s.error(e, operatorToString(e, b1, b2, IncompatibleTypes))
	return mkUndefined()
}
