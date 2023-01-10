package imp

func makeExecutionContext() ExecutionContext {
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel),
	}
	return context
}

type ExecutionConsumer func(value Val)

func executeAst(program Stmt, consumer ExecutionConsumer) Closure[Val] {
	context := makeExecutionContext()
	closure := makeRootValueClosure(context)
	go func() {
		program.eval(closure)
		close(context.out)
		if len(closure.getErrorStack()) == 0 {
			context.signal <- true
		} else {
			print("==== Runtime error:")
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

func Execute(source string, ignoreTypecheck bool) {
	tokens, err := tokenize(source)
	if err != nil {
		println("==== Tokenization error:")
		println(err.Error())
		return
	}
	program, err := parseFromTokens(tokens)
	if err != nil {
		println("==== Parsing error:")
		println(err.Error())
		return
	}

	closure := makeRootTypeClosure()
	check := program.check(closure)
	if !check {
		println("==== Typecheck error:")
		println(closure.errorStackToString())
		if !ignoreTypecheck {
			return
		}
	}

	if program != nil {
		print("\n\nInterpeted AST: \n{\n" + indent(program.pretty()) + "}\n\nOutput:\n")
		executeAst(program, func(value Val) {
			println(showVal(value))
		})
	}
}
