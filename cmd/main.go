package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer func() {
		pprof.StopCPUProfile()
	}()

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
	parser := parser.NewParser(tokens, "")
	ast := parser.CreateAST()
	env := enviroment.NewGlobalEnv()
	std.Load(env)
	if _, err := runtime.Interpreter(ast, env); err != nil {
		fmt.Println(err)
	}
}
