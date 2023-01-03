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
	get(key string) T
	set(key string, value T)
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
	_, declaredLocally := (closure.stateMap)[key]
	if !declaredLocally {
		if closure.parentClosure != nil {
			return closure.parentClosure.has(key)
		}
	}
	return declaredLocally
}

func (closure *ClosureState[T]) get(key string) T {
	value, declaredLocally := closure.stateMap[key]
	if declaredLocally {
		return value
	} else {
		if closure.parentClosure != nil {
			return closure.parentClosure.get(key)
		}
	}
	return value
}

func (closure *ClosureState[T]) set(key string, value T) {
	closure.stateMap[key] = value
}

// Interface

type Exp interface {
	pretty() string
	eval(s ValState) Val
	infer(t ClosureState[Type]) Type
}

type Stmt interface {
	pretty() string
	eval(s ValState)
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
