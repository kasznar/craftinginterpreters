package src

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt)
	VisitPrintStmt(stmt PrintStmt)
}

type ExpressionStmt struct {
	expression Expr
}

func (e ExpressionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	expression Expr
}

func (p PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(p)
}
