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

func (p *Parser) parse() (statements []Stmt, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("syntax error", r)
		}
	}()

	// todo: what happens if this initializer is removed?
	statements = make([]Stmt, 0)

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements, nil
}

func (p *Parser) declaration() Stmt {
	// todo: error handling

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable Name")

	var initializer Expr

	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return VarStmt{initializer: &initializer, name: name}
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return BlockStmt{statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ;, after value")
	return PrintStmt{expression: value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return ExpressionStmt{expression: expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

// todo add here assignment
func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if target, ok := expr.(VariableExpr); ok {
			return AssignExpr{Name: target.Name, Value: value}
		}

		// todo: error handling
		panic(fmt.Errorf(equals.String(), "Invalid assignment target."))
	}

	return expr
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

	if p.match(IDENTIFIER) {
		return VariableExpr{Name: p.previous()}
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
