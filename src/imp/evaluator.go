package imp

import (
	"errors"
	"fmt"
	"strconv"
)

// Evaluator

/////////////////////////
// Stmt instances

func (stmt Print) eval(s Closure[Val]) {
	stmt.out <- showVal(stmt.exp.eval(s))
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
			// TODO: runtime error
			fmt.Print("Could not assign to: " + identifier)
		}
	}
}

// eval

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
		fmt.Printf("if-then-else eval fail")
	}
}

func (while While) eval(s Closure[Val]) {
	conditionHolds := true
	for conditionHolds {
		v := while.cond.eval(s)
		if v.flag == ValueBool {
			if v.valB {
				while.stmt.eval(s.makeChild())
			} else {
				conditionHolds = false
			}
		}
	}
}

// Maps are represented via points.
// Hence, maps are passed by "reference" and the update is visible for the caller as well.
func (decl Decl) eval(s Closure[Val]) {
	value := decl.rhs.eval(s)
	identifier := (string)(decl.lhs)
	s.declare(identifier, value) // TODO: runtime type validity check?
}

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
		return mkRuntimeError(errors.New("Unsupported type for equality check: " +
			strconv.FormatInt((int64)(e1.flag), 10))) //TODO: proper type to stirng conversion
	} else {
		// type mismatch! (should have been caught by the typechecker)
		// throw runtime error
		return mkRuntimeError(errors.New("Equals type mismatch"))
	}
}

func (exp LessThan) eval(s Closure[Val]) Val {
	lhs := exp[0].eval(s)
	rhs := exp[0].eval(s)
	if isRuntimeTypeCompatible(lhs, rhs) && lhs.flag == ValueInt {
		return mkBool(lhs.valI < rhs.valI)
	}
	return mkRuntimeError(errors.New("Incompatible variable types"))
}

func (exp Var) eval(s Closure[Val]) Val {
	identifier := (string)(exp)
	if s.has(identifier) {
		return s.get(identifier)
	} else {
		return mkRuntimeError(errors.New("Variable does not exist in this context: " + identifier))
	}
}

func (exp Not) eval(s Closure[Val]) Val {
	return mkBool(!exp[0].eval(s).valB) //TODO: implement eval time checks
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
	return mkUndefined()
}

func (e Plus) eval(s Closure[Val]) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
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
	return mkUndefined()
}
