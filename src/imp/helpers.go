package imp

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

/*
	Precedence rules:
	Standard precedence

	Normalize:
	plus ::= mult plusRhs
	plusRhs::= + plus |
	mult ::= det multRhs
	multRhs ::= * det mulRhs |

	Normalized expression grammar
	exp ::= dexp | dexp rhs
	dexp ::=
		| val
		| ! dexp
		| ( dexp )

	val ::= 0 | 1 | -1 | ...
		| "true" | "false"
		| vars

	rhs ::= + exp
		|	* exp
		|	|| exp
		| 	< exp


	Normalized statement grammar:
	statement ::= concreteStatement | concreteStatement; statement
	concreteStatement ::= vars ... |
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
	case v.flag == ValueBool:
		s = Bool(v.valB).pretty()
	case v.flag == Undefined:
		s = "Undefined"
	default:
		s = "UnknownError"
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

func lessThan(x, y Exp) Exp {
	return LessThan{x, y}
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

func not(x Exp) Exp {
	return Not{x}
}

func sequenceStatement(stmt1 Stmt, stmt2 Stmt) Stmt {
	return Seq{
		stmt1, stmt2,
	}
}

func declarationStatement(name string, exp Exp) Stmt {
	return Decl{
		lhs: name,
		rhs: exp,
	}
}

func assignmentStatement(name string, exp Exp) Stmt {
	return Assign{
		lhs: name,
		rhs: exp,
	}
}

func whileStatement(cond Exp, body Stmt) Stmt {
	return While{
		cond: cond,
		stmt: body,
	}
}

func printStatement(x Exp) Stmt {
	return Print{
		exp: x,
	}
}

func variableExpression(name string) Exp {
	return Var(name)
}
