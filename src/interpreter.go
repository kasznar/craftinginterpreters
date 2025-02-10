package src

import "fmt"

type Interpreter struct{}

func (i *Interpreter) Interpret(expr Expr) {
	// todo: error handling
	value := i.evaluate(expr)

	fmt.Println("value: ", value)
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
