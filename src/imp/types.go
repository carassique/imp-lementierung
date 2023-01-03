package imp

// Values

type Kind int

const (
	ValueInt     Kind = 0
	ValueBool    Kind = 1
	Undefined    Kind = 2
	RuntimeError Kind = 3
)

type Val struct {
	flag Kind
	valI int
	valB bool
	err  error
}

// Types

type Type int

const (
	TyIllTyped Type = 0
	TyInt      Type = 1
	TyBool     Type = 2
)

// Value State is a mapping from variable names to values
type ValState map[string]Val

type RuntimeVariableState interface {
	isVariableDeclared(variableName string) bool
	getVariableType(variableName string) Kind
	hasType(Kind) bool
	getVariableValue(variableName string) Val
}

// Value State is a mapping from variable names to types
type TyState map[string]Type

func makeRootTypeClosure() TypeClosure {
	return TypeClosure{
		typeMap: make(TyState, 0),
	}
}

type TypeClosure struct {
	typeMap       TyState
	parentClosure *TypeClosure
}

type TypeVariableState interface {
	isVariableDeclared(variableName string) bool
	getVariableType(variableName string) Type
	declareVariable(variableName string, t Type)
	makeNewClosure() TypeClosure
}

func (closure *TypeClosure) makeNewClosure() TypeClosure {
	return TypeClosure{
		typeMap:       make(TyState, 0),
		parentClosure: closure,
	}
}

func (closure *TypeClosure) isVariableDeclared(variableName string) bool {
	_, declaredLocally := (closure.typeMap)[variableName]
	if !declaredLocally {
		if closure.parentClosure != nil {
			return closure.parentClosure.isVariableDeclared(variableName)
		}
	}
	return declaredLocally
}

func (closure *TypeClosure) getVariableType(variableName string) Type {
	val, declaredLocally := closure.typeMap[variableName]
	if declaredLocally {
		return val
	} else {
		if closure.parentClosure != nil {
			return closure.parentClosure.getVariableType(variableName)
		}
	}
	return TyIllTyped
}

func (closure *TypeClosure) declareVariable(variableName string, t Type) {
	closure.typeMap[variableName] = t
}

// Interface

type Exp interface {
	pretty() string
	eval(s ValState) Val
	infer(t TypeClosure) Type
}

type Stmt interface {
	pretty() string
	eval(s ValState)
	check(t TypeClosure) bool
}

// Statement cases (incomplete)

type Seq [2]Stmt
type Decl struct {
	lhs string
	rhs Exp
}
type IfThenElse struct {
	cond     Exp
	thenStmt Stmt
	elseStmt Stmt
}

type While struct {
	cond Exp
	stmt Stmt
}

type Assign struct {
	lhs string
	rhs Exp
}

type PrintChannel chan string
type SignalChannel chan bool

type Print struct {
	exp Exp
	out PrintChannel
}

// Expression cases (incomplete)

type Bool bool
type Num int
type Not [1]Exp
type Mult [2]Exp
type Plus [2]Exp
type And [2]Exp
type Or [2]Exp
type LessThan [2]Exp
type Equals [2]Exp
type Var string
