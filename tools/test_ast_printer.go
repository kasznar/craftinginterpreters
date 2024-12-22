package main

import . "craftinginterpreters/src"

func main() {
	expression := BinaryExpr{
		Left: UnaryExpr{
			NewToken(MINUS, "-", nil, 1),
			LiteralExpr{123},
		},
		Operator: NewToken(STAR, "*", nil, 1),
		Right:    GroupingExpr{LiteralExpr{45.67}},
	}

	println(AstPrinter{}.Print(expression))
}
