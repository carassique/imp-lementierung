package imp

func Execute(source string) {
	context := ExecutionContext{
		out:    make(PrintChannel, 1000),
		signal: make(SignalChannel, 0),
	}
	tokens, _ := tokenize(source)
	program, _ := parseFromTokens(tokens)
	closure := makeRootTypeClosure()
	program.check(closure)
	if program != nil {
		print("\n\n" + program.pretty() + "\n")
		execClosure := makeRootValueClosure(context)
		go func() {
			program.eval(execClosure)
			close(context.out)
			if len(execClosure.getErrorStack()) == 0 {
				context.signal <- true
			} else {
				print(execClosure.errorStackToString())
				context.signal <- false
			}
		}()

		for {
			line, more := <-context.out
			if more == false {
				break
			} else {
				println(line)
			}
		}
		for {
			<-context.signal
			break
		}
	}
}
