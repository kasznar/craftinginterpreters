package src

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(expr BinaryExpr) any
	VisitGroupingExpr(expr GroupingExpr) any
	VisitLiteralExpr(expr LiteralExpr) any
	VisitUnaryExpr(expr UnaryExpr) any
	VisitVariableExpr(expr VariableExpr) any
	VisitAssignExpr(expr AssignExpr) any
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func (g GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (e LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name Token
}

func (e VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (e AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(e)
}
