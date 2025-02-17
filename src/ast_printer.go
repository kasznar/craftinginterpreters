package src

import "fmt"

type AstPrinter struct{}

func (a AstPrinter) VisitAssignExpr(expr AssignExpr) any {
	//TODO implement me
	panic("implement me")
}

func (a AstPrinter) VisitVariableExpr(expr VariableExpr) any {
	//TODO implement me
	panic("implement me")
}

func (a AstPrinter) VisitBinaryExpr(expr BinaryExpr) any {
	return a.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (a AstPrinter) VisitGroupingExpr(expr GroupingExpr) any {
	return a.parenthesize("group", expr.Expression)
}

func (a AstPrinter) VisitLiteralExpr(expr LiteralExpr) any {
	if expr.Value == nil {
		return "nil"
	}

	return fmt.Sprint(expr.Value)
}

func (a AstPrinter) VisitUnaryExpr(expr UnaryExpr) any {
	return a.parenthesize(expr.Operator.lexeme, expr.Right)
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	result := "(" + name

	for i := 0; i < len(exprs); i++ {
		expr := exprs[i]
		result = result + " " + expr.Accept(a).(string)

	}

	result = result + ")"

	return result
}
