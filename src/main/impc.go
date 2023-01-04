package main

import (
	"imp"
	"os"
)

func main() {
	args := os.Args[1:]
	acc := ""
	for _, arg := range args {
		acc += arg + " "
	}
	print("Input: " + acc)
	imp.Execute(acc)
}
