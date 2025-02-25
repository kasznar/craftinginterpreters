package main

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

func main() {
	one := NewEnvironment(nil)
	two := one

	one.values["hello"] = "hello"

	fmt.Println(two)
}
