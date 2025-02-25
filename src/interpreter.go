package src

import "fmt"

type Interpreter struct {
	globals     *Environment
	environment *Environment
}

func NewInterpreter() Interpreter {
	globals := NewEnvironment(nil)

	globals.define("clock", Clock{})

	environment := globals
	return Interpreter{globals, environment}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	// todo: error handling
	for j := 0; j < len(statements); j++ {
		stmt := statements[j]
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment

	i.environment = environment

	defer func() {
		i.environment = previous
	}()

	for j := 0; j < len(statements); j++ {
		statement := statements[j]
		i.execute(statement)
	}

}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}

	if boolean, ok := object.(bool); ok {
		return boolean
	}

	return true
}

func (i *Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	// todo: how go implement this vs java?
	return a == b
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}

	panic(fmt.Errorf("%+v operand must be a number", operator))
}

func (i *Interpreter) checkNumberOperands(operator Token, left any, right any) {
	if _, leftOk := left.(float64); leftOk {
		if _, rightOk := right.(float64); rightOk {
			return
		}
	}

	panic(fmt.Errorf("%+v operands must be a numbers", operator))
}

func (i *Interpreter) VisitExpressionStmt(stmt ExpressionStmt) {
	i.evaluate(stmt.expression)
}

func (i *Interpreter) VisitPrintStmt(stmt PrintStmt) {
	value := i.evaluate(stmt.expression)
	fmt.Println(value)
}

func (i *Interpreter) VisitReturnStmt(stmt ReturnStmt) {
	var value any = nil
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}

	panic(&Return{value})
}

func (i *Interpreter) VisitVarStmt(stmt VarStmt) {
	var value any

	// todo: need to dereference otherwise not going to work
	// https://go.dev/doc/faq#nil_error
	if *stmt.initializer != nil {
		value = i.evaluate(*stmt.initializer)
	}

	i.environment.define(stmt.name.lexeme, value)
}

func (i *Interpreter) VisitAssignExpr(expr AssignExpr) any {
	value := i.evaluate(expr.Value)
	i.environment.assign(expr.Name, value)
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64)
	case PLUS:
		leftNum, leftOk := left.(float64)
		rightNum, rightOk := right.(float64)

		if leftOk && rightOk {
			return leftNum + rightNum
		}

		// note: more concise if statement
		if leftStr, leftOk := left.(string); leftOk {
			if rightStr, rightOk := right.(string); rightOk {
				return leftStr + rightStr
			}
		}

		panic(fmt.Errorf("%+v operands must be two numbers or two strings", expr.Operator))
	case SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.tokenType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr VariableExpr) any {
	return i.environment.get(expr.Name)
}

func (i *Interpreter) VisitBlockStmt(stmt BlockStmt) {
	blockEnv := NewEnvironment(i.environment)
	i.executeBlock(stmt.statements, blockEnv)
}

func (i *Interpreter) VisitIfStmt(stmt IfStmt) {
	if i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if *stmt.elseBranch != nil {
		i.execute(*stmt.elseBranch)
	}
}

func (i *Interpreter) VisitLogicalExpr(expr LogicalExpr) any {
	left := i.evaluate(expr.left)

	if expr.operator.tokenType == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) VisitWhileStmt(stmt WhileStmt) {
	for i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
}

func (i *Interpreter) VisitCallExpr(expr CallExpr) any {
	callee := i.evaluate(expr.callee)

	arguments := make([]any, 0)

	for _, argument := range expr.arguments {
		arguments = append(arguments, i.evaluate(argument))
	}
	function, ok := callee.(LoxCallable)

	if !ok {
		panic(fmt.Errorf("can only call functions and classes. %+v", expr.paren))
	}

	if len(arguments) != function.Arity() {
		panic(fmt.Errorf("wrong number of arguments"))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitFunctionStmt(stmt FunctionStmt) {
	function := LoxFunction{stmt}
	i.environment.define(stmt.name.lexeme, function)
}
