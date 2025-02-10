package src

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

func (l *Lox) RunFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Something ain't right", err)
		return
	}
	s := string(file)
	r := []rune(s)
	l.run(r)

	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		r := []rune(text)
		l.run(r)
		l.hadError = false
	}
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Printf("[line %s] Error %s : %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) run(source []rune) {
	scanner := makeScanner(source)
	tokens := scanner.scanTokens()
	parser := makeParser(tokens)
	expression, err := parser.parse()

	if l.hadError || err != nil {
		// todo: use report
		println(err)
		return
	}

	println(expression)
	println(AstPrinter{}.Print(expression))

	interpreter := Interpreter{}
	interpreter.Interpret(expression)
}
