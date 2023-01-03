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

func makeRootValueClosure() ClosureState[Val] {
	makeStateMap := func() ClosureStateMap[Val] {
		return make(map[string]Val, 0)
	}
	closureState := ClosureState[Val]{
		makeStateMap: makeStateMap,
		stateMap:     makeStateMap(),
	}
	return closureState
}

// Value State is a mapping from variable names to types
type TyState map[string]Type

func makeRootTypeClosure() ClosureState[Type] {
	makeStateMap := func() ClosureStateMap[Type] {
		return make(map[string]Type, 0)
	}
	closureState := ClosureState[Type]{
		makeStateMap: makeStateMap,
		stateMap:     makeStateMap(),
	}
	return closureState
}

type ClosureStateMap[T any] map[string]T

type ClosureState[T any] struct {
	makeStateMap  func() ClosureStateMap[T]
	stateMap      ClosureStateMap[T]
	parentClosure *ClosureState[T]
}

type Closure[T any] interface {
	has(key string) bool
	hasLocal(key string) bool
	get(key string) T
	getLocal(key string) (T, bool)
	setLocal(key string, value T)
	assign(key string, value T)
	declare(key string, value T)
	makeChild() Closure[T]
}

func (closure *ClosureState[T]) makeChild() ClosureState[T] {
	return ClosureState[T]{
		makeStateMap:  closure.makeStateMap,
		stateMap:      closure.makeStateMap(),
		parentClosure: closure,
	}
}

func (closure *ClosureState[T]) has(key string) bool {
	declaredLocally := closure.hasLocal(key)
	if !declaredLocally {
		if closure.parentClosure != nil {
			return closure.parentClosure.has(key)
		}
	}
	return declaredLocally
}

func (closure *ClosureState[T]) hasLocal(key string) bool {
	_, declaredLocally := closure.getLocal(key)
	return declaredLocally
}

func (closure *ClosureState[T]) getLocal(key string) (T, bool) {
	value, ok := closure.stateMap[key]
	return value, ok
}

func (closure *ClosureState[T]) get(key string) T {
	value, declaredLocally := closure.getLocal(key)
	if declaredLocally {
		return value
	} else {
		if closure.parentClosure != nil {
			return closure.parentClosure.get(key)
		}
	}
	return value
}

func (closure *ClosureState[T]) setLocal(key string, value T) {
	closure.stateMap[key] = value
}

func (closure *ClosureState[T]) declare(key string, value T) {
	// Notice: will overwrite previous declarations in the same closure scope
	closure.setLocal(key, value)
}

func (closure *ClosureState[T]) assign(key string, value T) {
	if closure.hasLocal(key) {
		closure.setLocal(key, value)
	} else {
		// Will try to assign
		// TODO: error reporting
		closure.parentClosure.assign(key, value)
	}
}

// Interface

type OffenderType string

const (
	Expression OffenderType = "expression"
	Statement  OffenderType = "statement"
)

type TypecheckDiagnostics struct {
	validated           bool
	offenderType        OffenderType
	offendingStatement  Stmt
	offendingExpression Exp
	description         string
}

func mkValid() TypecheckDiagnostics {
	return TypecheckDiagnostics{
		validated: true,
	}
}

func mkDiagStatement(stmt Stmt, desc string) TypecheckDiagnostics {
	return TypecheckDiagnostics{
		validated:          false,
		offenderType:       Statement,
		offendingStatement: stmt,
		description:        desc,
	}
}

func mkDiagExpression(exp Exp, desc string) TypecheckDiagnostics {
	return TypecheckDiagnostics{
		validated:           false,
		offenderType:        Expression,
		offendingExpression: exp,
		description:         desc,
	}
}

type Exp interface {
	pretty() string
	eval(s ClosureState[Val]) Val
	infer(t ClosureState[Type]) Type
}

type Stmt interface {
	pretty() string
	eval(s ClosureState[Val])
	check(t ClosureState[Type]) bool
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
