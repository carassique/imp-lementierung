package imp

import (
	"fmt"
)

func ex1() {
	ast := plus(mult(number(1), number(2)), number(0))

	run(ast)
}

func ex2() {
	ast := and(boolean(false), number(0))
	run(ast)
}

func ex3() {
	ast := or(boolean(false), number(0))
	run(ast)
}

func ex4() {
	s := make(map[string]Val)
	program := IfThenElse{
		cond: boolean(true),
		thenStmt: Assign{
			lhs: "variableName1",
			rhs: number(1),
		},
		elseStmt: Assign{
			lhs: "variableName1",
			rhs: number(2),
		},
	}
	program.eval(makeRootValueClosure(ExecutionContext{}))
	println("\n*******")
	println(showVal(s["variableName1"]))
	// ast := plus(number(1), number(2))
	// run(ast)
	println("\n*******")
}

func ex5() {
	s := make(map[string]Val)

	condition := (LessThan)([2]Exp{number(0),
		(Var)("iterator")})

	wh := While{
		cond: condition,
		stmt: Seq{
			Assign{
				lhs: "iterator",
				rhs: plus((Var)("iterator"), number(-1)),
			},
			Print{
				exp: (Var)("iterator"),
			},
		},
	}

	seq := Seq{Assign{
		lhs: "iterator",
		rhs: number(10),
	}, wh}
	seq.eval(makeRootValueClosure(ExecutionContext{}))
	println("\n*******")
	println(showVal(s["iterator"]))
	println("\n*******")
}

func ex6() {
	condition := (LessThan)([2]Exp{number(0),
		(Var)("iterator")})

	wh := While{
		cond: condition,
		stmt: Seq{
			Assign{
				lhs: "iterator",
				rhs: plus((Var)("iterator"), number(-1)),
			},
			Print{
				exp: (Var)("iterator"),
			},
		},
	}

	seq := Seq{Assign{
		lhs: "iterator",
		rhs: number(10),
	}, wh}

	runStatement(seq)
}

func main() {

	fmt.Printf("\n")

	// ex1()
	// ex2()
	// ex3()
	// ex4()
	// ex5()
	// ex6()
}
