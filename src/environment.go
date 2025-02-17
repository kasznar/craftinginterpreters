package src

import "fmt"

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment {
	values := map[string]any{}

	return Environment{values}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) any {
	value, has := e.values[name.lexeme]

	if has {
		return value
	}

	panic(fmt.Errorf(name.String(), "Undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) assign(name Token, value any) {

	if _, has := e.values[name.lexeme]; has {
		e.values[name.lexeme] = value
		return
	}

	panic(fmt.Errorf(name.String(), "Undefined variable '"+name.lexeme+"'."))

}
