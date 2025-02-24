package main

import "fmt"

type Return struct {
	value any
}

func throwing() {
	panic(&Return{5})
}

func getValue() (returnValue any) {

	defer func() {
		if exception := recover(); exception != nil {
			rv, ok := exception.(*Return)
			if !ok {
				panic("not ok")
			}
			returnValue = rv.value
		}

	}()

	throwing()
	return nil
}

func main() {
	val := getValue()
	val = getValue()
	val = getValue()
	fmt.Println(val)
}
