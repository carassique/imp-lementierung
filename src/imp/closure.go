package imp

import (
	"reflect"
	"strings"
)

type ClosureStateMap[T any] map[string]T

type ErrorStack[T any] []ClosureError[T]

type ClosureError[T any] struct {
	reason              string
	closure             Closure[T]
	offenderType        OffenderType
	offendingStatement  *Stmt
	offendingExpression *Exp
}

type ClosureState[T any] struct {
	makeStateMap     func() ClosureStateMap[T]
	name             string
	stateMap         ClosureStateMap[T]
	parentClosure    *ClosureState[T]
	errorStack       *ErrorStack[T]
	executionContext *ExecutionContext
	interrupted      *bool
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
	error(offender interface{}, reason string)
	pushError(err ClosureError[T])
	getErrorStack() ErrorStack[T]
	errorStackToString() string
	getExecutionContext() ExecutionContext
	isInterrupted() bool
}

func makeHeader(t interface{}) string {
	return "[" + reflect.TypeOf(t).Name() + "] "
}

func closureErrorToString[T any](err ClosureError[T]) string {
	switch err.offenderType {
	case Statement:
		return makeHeader(*err.offendingStatement) + err.reason
	case Expression:
		return makeHeader(*err.offendingExpression) + err.reason
	default:
		return "[] " + err.reason
	}
}

func (closure *ClosureState[T]) errorStackToString() string {
	errorStack := *closure.errorStack
	var sb strings.Builder
	sb.WriteString("\n============== ERROR STACK " + closure.name + " ====================\n")
	for _, err := range errorStack {
		sb.WriteString(closureErrorToString(err) + "\n")
	}
	sb.WriteString("============== ERROR STACK END " + closure.name + " =================\n")
	return sb.String()
}

func (closure *ClosureState[T]) isInterrupted() bool {
	return *closure.interrupted
}

func (closure *ClosureState[T]) getExecutionContext() ExecutionContext {
	return *closure.executionContext
}

const DEFAULT_RUNTIME_ERROR_STACK_LENGTH = 5

func (closure *ClosureState[T]) pushError(err ClosureError[T]) {
	stack := *closure.errorStack
	*closure.errorStack = append(stack, err)
}

func (closure *ClosureState[T]) error(offender interface{}, reason string) {
	switch v := offender.(type) {
	case Stmt:
		closure.pushError(ClosureError[T]{
			reason:             reason,
			closure:            closure,
			offenderType:       Statement,
			offendingStatement: &v,
		})
	case Exp:
		closure.pushError(ClosureError[T]{
			reason:              reason,
			closure:             closure,
			offenderType:        Expression,
			offendingExpression: &v,
		})
	default:
		closure.pushError(ClosureError[T]{
			reason:       reason,
			closure:      closure,
			offenderType: Unsupported,
		})
	}
	if closure.executionContext != nil {
		if len(*closure.errorStack) >= DEFAULT_RUNTIME_ERROR_STACK_LENGTH && !closure.isInterrupted() {
			*closure.interrupted = true
		}
	}
}

func (closure *ClosureState[T]) makeChild() Closure[T] {
	return &ClosureState[T]{
		makeStateMap:     closure.makeStateMap,
		stateMap:         closure.makeStateMap(),
		errorStack:       closure.errorStack,
		parentClosure:    closure,
		executionContext: closure.executionContext,
		interrupted:      closure.interrupted,
	}
}

func (closure *ClosureState[T]) getErrorStack() ErrorStack[T] {
	return *closure.errorStack
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
		closure.parentClosure.assign(key, value)
	}
}

type OffenderType string

const (
	Expression  OffenderType = "expression"
	Statement   OffenderType = "statement"
	Unsupported OffenderType = "unsupported"
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

func makeRootTypeClosure() Closure[Type] {
	errorStack := make(ErrorStack[Type], 0)
	makeStateMap := func() ClosureStateMap[Type] {
		return make(map[string]Type, 0)
	}
	closureState := ClosureState[Type]{
		makeStateMap: makeStateMap,
		stateMap:     makeStateMap(),
		errorStack:   &errorStack,
		name:         "TYPE-CHECKER",
	}
	return &closureState
}

func makeRootValueClosure(context ExecutionContext) Closure[Val] {
	errorStack := make(ErrorStack[Val], 0)
	makeStateMap := func() ClosureStateMap[Val] {
		return make(map[string]Val, 0)
	}
	interrupted := false
	closureState := ClosureState[Val]{
		makeStateMap:     makeStateMap,
		stateMap:         makeStateMap(),
		errorStack:       &errorStack,
		executionContext: &context,
		interrupted:      &interrupted,
		name:             "EVALUATOR",
	}
	return &closureState
}
