package main

import (
	"strings"
	"syscall/js"

	. "la/pkg/lamb"
)

func run(this js.Value, args []js.Value) any {
	// Ensure js repl running, don't panic
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	if len(args) > 0 {
		code := args[0].String()
		var b strings.Builder
		Run(code, &b)
		return b.String()
	}

	return "no code provided"
}

func main() {
	c := make(chan struct{})
	js.Global().Set("run", js.FuncOf(run)) // Expose the Go function to the global JS scope
	<-c
}
