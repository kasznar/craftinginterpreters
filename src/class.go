package src

import (
	"fmt"
)

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

// Arity todo: pointer receiver or no pointer receiver
func (c LoxClass) Arity() int {
	initializer := c.findMethod("init")

	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}

func (c LoxClass) Call(interpreter *Interpreter, arguments []any) any {
	instance := &LoxInstance{c, make(map[string]any)}
	initializer := c.findMethod("init")

	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (c LoxClass) String() string {
	return c.name
}

func (c LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil
}

type LoxInstance struct {
	class  LoxClass
	fields map[string]any
}

func (i LoxInstance) String() string {
	return i.class.name + " instance"
}

func (i *LoxInstance) get(name Token) any {
	if value, ok := i.fields[name.lexeme]; ok {
		return value
	}

	method := i.class.findMethod(name.lexeme)

	if method != nil {
		return method.Bind(i)
	}

	panic(fmt.Errorf("Undefined property '" + name.lexeme + "'."))
}

func (i LoxInstance) set(name Token, value any) {
	i.fields[name.lexeme] = value
}
