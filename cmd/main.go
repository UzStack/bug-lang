package main

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
)

func main() {
	tokenize := lexar.NewTokenize()
	tokens := tokenize.Tokenize(`
	var age = 20;
	print(age(name));
	`)
	parser := parser.NewParser(tokens)
	parser.CreateAST()
	fmt.Print("BugLang Forever\n")
}
