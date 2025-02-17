package src

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt)
	VisitPrintStmt(stmt PrintStmt)
	VisitVarStmt(stmt VarStmt)
	VisitBlockStmt(stmt BlockStmt)
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

func (s PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(s)
}

type VarStmt struct {
	name        Token
	initializer *Expr
}

func (s VarStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarStmt(s)
}

type BlockStmt struct {
	statements []Stmt
}

func (s BlockStmt) Accept(visitor StmtVisitor) {
	visitor.VisitBlockStmt(s)
}
