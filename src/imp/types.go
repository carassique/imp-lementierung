package imp

// Values

type Kind string

const (
	ValueInt     Kind = "ValueInt"
	ValueBool    Kind = "ValueBool"
	Undefined    Kind = "Undefined"
	RuntimeError Kind = "RuntimeError"
)

type Val struct {
	flag Kind
	valI int
	valB bool
	err  error
}

// Types

type Type string

const (
	TyIllTyped Type = "TyIllTyped"
	TyInt      Type = "TyInt"
	TyBool     Type = "TyBool"
)

// Value State is a mapping from variable names to values
type ValState map[string]Val

// Value State is a mapping from variable names to types
type TyState map[string]Type

// Interface

type Exp interface {
	pretty() string
	eval(s Closure[Val]) Val
	infer(t Closure[Type]) Type
}

type Stmt interface {
	pretty() string
	eval(s Closure[Val])
	check(t Closure[Type]) bool
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

type PrintChannel chan Val
type SignalChannel chan bool

type Print struct {
	exp Exp
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
