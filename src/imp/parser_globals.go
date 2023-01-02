package imp

var operators InfixOperators

func op() InfixOperators {
	if !operators.initialized {
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
			terminal:       ADD,
			higherPriority: &mult,
		}
		and := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return And{
					lhs, rhs,
				}
			},
			terminal:       AND,
			higherPriority: &plus,
		}
		or := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Or{
					lhs, rhs,
				}
			},
			terminal:       OR,
			higherPriority: &and,
		}
		lessThan := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return LessThan{
					lhs, rhs,
				}
			},
			terminal:       LESS_THAN,
			higherPriority: &or,
		}
		equals := InfixOperator{
			make: func(lhs, rhs Exp) Exp {
				return Equals{
					lhs, rhs,
				}
			},
			terminal:       EQUALS,
			higherPriority: &lessThan,
		}
		initialized := true
		operators = InfixOperators{
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
