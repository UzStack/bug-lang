package main

import (
	"os"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
)

func main() {
	code, err := os.ReadFile("example.bug")
	if err != nil {
		panic(err.Error())
	}
	tokenize := lexar.NewTokenize()
	tokens := tokenize.Tokenize(string(code))
	parser := parser.NewParser(tokens)
	ast := parser.CreateAST()
	env := enviroment.NewGlobalEnv()
	runtime.Interpreter(ast, env)
}
