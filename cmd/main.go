package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func main() {
	args := os.Args

	// profiling faqat so'ralganda: buglang --profile main.bug
	if len(args) > 1 && args[1] == "--profile" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		args = args[1:]
	}

	// fayl berilmasa — interaktiv REPL
	if len(args) <= 1 {
		repl()
		return
	}

	code, err := os.ReadFile(args[1])
	if err != nil {
		panic(err.Error())
	}
	tokenize := lexar.NewTokenize()
	tokens, err := tokenize.Tokenize(string(code))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	parser := parser.NewParser(tokens, "")
	ast := parser.CreateAST()
	env := enviroment.NewGlobalEnv()
	std.Load(env)
	if _, err := runtime.Interpreter(ast, env); err != nil {
		fmt.Println(err)
	}
}

func repl() {
	env := enviroment.NewGlobalEnv()
	std.Load(env)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("bug-lang REPL — chiqish uchun 'exit' yoki Ctrl+D")
	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			return
		}
		evalLine(line, env)
	}
}

// evalLine bitta REPL qatorini ishlaydi; xato/panic REPLni to'xtatmaydi.
func evalLine(line string, env *enviroment.Enviroment) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("xato:", r)
		}
	}()
	tokenize := lexar.NewTokenize()
	tokens, err := tokenize.Tokenize(line)
	if err != nil {
		fmt.Println("xato:", err.Error())
		return
	}
	ast, ok := parser.NewParser(tokens, "").CreateAST().(*parser.ProgramNode)
	if !ok {
		return
	}
	for _, stmt := range ast.Body {
		res, err := runtime.Interpreter(stmt, env)
		if err != nil {
			fmt.Println("xato:", err.Error())
			return
		}
		// ifoda natijasini chop etamiz (null bo'lmasa)
		if obj, ok := res.(types.Object); ok {
			if _, isNull := obj.(*types.NullValue); !isNull {
				std.Print(obj)
				fmt.Println()
			}
		}
	}
}
