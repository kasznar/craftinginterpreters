package src

import "fmt"

type Scope map[string]bool

type FunctionType int

const (
	NONE_FUNCTION FunctionType = iota
	FUNCTION
	METHOD
	INITIALIZER
)

type ClassType int

const (
	NOT_CLASS ClassType = iota
	IS_CLASS
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []Scope
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []Scope{},
		currentFunction: NONE_FUNCTION,
		currentClass:    NOT_CLASS,
	}
}

func (r *Resolver) peekScope() Scope {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveStmtList(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(Scope))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.peekScope()
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.peekScope()[name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *FunctionStmt, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmtList(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) any {
	if len(r.scopes) > 0 {
		if val, ok := r.peekScope()[expr.Name.lexeme]; ok && !val {
			panic("Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) any {
	r.resolveExpr(expr.callee)

	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(expr *GetExpr) any {
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) any {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) any {
	if r.currentClass == NOT_CLASS {
		panic(fmt.Errorf(expr.keyword.lexeme + " Can't use 'this' outside of a class."))
	}

	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExpressionStmt) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) {
	r.declare(stmt.name)
	if *stmt.initializer != nil {
		// todo: pointer good?
		r.resolveExpr(*stmt.initializer)
	}
	r.define(stmt.name)
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) {
	r.beginScope()
	r.resolveStmtList(stmt.statements)
	r.endScope()
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(*stmt.elseBranch)
	}
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) {
	if r.currentFunction == NONE_FUNCTION {
		panic("Can't return from top-level code.")
	}

	// todo: optional value?
	if stmt.value != nil {
		if r.currentFunction == INITIALIZER {
			panic(fmt.Errorf("can't return a value from an initializer"))
		}
		r.resolveExpr(*stmt.value)
	}
}

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) {
	enclosingClass := r.currentClass
	r.currentClass = IS_CLASS

	r.declare(stmt.name)
	r.define(stmt.name)

	r.beginScope()
	r.peekScope()["this"] = true

	for _, method := range stmt.methods {
		declaration := METHOD
		if method.name.lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()
	r.currentClass = enclosingClass
}
