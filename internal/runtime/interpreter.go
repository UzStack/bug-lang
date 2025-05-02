package runtime

import (
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func Interpreter(astBody any, env *enviroment.Enviroment) any {
	switch node := astBody.(type) {
	case *parser.NumberLiteralNode:
		return &types.RuntimeValue{
			Type:  types.Number,
			Value: node.Value,
		}
	case *parser.StringLiteralNode:
		return &types.RuntimeValue{
			Type:  types.Number,
			Value: node.Value,
		}
	case *parser.Program:
		return EvalProgram(node, env)
	case *parser.IdentifierStatement:
		return EvalIdentifier(node, env)
	case *parser.CallStatement:
		return CallStatement(node, env)
	case *parser.VariableDeclaration:
		return VariableDeclaration(node, env)
	case *parser.BinaryExpression:
		return EvalBinaryExpression(node, env)
	default:
		// fmt.Printf("Tip: %T", astBody)
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

func EvalBinaryExpression(node *parser.BinaryExpression, env *enviroment.Enviroment) any {
	left, _ := utils.Str2Int(Interpreter(node.Left, env).(*types.RuntimeValue).Value)
	right, _ := utils.Str2Int(Interpreter(node.Right, env).(*types.RuntimeValue).Value)
	var value any
	switch node.Operator {
	case "+":
		value = left + right
	case "-":
		value = left - right
	case "*":
		value = left * right
	case "/":
		value = left / right
	case "%":
		value = left % right
	}

	return &types.RuntimeValue{
		Type:  types.String,
		Value: value,
	}
}

func VariableDeclaration(node *parser.VariableDeclaration, env *enviroment.Enviroment) any {
	env.DeclareVariable(node.Name, Interpreter(node.Value, env), node.Line)
	return nil
}

func CallStatement(node *parser.CallStatement, env *enviroment.Enviroment) any {
	var args []any
	fn, _ := env.GetVariable(node.Caller.Name, -1).(*types.NativeFunctionDeclaration)
	for _, arg := range node.Args {
		args = append(args, Interpreter(arg, env))
	}
	call := fn.Call.(func(...any))
	call(args...)
	return nil
}
