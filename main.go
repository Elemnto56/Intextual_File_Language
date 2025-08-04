package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args

	// TUI
	if len(args) == 1 {
		//TUI()
		fmt.Println("You forgot the filename")
		os.Exit(1)
	}

	// Standard Path
	if len(args) > 2 {
		panic("Usage ./lexer <Intext File>")
	}

	filename := args[1]

	if !strings.HasSuffix(filename, ".itx") {
		panic("File is not an Intext File")
	}

	Lexer(filename)
	Parser()
	Validator()
	Interpreter()
}
