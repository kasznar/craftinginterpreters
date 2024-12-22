package main

import (
	"craftinginterpreters/src"
	"fmt"
	"os"
)

func main() {
	lox := src.Lox{}

	if len(os.Args) > 2 {
		fmt.Println("Usage: jlox [script]")
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
