package runtime

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func Interpreter(astBody any, env *enviroment.Enviroment) any {
	switch node := astBody.(type) {
	case *parser.Program:
		return EvalProgram(node, env)
	case *parser.IdentifierStatement:
		return EvalIdentifier(node, env)
	case *parser.CallStatement:
		return CallStatement(node, env)
	case *parser.VariableDeclaration:
		return VariableDeclaration(node, env)
	default:
		fmt.Printf("Tip: %T", astBody)
	}
	return nil
}

func EvalIdentifier(node *parser.IdentifierStatement, env *enviroment.Enviroment) any {
	name, _ := node.Value.(string)
	return env.GetVariable(name, -1)
}

func EvalProgram(node *parser.Program, env *enviroment.Enviroment) any {
	var lastInterpreted any
	for _, statement := range node.Body {
		lastInterpreted = Interpreter(statement, env)
	}
	return lastInterpreted
}

func VariableDeclaration(node *parser.VariableDeclaration, env *enviroment.Enviroment) any {
	value, _ := node.Value.(map[string]any)
	env.DeclareVariable(node.Name, value["value"], node.Line)
	return nil
}

func CallStatement(node *parser.CallStatement, env *enviroment.Enviroment) any {
	var args []any
	name, _ := node.Name.(string)
	fn, _ := env.GetVariable(name, -1).(*types.NativeFunctionDeclaration)
	for _, arg := range node.Args {
		args = append(args, Interpreter(arg, env))
	}
	call := fn.Call.(func(...any))
	call(args...)
	return nil
}
