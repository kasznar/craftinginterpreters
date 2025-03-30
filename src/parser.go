package src

import (
	"fmt"
	"reflect"
)

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
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("functions")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")
	p.consume(LEFT_BRACE, "Expect '{' before class body.")

	var methods []*FunctionStmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(RIGHT_BRACE, "Expect '}' after class body.")

	return &ClassStmt{name, methods}
}

func (p *Parser) function(kind string) *FunctionStmt {
	name := p.consume(IDENTIFIER, "Expect "+kind+" name.")

	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	parameters := make([]Token, 0)

	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				panic("too many parameters")
			}

			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))

			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()

	return &FunctionStmt{name, parameters, body}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable Name")

	// todo: nil litral?
	var initializer Expr

	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	// todo: how is this returning Stmt instead of *Stmt??
	return &VarStmt{initializer: &initializer, name: name}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return &WhileStmt{condition, body}
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
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return &BlockStmt{statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	// todo: what is the default value here?
	var initializer *Stmt

	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		varDec := p.varDeclaration()
		initializer = &varDec
	} else {
		exprStmt := p.expressionStatement()
		initializer = &exprStmt
	}

	var condition *Expr

	if !p.check(SEMICOLON) {
		conditionExpr := p.expression()
		condition = &conditionExpr
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment *Expr
	if !p.check(RIGHT_PAREN) {
		incrementExpr := p.expression()
		increment = &incrementExpr
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &BlockStmt{
			statements: []Stmt{
				body,
				&ExpressionStmt{expression: *increment},
			},
		}
	}

	if condition == nil {
		literalExpr := LiteralExpr{Value: true}
		// todo: double check this
		casted := reflect.ValueOf(literalExpr).Interface().(Expr)
		condition = &casted
	}

	body = &WhileStmt{condition: *condition, body: body}

	if initializer != nil {
		body = &BlockStmt{
			statements: []Stmt{
				*initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt

	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{condition, thenBranch, &elseBranch}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ;, after value")
	return &PrintStmt{expression: value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value *Expr = nil

	if !p.check(SEMICOLON) {
		expr := p.expression()
		value = &expr
	}

	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &ReturnStmt{keyword, *value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &ExpressionStmt{expression: expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if target, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{Name: target.Name, Value: value}
		} else if target, ok := expr.(*GetExpr); ok {
			return &SetExpr{target.object, target.name, value}
		}

		// todo: error handling
		panic(fmt.Errorf(equals.String(), "Invalid assignment target."))
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()

		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()

		return &UnaryExpr{Operator: operator, Right: right}
	}

	return p.call()
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)

	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				panic(fmt.Errorf("can't have more than 255 arguments. %+v", p.peek()))
			}

			arguments = append(arguments, p.expression())

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return &CallExpr{callee, paren, arguments}
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = &GetExpr{expr, name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{Value: false}
	}
	if p.match(TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{Value: p.previous().literal}
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{Name: p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{Expression: expr}
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
