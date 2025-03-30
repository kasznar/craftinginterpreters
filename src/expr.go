package src

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitVariableExpr(expr *VariableExpr) any
	VisitAssignExpr(expr *AssignExpr) any
	VisitLogicalExpr(expr *LogicalExpr) any
	VisitCallExpr(expr *CallExpr) any
	VisitGetExpr(expr *GetExpr) any
	VisitSetExpr(expr *SetExpr) any
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name Token
}

func (e *VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (e *AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(e)
}

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (e *LogicalExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpr(e)
}

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (e *CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpr(e)
}

type GetExpr struct {
	object Expr
	name   Token
}

func (e *GetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGetExpr(e)
}

type SetExpr struct {
	object Expr
	name   Token
	value  Expr
}

func (e *SetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetExpr(e)
}
