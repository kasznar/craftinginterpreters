package main

import "fmt"

type Message struct {
	text string
}

func compareStructs() {
	fmt.Println("COMPARE STRUCT VALUES")
	message1 := Message{"hello1"}
	message2 := Message{"hello1"}
	fmt.Println(message1 == message2)
	fmt.Println(&message1 == &message2)
	fmt.Println(Message{"hello"} == Message{"hello"})
}

func mapLookUpByStruct() {
	fmt.Println("MAP LOOKUP BY STRUCT")

	message1 := Message{"hello1"}
	message2 := Message{"hello1"}

	message3 := Message{"hello3"}

	myMap := map[Message]int{}

	myMap[message1] = 1
	// changes the same one
	myMap[message2] = 2

	fmt.Println(myMap[message1])
	fmt.Println(myMap[message2])

	// default value is returned
	fmt.Println(myMap[message3])

	fmt.Println(myMap[Message{"hello1"}])

	fmt.Println(myMap)
}

func mapLookUpByPointer() {
	fmt.Println("MAP LOOKUP BY REFERENCE")

	message1 := &Message{"hello1"}
	message2 := &Message{"hello1"}

	message3 := &Message{"hello3"}

	myMap := map[*Message]int{}

	myMap[message1] = 1
	// no longer changes the same one
	myMap[message2] = 2

	fmt.Println(myMap[message1])
	fmt.Println(myMap[message2])

	// default value is returned
	fmt.Println(myMap[message3])

	fmt.Println(myMap[&Message{"hello1"}])

	fmt.Println(myMap)

}

func slicesShareUnderlyingArray() {
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)

	b[0] = "XXX"
	fmt.Println(a, b)
	fmt.Println(names)
}
func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func lengthAndCapacity() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)
}

func containsKey() {
	m := map[string]bool{"one": true, "two": false}

	_, ok := m["two"]
	fmt.Println("two ", ok)

	_, ok = m["three"]
	fmt.Println("two ", ok)
}

func main() {
	//compareStructs()
	mapLookUpByStruct()
	mapLookUpByPointer()
	//slicesShareUnderlyingArray()
	//lengthAndCapacity()
	//containsKey()
}
