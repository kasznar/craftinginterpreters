package src

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    map[string]any{},
		enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) any {
	value, has := e.values[name.lexeme]

	if has {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(fmt.Errorf(name.String(), "Undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) assign(name Token, value any) {

	if _, has := e.values[name.lexeme]; has {
		e.values[name.lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	panic(fmt.Errorf(name.String(), "Undefined variable '"+name.lexeme+"'."))
}

func (e *Environment) assignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e

	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}

	return environment
}
