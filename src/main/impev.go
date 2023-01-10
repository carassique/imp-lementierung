package main

import (
	"imp"
	"os"
)

func getHelpPrompt() {
	println("Try running \"impev help\" to get more information about available parameters.")
}

func doExecute(input string, ignoreTypecheck bool) {
	println("Input: " + input)
	imp.Execute(input, ignoreTypecheck)
}

func tryReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}

func expect(stack imp.Stack[string], terminal string) bool {
	return !stack.IsEmpty() && stack.Pop() == terminal
}

func expectPeek(stack imp.Stack[string], terminal string) bool {
	return !stack.IsEmpty() && stack.Peek() == terminal
}

func tryPop(stack imp.Stack[string]) (bool, string) {
	if stack.IsEmpty() {
		return false, ""
	}
	return true, stack.Pop()
}

func parseIgnoreTypecheck(stack imp.Stack[string]) bool {
	ok, flag := tryPop(stack)
	if ok {
		if flag == "-i" {
			return true
		}
		println("Expected flag -i, received: " + flag)
		getHelpPrompt()
		os.Exit(2)
	}
	return false
}

func tryWithFile(stack imp.Stack[string]) {
	if expect(stack, "-f") {
		// Assume next one is file
		ok, file := tryPop(stack)
		if ok {
			source, err := tryReadFile(file)
			if err != nil {
				println("Could not read file '" + file + "'")
			} else {
				doExecute(source, parseIgnoreTypecheck(stack))
			}
		} else {
			println("Please specify file name")
			getHelpPrompt()
		}
	} else {
		println("Incompatible arguments. Expected -f for 'file' specification.")
		getHelpPrompt()
	}
}

func main() {
	args := os.Args[1:]
	stack := imp.MakeStack(args...)

	if stack.IsEmpty() {
		println("Not enough arguments.")
		getHelpPrompt()
		return
	}
	if expectPeek(stack, "-f") {
		tryWithFile(stack)
	} else if expectPeek(stack, "help") {
		print(`
			Evaluator can be run in two ways:
			(1) Inline IMP code, e.g.:
			 impev "{ print 123 }"
			(2) File specification, e.g.:
			 impev -f "src/mysource.imp"

			Available additional flags:
			-i  Ignore typecheck
			will ignore typecheck results and run the evaluator either way
			has to be the last flag, e.g.:
			 impev -f "src/mysource.imp" -i
		`)
	} else {
		// Assume inline code
		doExecute(stack.Pop(), parseIgnoreTypecheck(stack))
	}
}
