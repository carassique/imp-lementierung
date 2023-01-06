package main

import (
	"imp"
	"os"
)

func doExecute(input string) {
	println("Input: " + input)
	imp.Execute(input)
}

func tryReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		// Assume code
		doExecute(args[0])
	} else if len(args) == 2 {
		if args[0] == "-f" {
			// Assume next one is file
			file := args[1]
			source, err := tryReadFile(file)
			if err != nil {
				println("Could not read file '" + file + "'")
			} else {
				doExecute(source)
			}
		} else {
			println("Incompatible arguments. Expected -f for 'file' specification.")
		}
	} else {
		println("Not enough arguments.")
	}

}
