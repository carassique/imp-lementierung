package imp

var operators InfixOperators

func op() InfixOperators {
	if !operators.initialized {
		lessThan := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return LessThan{
					lhs, rhs,
				}
			},
			terminal: LESS_THAN,
		}
		equals := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Equals{
					lhs, rhs,
				}
			},
			terminal: EQUALS,
		}
		mult := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Mult{
					lhs,
					rhs,
				}
			},
			terminal: MULTIPLY,
		}
		plus := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Plus{
					lhs, rhs,
				}
			},
			terminal: ADD,
		}

		and := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return And{
					lhs, rhs,
				}
			},
			terminal: AND,
		}
		or := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Or{
					lhs, rhs,
				}
			},
			terminal: OR,
		}

		precedence := [...]*InfixOperator{
			&mult,
			&plus,
			&lessThan,
			&equals,
			&and,
			&or,
		}

		var prevPtr *InfixOperator
		for _, entry := range precedence {
			entry.higherPriority = prevPtr
			prevPtr = entry
		}

		initialized := true
		operators = InfixOperators{
			lowest:      *prevPtr,
			initialized: initialized,
			plus:        plus,
			mult:        mult,
			or:          or,
			and:         and,
			equals:      equals,
			lessThan:    lessThan,
		}
	}
	return operators
}
