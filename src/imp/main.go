package imp

import (
	"fmt"
)

// Simple imperative language

/*
vars       Variable names, start with lower-case letter

prog      ::= block
block     ::= "{" statement "}"
statement ::=  statement ";" statement           -- Command sequence
            |  vars ":=" exp                     -- Variable declaration
            |  vars "=" exp                      -- Variable assignment
            |  "while" exp block                 -- While
            |  "if" exp block "else" block       -- If-then-else
            |  "print" exp                       -- Print

exp ::= 0 | 1 | -1 | ...     -- Integers
     | "true" | "false"      -- Booleans
     | exp "+" exp           -- Addition
     | exp "*" exp           -- Multiplication
     | exp "||" exp          -- Disjunction
     | exp "&&" exp          -- Conjunction
     | "!" exp               -- Negation
     | exp "==" exp          -- Equality test
     | exp "<" exp           -- Lesser test
     | "(" exp ")"           -- Grouping of expressions
     | vars                  -- Variables
*/

func mkInt(x int) Val {
	return Val{flag: ValueInt, valI: x}
}
func mkBool(x bool) Val {
	return Val{flag: ValueBool, valB: x}
}
func mkUndefined() Val {
	return Val{flag: Undefined}
}

func showVal(v Val) string {
	var s string
	switch {
	case v.flag == ValueInt:
		s = Num(v.valI).pretty()
		if v.valI < 0 {
			s = "(" + s + ")"
		}
	case v.flag == ValueBool:
		s = Bool(v.valB).pretty()
	case v.flag == Undefined:
		s = "Undefined"
	}
	return s
}

func showType(t Type) string {
	var s string
	switch {
	case t == TyInt:
		s = "Int"
	case t == TyBool:
		s = "Bool"
	case t == TyIllTyped:
		s = "Illtyped"
	}
	return s
}

// Helper functions to build ASTs by hand

func number(x int) Exp {
	return Num(x)
}

func boolean(x bool) Exp {
	return Bool(x)
}

func plus(x, y Exp) Exp {
	return (Plus)([2]Exp{x, y})

	// The type Plus is defined as the two element array consisting of Exp elements.
	// Plus and [2]Exp are isomorphic but different types.
	// We first build the AST value [2]Exp{x,y}.
	// Then cast this value (of type [2]Exp) into a value of type Plus.

}

func mult(x, y Exp) Exp {
	return (Mult)([2]Exp{x, y})
}

func and(x, y Exp) Exp {
	return (And)([2]Exp{x, y})
}

func or(x, y Exp) Exp {
	return (Or)([2]Exp{x, y})
}

// Examples

func run(e Exp) {
	s := make(map[string]Val)
	t := make(map[string]Type)
	fmt.Printf("\n ******* ")
	fmt.Printf("\n %s", e.pretty())
	fmt.Printf("\n %s", showVal(e.eval(s)))
	fmt.Printf("\n %s", showType(e.infer(t)))
}

func runStatement(e Stmt) {
	s := make(map[string]Val)
	fmt.Printf("\n ******* ")
	fmt.Printf("\n %s \n", e.pretty())
	e.eval(s)
}

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
	program.eval(s)
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
	seq.eval(s)
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

	ex1()
	ex2()
	ex3()
	ex4()
	ex5()
	ex6()
}
