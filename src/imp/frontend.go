package imp

func makeExecutionContext() ExecutionContext {
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel),
	}
	return context
}

type ExecutionConsumer func(value Val)

type ExecutionValidatorResult struct {
	failed *bool
}

type ValidatedChannel chan bool

func executeAst(program Stmt, consumer ExecutionConsumer) Closure[Val] {
	context := makeExecutionContext()
	closure := makeRootValueClosure(context)
	go func() {
		program.eval(closure)
		close(context.out)
		if len(closure.getErrorStack()) == 0 {
			context.signal <- true
		} else {
			print(closure.errorStackToString())
			context.signal <- false
		}
	}()
	for {
		line, more := <-context.out
		if !more {
			break
		} else {
			consumer(line)
		}
	}
	for {
		<-context.signal
		break
	}
	return closure
}

func Execute(source string) {

	tokens, _ := tokenize(source)
	program, _ := parseFromTokens(tokens)
	closure := makeRootTypeClosure()
	program.check(closure)

	if program != nil {
		print("\n\n" + program.pretty() + "\n")
		executeAst(program, func(value Val) {
			println(showVal(value))
		})
	}
}
