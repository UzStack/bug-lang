package runtime

import (
	"reflect"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
	"github.com/k0kubun/pp"
)

func Init(ast any, env *enviroment.Enviroment) any {
	env.DeclareVariable("print", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: std.Print,
	}, -1)
	env.DeclareVariable("input", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: std.Input,
	}, -1)
	env.DeclareVariable("true", &types.RuntimeValue{
		Type:  "variable",
		Value: true,
	}, -1)
	env.DeclareVariable("false", &types.RuntimeValue{
		Type:  "variable",
		Value: false,
	}, -1)
	env.DeclareVariable("null", &types.RuntimeValue{
		Type:  "variable",
		Value: nil,
	}, -1)
	return Interpreter(ast, env)
}

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
		return EvalCallStatement(node, env)
	case *parser.VariableDeclaration:
		return VariableDeclaration(node, env)
	case *parser.BinaryExpression:
		return EvalBinaryExpression(node, env)
	case *parser.FunctionDeclaration:
		return EvalFunctionDeclaration(node, env)
	case *parser.IfStatement:
		return EvalIfStatement(node, env)
	case *parser.ElseIfStatement:
		return EvalElseIfStatement(node, env)
	case *parser.ElseStatement:
		return EvalElseStatement(node, env)
	case *parser.ForStatement:
		return EvalForStatement(node, env)
	case *parser.AssignmentExpression:
		return EvalAssignmentExpression(node, env)
	case *parser.ReturnStatement:
		return EvalReturnStatement(node, env)
	case *parser.MemberExpression:
		return EvalMemberExpression(node, env)
	case *parser.ArrayExpression:
		return EvalArrayExpression(node, env)
	case *parser.ClassDeclaration:
		return EvalClassDeclaration(node, env)
	case *parser.MapExpression:
		return EvalMapDeclaration(node, env)
	case *parser.ObjectExpression:
		return EvalObjectExpression(node, env)
	default:
		// fmt.Printf("Tip: %T", astBody)
	}
	return nil
}

func EvalObjectExpression(node *parser.ObjectExpression, env *enviroment.Enviroment) any {
	scope := enviroment.NewEnv(env)
	className := node.Name.(*lexar.Token).Value.(string)
	for _, method := range env.GetVariable(className, -1).(*parser.ClassDeclaration).Methods {
		EvalFunctionDeclaration(method, scope)
	}
	return types.NewObject(className, scope)
}

func EvalClassDeclaration(node *parser.ClassDeclaration, env *enviroment.Enviroment) any {
	env.DeclareVariable(node.Name.(*lexar.Token).Value.(string), node, node.Line)
	return nil
}

func VariableDeclaration(node *parser.VariableDeclaration, env *enviroment.Enviroment) any {
	env.DeclareVariable(node.Name, Interpreter(node.Value, env), node.Line)
	return nil
}

func EvalMapDeclaration(node *parser.MapExpression, env *enviroment.Enviroment) any {
	values := make(map[string]any)
	for key, item := range node.Values {
		values[key] = Interpreter(item, env)
	}
	// return types.NewObject(values)
	return nil
}
func EvalArrayExpression(node *parser.ArrayExpression, env *enviroment.Enviroment) any {
	var values []any
	for _, item := range node.Values {
		values = append(values, Interpreter(item, env))
	}
	return &types.ArrayValue{
		Values: values,
	}
}

func EvalReturnStatement(node *parser.ReturnStatement, env *enviroment.Enviroment) any {
	return &types.ReturnValue{
		Value: Interpreter(node.Value, env),
	}
}
func EvalMemberExpression(node *parser.MemberExpression, env *enviroment.Enviroment) any {
	if node.Computed {
		switch t := Interpreter(node.Left, env).(type) {
		case *types.ArrayValue:
			index, _ := utils.Str2Int(Interpreter(node.Prop, env).(*types.RuntimeValue).Value)
			return t.Values[index]
		// case *types.ObjectValue:
		// return t.Values[Interpreter(node.Prop, env).(*types.RuntimeValue).Value.(string)]
		default:
			return nil
		}

	} else {
		left := Interpreter(node.Left, env).(*types.ObjectValue)
		name := node.Prop.(*parser.IdentifierStatement).Value.(string)
		return left.Enviroment.GetVariable(name, -1)
		// v := reflect.ValueOf(left)
		// return v.MethodByName(string(strings.ToUpper(name[:1]) + name[1:]))
	}
}

func EvalAssignmentExpression(node *parser.AssignmentExpression, env *enviroment.Enviroment) any {
	env.AssignmenVariable(node.Owner.(*parser.IdentifierStatement).Value.(string), Interpreter(node.Value, env), node.Statement.Line)
	return nil
}

func EvalForStatement(node *parser.ForStatement, env *enviroment.Enviroment) any {
	scope := enviroment.NewEnv(env)
	for Interpreter(node.Condition, env).(*types.RuntimeValue).Value.(bool) {
		for _, statement := range node.Body {
			result := Interpreter(statement, scope)
			if isReturn, response := IsReturn(result); isReturn {
				return response
			}
		}
	}
	return nil
}

func EvalIfStatement(node *parser.IfStatement, env *enviroment.Enviroment) any {
	if Interpreter(node.Condition, env).(*types.RuntimeValue).Value.(bool) {
		for _, statement := range node.Body {
			result := Interpreter(statement, env)
			if isReturn, response := IsReturn(result); isReturn {
				return response
			}
		}
		return &types.FlowValue{
			Catched: true,
			Type:    types.Flow,
		}
	}
	for _, child := range node.Childs {
		result := Interpreter(child, env)
		if result.(*types.FlowValue).Catched {
			return &types.FlowValue{
				Catched: true,
				Type:    types.Flow,
			}
		}
	}
	return &types.FlowValue{
		Catched: false,
		Type:    types.Flow,
	}
}

func EvalElseIfStatement(node *parser.ElseIfStatement, env *enviroment.Enviroment) any {
	if Interpreter(node.Condition, env).(*types.RuntimeValue).Value.(bool) {
		for _, statement := range node.Body {
			result := Interpreter(statement, env)
			if isReturn, response := IsReturn(result); isReturn {
				return response
			}
		}
		return &types.FlowValue{
			Type:    types.Flow,
			Catched: true,
		}
	}
	return &types.FlowValue{
		Type:    types.Flow,
		Catched: false,
	}
}
func EvalElseStatement(node *parser.ElseStatement, env *enviroment.Enviroment) any {
	for _, statement := range node.Body {
		result := Interpreter(statement, env)
		if isReturn, response := IsReturn(result); isReturn {
			return response
		}
	}
	return &types.FlowValue{
		Catched: true,
		Type:    types.Flow,
	}
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
	var value any
	if utils.InArray(node.Operator, []any{"&&", "||"}) {
		left := Interpreter(node.Left, env).(*types.RuntimeValue).Value.(bool)
		right := Interpreter(node.Right, env).(*types.RuntimeValue).Value.(bool)
		switch node.Operator {
		case "&&":
			value = left && right
		case "||":
			value = left || right
		}
	} else {
		left, _ := utils.Str2Int(Interpreter(node.Left, env).(*types.RuntimeValue).Value)
		right, _ := utils.Str2Int(Interpreter(node.Right, env).(*types.RuntimeValue).Value)
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
		case "==":
			value = left == right
		case ">=":
			value = left >= right
		case "<=":
			value = left <= right
		case "!=":
			value = left != right
		case ">":
			value = left > right
		case "<":
			value = left < right
		}
	}

	return &types.RuntimeValue{
		Type:  types.String,
		Value: value,
	}
}

func IsReturn(result any) (bool, any) {
	switch val := result.(type) {
	case *types.ReturnValue:
		return true, val
	}
	return false, nil
}

func EvalCallStatement(node *parser.CallStatement, env *enviroment.Enviroment) any {
	scope := enviroment.NewEnv(env)
	var args []any
	for _, arg := range node.Args {
		args = append(args, Interpreter(arg, scope))
	}
	switch v := Interpreter(node.Caller, scope).(type) {
	case *types.NativeFunctionDeclaration:
		fun := reflect.ValueOf(v.Call)
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
		}
		out := fun.Call(callArgs)
		var results = make([]any, len(out))
		for i, res := range out {
			results[i] = res.Interface()
		}
		return results
	case *types.FunctionDeclaration:
		var result any
		for _, statement := range v.Body {
			result = Interpreter(statement, scope)
			if isReturn, response := IsReturn(result); isReturn {
				return response.(*types.ReturnValue).Value
			}
		}
		return result
	case reflect.Value:
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
		}
		return v.Call(callArgs)
	default:
		pp.Print(v)
	}
	return nil

}
