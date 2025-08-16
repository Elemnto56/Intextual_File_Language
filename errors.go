package main

import (
	"fmt"
	"os"
)

// Presets
var typemismatch string = "Check your code, and make sure to use the right conversions"

type ItxErr struct {
	Name        string
	Line        int
	Expression  string
	Description string
	AES         bool
	Hint        string
}

func NewError(name string, line int, expr string, desc string, approx bool, hint string) ItxErr {
	return ItxErr{
		Name:        name,
		Line:        line,
		Expression:  expr,
		Description: desc,
		AES:         approx,
		Hint:        hint,
	}
}

// Colors
var Turq, Green, Yellow, Red, Reset string = "\033[1;96m", "\033[1;92m", "\033[1;93m", "\033[31m", "\033[0m"

func (e ItxErr) Throw() {
	out := fmt.Sprintf("╭─[%sError: %s%s]────────────\n", Red, e.Name, Reset)
	if e.Line == 0 {
		out += fmt.Sprintf("│ %sLine:%s \033[3m%s\033[23m\n", Turq, Reset, "Undetermined")
	} else {
		out += fmt.Sprintf("│ %sLine:%s %d\n", Turq, Reset, e.Line)
	}
	out += "│\n"
	out += fmt.Sprintf("│ \033[0;35mIn expression:\033[0m %s\n", e.Expression)
	if e.AES {
		out += "│ \033[4;36m*The above expression is an approximation of what you typed\033[0m\n"
	}
	out += fmt.Sprintf("│\n│ \033[0;34mDesc:\033[0m %s\n", e.Description)
	if e.Hint != "" {
		out += fmt.Sprintf("│ %sHint:%s %s\n", Green, Reset, e.Hint)
	}
	out += "╰────────────────────────────\n"
	fmt.Println(out)
	os.Exit(1)
}
