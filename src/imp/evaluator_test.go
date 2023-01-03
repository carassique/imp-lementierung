package imp

import "testing"

func testEvaluator(t *testing.T) {
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
	t.Error("Evaluator error")
}

func TestWhile(t *testing.T) {

	// Infinite loop
	// counter := 2
	// for counter < 10 {
	// 	counter := counter + 1 // Receives counter from the outer scope
	// 	t.Log(counter)
	// 	counter = 10
	// 	t.Log(counter)
	// }
}
