package main

import (
	"fmt"

	"github.com/chzyer/readline"
)

// Definition:
//
//	variable
//	abstraction
//	application
type (
	variable    string
	abstraction struct {
		param variable
		body  term
	}
)

type application struct {
	left  term
	right term
}
type term interface {
	isTerm()
}

func (v variable) isTerm()    {}
func (a abstraction) isTerm() {}
func (a application) isTerm() {}

func (a abstraction) String() string {
	return fmt.Sprintf("(λ%s.%s)", a.param, a.body)
}

func (a application) String() string {
	return fmt.Sprintf("(%s %s)", a.left, a.right)
}

// Capture-avoiding substitution M[x := N]
// - x[x := N] = N; y[x := N] = y if y != x
// - (M1 M2)[x := N] = (M1[x := N])(M2[x := N])
// - (λy.M)[x := N]
//   - If y == x, then λx.M (no change)
//   - If y ∉ FV(N), then λy.(M[x := N]) (FV(N): free variables of N)
//   - If y ∈ FV(N), then λy'.M[y := y'], where y' is fresh, then continue as above
func substitute(m term, x variable, n term) term {
	switch m := m.(type) {
	case variable:
		if m != x {
			return m
		}
		return n
	case application:
		return application{
			left:  substitute(m.left, x, n),
			right: substitute(m.right, x, n),
		}
	case abstraction:
		if m.param == x {
			return m
		}
		if !isFreeVariable(m.param, n) {
			return abstraction{
				param: m.param,
				body:  substitute(m.body, x, n),
			}
		}
		freshVar := getFreshVariable()
		newAbs := abstraction{
			param: freshVar,
			body:  substitute(m.body, m.param, freshVar),
		}
		return substitute(newAbs, x, n)
	default:
		panic("unrecognized term")
	}
}

var freshVarCounter int = 0

func getFreshVariable() variable {
	res := fmt.Sprintf("_%d", freshVarCounter)
	freshVarCounter++
	return variable(res)
}

func isFreeVariable(x variable, t term) bool {
	switch t := t.(type) {
	case variable:
		return x == t
	case application:
		return isFreeVariable(x, t.left) || isFreeVariable(x, t.right)
	case abstraction:
		if t.param == x {
			return false
		}
		return isFreeVariable(x, t.body)
	default:
		panic("unrecognized term")
	}
}

// Beta reduction
func betaReduce(m application) (term, bool) {
	if a, ok := reduce(m.left).(abstraction); ok {
		return substitute(a.body, a.param, m.right), true
	}
	return m, false
}

func reduce(t term) term {
	switch t := t.(type) {
	case variable:
		return t
	case application:
		// Substitute without reducing argument
		if t, ok := betaReduce(t); ok {
			return reduce(t)
		}
		// For non-betareducable expr, reduce argument for best effort
		return application{
			left:  t.left,
			right: reduce(t.right),
		}
	case abstraction:
		// Inspect the body of abstraction for possible reduction
		return abstraction{
			param: t.param,
			body:  reduce(t.body),
		}
	default:
		panic("unrecognized term")
	}
}

type TokenType int

const (
	TokenLambda = iota
	TokenLeftParen
	TokenRightParen
	TokenDot
	TokenVar
	TokenLet
	TokenIn
	TokenEq
)

type Token struct {
	ttype  TokenType
	value  string
	line   int
	column int
}

type Tokenizer struct {
	cur    int
	str    string
	line   int
	column int
}

func newTokenizer(s string) Tokenizer {
	return Tokenizer{0, s, 0, -1}
}

func (t *Tokenizer) scan() ([]Token, bool) {
	s := t.str
	var res []Token
	for t.cur < len(s) {
		c := s[t.cur]
		switch c {
		case ' ':
			t.column++
			t.cur++
		case '\n':
			t.line++
			t.column = 0
			t.cur++
		case '\\':
			t.column++
			res = append(res, Token{TokenLambda, "\\", t.line, t.column})
			t.cur++
		case '=':
			t.column++
			res = append(res, Token{TokenEq, "=", t.line, t.column})
			t.cur++
		case '(':
			t.column++
			res = append(res, Token{TokenLeftParen, "(", t.line, t.column})
			t.cur++
		case ')':
			t.column++
			res = append(res, Token{TokenRightParen, ")", t.line, t.column})
			t.cur++
		case '.':
			t.column++
			res = append(res, Token{TokenDot, ".", t.line, t.column})
			t.cur++
		default:
			if t.consumeKeyword("let") {
				res = append(res, Token{TokenLet, "let", t.line, t.column})
				continue
			}
			if t.consumeKeyword("in") {
				res = append(res, Token{TokenIn, "in", t.line, t.column})
				continue
			}
			name, ok := t.scanVar()
			if !ok {
				fmt.Printf("expected a variable, got \"%s\"\n", string(c))
				return res, false
			}
			res = append(res, name)
		}
	}
	return res, true
}

func (t *Tokenizer) isEnd() bool {
	return t.cur >= len(t.str)
}

func (t *Tokenizer) consumeKeyword(s string) bool {
	oldCur := t.cur
	for _, c := range []byte(s) {
		if t.isEnd() || t.str[t.cur] != c {
			t.cur = oldCur
			return false
		}
		t.cur++
	}
	if t.isEnd() || t.peek() == ' ' || t.peek() == '\n' {
		return true
	}
	t.cur = oldCur
	return false
}

func (t *Tokenizer) peek() byte {
	return t.str[t.cur]
}

func (t *Tokenizer) scanVar() (Token, bool) {
	if !isVar(t.peek()) {
		fmt.Printf(
			"unrecognized token \"%s\", expect a letter as the start of an variable",
			string(t.peek()))
		return Token{}, false
	}
	var name []byte
	for t.cur < len(t.str) && isVar(t.peek()) {
		t.column++
		name = append(name, t.str[t.cur])
		t.cur++
	}
	return Token{TokenVar, string(name), t.line, t.column}, true
}

func isVar(c byte) bool {
	return isSpecial(c) || isDigit(c) || isLetter(c)
}

func isSpecial(c byte) bool {
	speicals := []byte("!@#$%^&*_+{}[]:;\"'<>?,/|~`-=")
	for _, s := range speicals {
		if c == s {
			return true
		}
	}
	return false
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// expr        ::= lambda | application | let
// let         ::= "let" var "=" expr "in" expr
// lambda      ::= "\" var "." expr
// application ::= atom { atom }
// atom        ::= var | "(" expr ")"
// var         ::= identifier
type Parser struct {
	cur    int
	tokens []Token
}

func (p *Parser) peek() Token {
	return p.tokens[p.cur]
}

func (p *Parser) advance() {
	p.cur++
}

func (p *Parser) consume(token TokenType) (Token, bool) {
	if p.isEnd() {
		return Token{}, false
	}
	res := p.peek()
	if res.ttype != token {
		return Token{}, false
	}
	p.advance()
	return res, true
}

func (p *Parser) isEnd() bool {
	return p.cur >= len(p.tokens)
}

func (p *Parser) printError(info string) {
	var line, column int
	var token Token
	if p.isEnd() && len(p.tokens) > 0 {
		token = p.tokens[len(p.tokens)-1]
	} else {
		token = p.peek()
	}
	line = token.line
	column = token.column
	fmt.Printf("%d:%d %s\n", line, column, info)
}

func (p *Parser) parse() (term, bool) {
	res, ok := p.parseExpr()
	if !p.isEnd() {
		p.printError(fmt.Sprintf("expect EOF, got \"%s\"", p.peek().value))
		return nil, false
	}
	return res, ok
}

// expr        ::= lambda | application | let
func (p *Parser) parseExpr() (term, bool) {
	if p.peek().ttype == TokenLambda {
		return p.parseLambda()
	}
	if p.peek().ttype == TokenLet {
		return p.parseLet()
	}
	return p.parseApplication()
}

// let         ::= "let" var "=" expr "in" expr
// let f = N in M == (\f.M) N
func (p *Parser) parseLet() (term, bool) {
	p.advance()
	bind, ok := p.consume(TokenVar)
	if !ok {
		p.printError("expect a variable for a let binding name")
		return nil, false
	}
	_, ok = p.consume(TokenEq)
	if !ok {
		p.printError("expect \"=\"")
		return nil, false
	}
	value, ok := p.parseExpr()
	if !ok {
		return nil, false
	}
	_, ok = p.consume(TokenIn)
	if !ok {
		p.printError("expect \"in\"")
		return nil, false
	}
	body, ok := p.parseExpr()
	if !ok {
		return nil, false
	}
	return application{
		left: abstraction{
			param: variable(bind.value),
			body:  body,
		},
		right: value,
	}, true
}

// lambda      ::= "\" var "." expr
func (p *Parser) parseLambda() (term, bool) {
	p.advance()
	id, ok := p.consume(TokenVar)
	if !ok {
		p.printError("expect a variable")
		return nil, false
	}
	param := variable(id.value)
	_, ok = p.consume(TokenDot)
	if !ok {
		p.printError("expect dot")
		return nil, false
	}
	body, ok := p.parseExpr()
	if !ok {
		return nil, false
	}
	return abstraction{param, body}, true
}

// application ::= atom { atom }
func (p *Parser) parseApplication() (term, bool) {
	atom, ok := p.parseAtom()
	if !ok {
		return nil, false
	}
	var res term = atom

	for !p.isEnd() &&
		(p.peek().ttype == TokenVar || p.peek().ttype == TokenLeftParen) {
		atom, ok = p.parseAtom()
		if !ok {
			return nil, false
		}
		res = application{res, atom}
	}

	return res, true
}

// atom        ::= var | "(" expr ")"
func (p *Parser) parseAtom() (term, bool) {
	if p.peek().ttype == TokenLeftParen {
		p.consume(TokenLeftParen)
		expr, ok := p.parseExpr()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(TokenRightParen)
		if !ok {
			p.printError("expect right parenthesis")
			return nil, false
		}
		return expr, true
	}

	id, ok := p.consume(TokenVar)
	if !ok {
		p.printError("expect a variable")
		return nil, false
	}
	return variable(id.value), true
}

func repl() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		tokenizer := newTokenizer(line)
		tokens, ok := tokenizer.scan()
		if !ok {
			fmt.Println()
			continue
		}
		parser := Parser{0, tokens}
		term, ok := parser.parse()
		if !ok {
			fmt.Println()
			continue
		}
		res := reduce(term)
		fmt.Println(res)
	}
}

func main() {
	repl()
}

// (\x.M)N == let x = N in M

// let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ 0
// (λf.(λx.(f x)))

// let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ (succ 0)
// (λf.(λx.(f (f x))))

// let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in succ (succ (succ 0))
// (λf.(λx.(f (f (f x)))))

// let 0 = \f.\x.x in let succ = \n.\f.\x.f (n f x) in let 1 = succ 0 in let + = \m.\n.\f.\x .m f (n f x) in (+ 1 (+ 1 (+ 1 1)))
// (λf.(λx.(f (f (f (f x))))))
