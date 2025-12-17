package lamb

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const EXT string = ".la"

func numberize(t Term) (int, bool) {
	a, ok := t.(Abstraction)
	if !ok {
		return 0, false
	}
	firstParam := a.Param

	a, ok = a.Body.(Abstraction)
	if !ok {
		return 0, false
	}
	secondParam := a.Param

	var count func(t Term) (int, bool)
	count = func(t Term) (int, bool) {
		switch t := t.(type) {
		case Abstraction:
			return 0, false
		case Variable:
			if t != secondParam {
				return 0, false
			}
			return 0, true
		case Application:
			left, ok := t.Left.(Variable)
			if !ok || left != firstParam {
				return 0, false
			}
			res, ok := count(t.Right)
			if !ok {
				return 0, false
			}
			return 1 + res, true
		default:
			panic("unrecognized term")
		}
	}

	return count(a.Body)
}

func preproc(s string) string {
	pat := `#use\s+(\w+)`
	re := regexp.MustCompile(pat)
	// text := "  #use std #use std "
	res := re.ReplaceAllStringFunc(s, func(match string) string {
		name := strings.Fields(match)[1] + EXT
		return loadFile(name)
	})
	return res
}

func Run(line string, w io.Writer) {
	line = preproc(line)
	tokenizer := NewTokenizer(line)
	tokens, ok := tokenizer.Scan()
	if !ok {
		fmt.Fprintln(w)
		return
	}
	if len(tokens) == 0 {
		return
	}
	parser := NewParser(tokens)
	term, ok := parser.Parse()
	if !ok {
		fmt.Fprintln(w)
		return
	}
	var rewrite int
	for {
		newTerm := Reduce(term)
		rewrite++

		irreducible := newTerm == term

		// Ensure printing for initially irreducible input
		if irreducible && rewrite == 1 || !irreducible {
			fmt.Fprintf(w, "R%d: %s", rewrite, newTerm)
			if n, ok := numberize(newTerm); ok {
				fmt.Fprintf(w, " -> %d", n)
			}
			fmt.Fprintf(w, "\n")
		}

		if irreducible {
			break
		}

		term = newTerm
	}
}

// func replWithReadline() {
// 	rl, err := readline.New("> ")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	for {
// 		line, err := rl.Readline()
// 		if err != nil {
// 			break
// 		}
// 		run(line)
// 	}
// }

func replBare() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !sc.Scan() {
			break
		}
		line := sc.Text()
		Run(line, os.Stdout)
	}
}

func Repl() {
	// replWithReadline()
	replBare()
}

func RunFile(name string, w io.Writer) {
	Run(loadFile(name), w)
}

func loadFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load file %s: %v", filename, err))
	}
	return string(content)
}
