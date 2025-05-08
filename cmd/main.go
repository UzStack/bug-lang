package main

import (
	"os"
	"runtime/pprof"
	"time"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/k0kubun/pp/v3"
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
	start := time.Now()
	tokenize := lexar.NewTokenize()
	tokens := tokenize.Tokenize(string(code))
	pp.Println(time.Since(start).Milliseconds())
	parser := parser.NewParser(tokens)
	ast := parser.CreateAST()
	pp.Println(time.Since(start).Milliseconds())
	env := enviroment.NewGlobalEnv()
	runtime.Init(ast, env)
	pp.Println(time.Since(start).Milliseconds())
}
