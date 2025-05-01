package main

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/lexar"
)

func main() {
	tokenize := lexar.NewTokenize()
	tokenize.Tokenize("age := 20")
	fmt.Print(tokenize.Get())
	fmt.Print("BugLang Forever\n")
}
