package runtime

import (
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
	"github.com/k0kubun/pp"
)

func Interpreter(astBody any, env *enviroment.Enviroment) any {
	switch node := astBody.(type) {
	case *parser.NumberLiteral:
		return &types.RuntimeValue{
			Type:  types.Number,
			Value: node.Value,
		}
	case *parser.StringLiteral:
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
	case *parser.FunctionDeclaration:
		return EvalFunctionDeclaration(node, env)
	default:
		// fmt.Printf("Tip: %T", astBody)
	}
	return nil
}

func EvalIdentifier(node *parser.IdentifierStatement, env *enviroment.Enviroment) any {
	name, _ := node.Value.(string)
	return env.GetVariable(name, -1)
}
func EvalFunctionDeclaration(node *parser.FunctionDeclaration, env *enviroment.Enviroment) any {
	fn := &types.FunctionDeclaration{
		Type:   types.Function,
		Body:   node.Body,
		Params: node.Params,
	}
	return env.DeclareVariable(node.Name, fn, node.Statement.Line)
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
	switch v := env.GetVariable(node.Caller.Name, -1).(type) {
	case *types.NativeFunctionDeclaration:
		for _, arg := range node.Args {
			args = append(args, Interpreter(arg, env))
		}
		call := v.Call.(func(...any))
		call(args...)
		return nil
	case *types.FunctionDeclaration:
		var result any
		for _, statement := range v.Body {
			result = Interpreter(statement, env)
		}
		return result
	default:
		pp.Print(v)
	}
	return nil

}
