package src

type LoxCallable interface {
	// Arity arity == number of arguments
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
}

type LoxFunction struct {
	declaration   *FunctionStmt
	closure       *Environment
	isInitializer bool
}

func (f LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.closure)
	environment.define("this", instance)
	return &LoxFunction{f.declaration, environment, f.isInitializer}
}

func (f LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f LoxFunction) Call(interpreter *Interpreter, arguments []any) (returnValue any) {
	environment := NewEnvironment(f.closure)

	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}

	defer func() {
		if exception := recover(); exception != nil {
			if f.isInitializer {
				returnValue = f.closure.getAt(0, "this")
				return
			}

			rv := exception.(*Return)
			returnValue = rv.value
		}
	}()

	interpreter.executeBlock(f.declaration.body, environment)

	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}
	return nil
}

func (f LoxFunction) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}

type Return struct {
	value any
}
