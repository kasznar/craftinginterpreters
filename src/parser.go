package src

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func makeParser(tokens []Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// todo: can modify receiver?
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().tokenType == tokenType
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for i := 0; i < len(tokenTypes); i++ {
		if p.check(tokenTypes[i]) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) parse() (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("syntax error", r)
		}
	}()

	return p.expression(), nil
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()

		return UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return LiteralExpr{Value: false}
	}
	if p.match(TRUE) {
		return LiteralExpr{Value: true}
	}
	if p.match(NIL) {
		return LiteralExpr{Value: nil}
	}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{Value: p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpr{Expression: expr}
	}

	// todo: better error object
	panic(fmt.Errorf("todo: Expect expression. %+v", p.peek()))
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic("todo: implement error handling " + message)
}
