package main

import (
	"os"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/k0kubun/pp"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		panic("Fayil kiritilmadi")
	}
	code, err := os.ReadFile(args[1])
	if err != nil {
		panic(err.Error())
	}
	tokenize := lexar.NewTokenize()
	tokens := tokenize.Tokenize(string(code))
	parser := parser.NewParser(tokens)
	ast := parser.CreateAST()
	pp.Print(ast)
	env := enviroment.NewGlobalEnv()
	runtime.Interpreter(ast, env)
}
