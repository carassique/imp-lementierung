package imp

import (
	"strconv"
	"strings"
)

// pretty print

func (stmt Print) pretty() string {
	return "print " + stmt.exp.pretty()
}

func indent(multiline string) string {
	split := strings.Split(multiline, "\n")
	indented := ""
	for _, line := range split {
		indented += "    " + line + "\n"
	}
	return indented
}

func (ite IfThenElse) pretty() string {
	return "if " + ite.cond.pretty() + " {\n" + indent(ite.thenStmt.pretty()) + "} else {\n" + indent(ite.elseStmt.pretty()) + "}"
}

func (stmt While) pretty() string {
	return "while " + stmt.cond.pretty() + " {\n" + indent(stmt.stmt.pretty()) + "}"
}

func (stmt Assign) pretty() string {
	return stmt.lhs + " = " + stmt.rhs.pretty()
}

func (stmt Seq) pretty() string {
	return stmt[0].pretty() + ";\n" + stmt[1].pretty()
}

func (decl Decl) pretty() string {
	return decl.lhs + " := " + decl.rhs.pretty()
}

/////////////////////////
// Exp instances

// pretty print

func (x Equals) pretty() string {
	return "(" + x[0].pretty() + " == " + x[1].pretty() + ")"
}

func (x LessThan) pretty() string {
	return "(" + x[0].pretty() + " < " + x[1].pretty() + ")"
}

func (x Var) pretty() string {
	return (string)(x)
}

func (x Not) pretty() string {
	return "!" + x[0].pretty()
}

func (x Bool) pretty() string {
	if x {
		return "true"
	} else {
		return "false"
	}

}

func (x Num) pretty() string {
	value := int(x)
	strValue := strconv.Itoa(value)
	if value < 0 {
		return "(" + strValue + ")"
	}
	return strValue
}

func (e Mult) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += " * "
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Plus) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += " + "
	x += e[1].pretty()
	x += ")"

	return x
}

func (e And) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += " && "
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Or) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += " || "
	x += e[1].pretty()
	x += ")"

	return x
}
