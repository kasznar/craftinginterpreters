package main

import (
	"fmt"
	"os"
	"strings"
)

func defineType(file os.File, baseName string, structName string, fieldList string) {}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + baseName + ".go"

	err := os.MkdirAll(outputDir, os.ModePerm)

	if err != nil {
		panic(err)
	}

	file, err := os.Create(path) // Create or truncate file

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString("package main\n") // Write content
	file.WriteString("\nWelcome to Go!")

	for i := 0; i < len(types); i++ {
		aType := types[i]

		name := strings.Split(aType, ":")[0]
		fields := strings.Split(aType, ":")[0]

		line := fmt.Sprintf("%s, %s", name, fields)

		file.WriteString(line)
	}
}

/*func main() {
	types := []string{"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Unary    : Token operator, Expr right"}

	defineAst("./generated", "expr", types)
}*/
