package main

import (
	"fmt"
	"os"

	. "la/pkg/lamb"

	"github.com/chzyer/readline"
)

func main() {
	// repl
	if len(os.Args) == 1 {
		Repl()
		return
	}

	// run from arg
	arg := os.Args[1]
	if arg == "-e" {
		if len(os.Args) < 3 {
			fmt.Println("expect code")
			return
		}
		code := os.Args[2]
		Run(code, os.Stdout)
		return
	}

	// run from file
	filename := arg
	RunFile(filename, os.Stdout)
}

func replWithReadline() {
	rl, err := readline.New("λ> ")
	if err != nil {
		panic(err)
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		Run(line, os.Stdout)
	}
}

func Repl() {
	replWithReadline()
	// replBare()
}
