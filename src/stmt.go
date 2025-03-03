package src

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt)
	VisitPrintStmt(stmt PrintStmt)
	VisitVarStmt(stmt VarStmt)
	VisitBlockStmt(stmt BlockStmt)
	VisitIfStmt(stmt IfStmt)
	VisitWhileStmt(stmt WhileStmt)
	VisitFunctionStmt(stmt FunctionStmt)
	VisitReturnStmt(stmt ReturnStmt)
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

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch *Stmt
}

func (s IfStmt) Accept(visitor StmtVisitor) {
	visitor.VisitIfStmt(s)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) Accept(visitor StmtVisitor) {
	visitor.VisitWhileStmt(s)
}

type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func (s FunctionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitFunctionStmt(s)
}

type ReturnStmt struct {
	keyword Token
	// todo: optional value?
	value Expr
}

func (s ReturnStmt) Accept(visitor StmtVisitor) {
	visitor.VisitReturnStmt(s)
}
