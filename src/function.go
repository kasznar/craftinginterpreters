package src

type LoxCallable interface {
	// Arity arity == number of arguments
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
}

type LoxFunction struct {
	declaration FunctionStmt
}

func (f LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f LoxFunction) Call(interpreter *Interpreter, arguments []any) any {
	environment := NewEnvironment(interpreter.globals)

	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.body, environment)

	return nil
}

func (f LoxFunction) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
