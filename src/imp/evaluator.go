package imp

import (
	"errors"
	"fmt"
	"strconv"
)

// Evaluator

/////////////////////////
// Stmt instances

func (stmt Print) eval(s ValState) {
	stmt.out <- showVal(stmt.exp.eval(s))
}

func (stmt Assign) eval(s ValState) {
	s[stmt.lhs] = stmt.rhs.eval(s)
}

// eval

func (stmt Seq) eval(s ValState) {
	stmt[0].eval(s)
	stmt[1].eval(s)
}

func (ite IfThenElse) eval(s ValState) {
	v := ite.cond.eval(s)
	if v.flag == ValueBool {
		switch {
		case v.valB:
			ite.thenStmt.eval(s)
		case !v.valB:
			ite.elseStmt.eval(s)
		}

	} else {
		fmt.Printf("if-then-else eval fail")
	}
}

func (while While) eval(s ValState) {
	conditionHolds := true
	for conditionHolds {
		v := while.cond.eval(s)
		if v.flag == ValueBool {
			if v.valB {
				while.stmt.eval(s)
			} else {
				conditionHolds = false
			}
		}
	}
}

// Maps are represented via points.
// Hence, maps are passed by "reference" and the update is visible for the caller as well.
func (decl Decl) eval(s ValState) {
	v := decl.rhs.eval(s)
	x := (string)(decl.lhs)
	s[x] = v
}

func (exp Equals) eval(s ValState) Val {
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

func (exp LessThan) eval(s ValState) Val {
	return mkBool(exp[0].eval(s).valI < exp[1].eval(s).valI) // TOOD: implement checks
}

func (exp Var) eval(s ValState) Val {
	variableName := (string)(exp)
	return s[variableName] //TODO: implement eval time checks (variable exists?)
}

func (exp Not) eval(s ValState) Val {
	return mkBool(!exp[0].eval(s).valB) //TODO: implement eval time checks
}

func (x Bool) eval(s ValState) Val {
	return mkBool((bool)(x))
}

func (x Num) eval(s ValState) Val {
	return mkInt((int)(x))
}

func (e Mult) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI * n2.valI)
	}
	return mkUndefined()
}

func (e Plus) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
	return mkUndefined()
}

func (e And) eval(s ValState) Val {
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

func (e Or) eval(s ValState) Val {
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
