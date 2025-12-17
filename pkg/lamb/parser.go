package lamb

import "fmt"

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

func NewParser(tokens []Token) Parser {
	return Parser{0, tokens}
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

func (p *Parser) Parse() (Term, bool) {
	res, ok := p.parseExpr()
	if !p.isEnd() {
		p.printError(fmt.Sprintf("expect EOF, got \"%s\"", p.peek().value))
		return nil, false
	}
	return res, ok
}

// expr        ::= lambda | application | let
func (p *Parser) parseExpr() (Term, bool) {
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
func (p *Parser) parseLet() (Term, bool) {
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
	return Application{
		Left: Abstraction{
			Param: Variable(bind.value),
			Body:  body,
		},
		Right: value,
	}, true
}

// lambda      ::= "\" var "." expr
func (p *Parser) parseLambda() (Term, bool) {
	p.advance()
	id, ok := p.consume(TokenVar)
	if !ok {
		p.printError("expect a variable")
		return nil, false
	}
	param := Variable(id.value)
	_, ok = p.consume(TokenDot)
	if !ok {
		p.printError("expect dot")
		return nil, false
	}
	body, ok := p.parseExpr()
	if !ok {
		return nil, false
	}
	return Abstraction{param, body}, true
}

// application ::= atom { atom }
func (p *Parser) parseApplication() (Term, bool) {
	atom, ok := p.parseAtom()
	if !ok {
		return nil, false
	}
	var res Term = atom

	for !p.isEnd() &&
		(p.peek().ttype == TokenVar || p.peek().ttype == TokenLeftParen) {
		atom, ok = p.parseAtom()
		if !ok {
			return nil, false
		}
		res = Application{res, atom}
	}

	return res, true
}

// atom        ::= var | "(" expr ")"
func (p *Parser) parseAtom() (Term, bool) {
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
	return Variable(id.value), true
}
