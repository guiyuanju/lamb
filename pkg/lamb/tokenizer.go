package lamb

import (
	"fmt"
	"slices"
)

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

func NewTokenizer(s string) Tokenizer {
	return Tokenizer{0, s, 0, -1}
}

func (t *Tokenizer) Scan() ([]Token, bool) {
	s := t.str
	var res []Token
	for t.cur < len(s) {
		c := s[t.cur]

		switch c {
		case ' ':
			t.column++
			t.advance()

		case '\n':
			t.line++
			t.column = 0
			t.advance()

		case '\\':
			t.column++
			res = append(res, Token{TokenLambda, "\\", t.line, t.column})
			t.advance()

		case '=':
			t.column++
			res = append(res, Token{TokenEq, "=", t.line, t.column})
			t.advance()

		case '(':
			t.column++
			res = append(res, Token{TokenLeftParen, "(", t.line, t.column})
			t.advance()

		case ')':
			t.column++
			res = append(res, Token{TokenRightParen, ")", t.line, t.column})
			t.advance()

		case '.':
			t.column++
			res = append(res, Token{TokenDot, ".", t.line, t.column})
			t.advance()

		default:
			if t.consume("//") {
				t.consumeTill('\n')
				t.advance()
				continue
			}
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

func (t *Tokenizer) advance() {
	t.cur++
}

func (t *Tokenizer) consume(s string) bool {
	oldCur := t.cur
	for _, c := range []byte(s) {
		if t.isEnd() || t.str[t.cur] != c {
			t.cur = oldCur
			return false
		}
		t.cur++
	}
	return true
}

func (t *Tokenizer) consumeTill(c byte) {
	for !t.isEnd() && t.peek() != c {
		t.cur++
	}
}

func (t *Tokenizer) consumeKeyword(s string) bool {
	oldCur := t.cur
	ok := t.consume(s)
	if !ok {
		return false
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
	return slices.Contains(speicals, c)
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
