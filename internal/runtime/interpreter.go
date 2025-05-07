package runtime

import (
	"reflect"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func Init(ast any, env *enviroment.Enviroment) any {
	std.Load(env)
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
	case *parser.Module:
		return EvalModuleStatement(node, env)
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
		return EvalFunctionDeclaration(node, env, nil)
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

func EvalModuleStatement(node *parser.Module, env *enviroment.Enviroment) any {
	scope := enviroment.NewGlobalEnv()
	std.Load(scope)
	var lastResult any
	if module := enviroment.Modules.Get(node.Path); module != nil {
		return env.DeclareVariable(node.Name, module, -1)
	}
	for _, stmt := range node.Body {
		lastResult = Interpreter(stmt, scope)
	}
	env.DeclareVariable(node.Name, scope.Variables, -1)
	enviroment.Modules.Add(node.Path, scope.Variables)
	return lastResult
}

func EvalObjectExpression(node *parser.ObjectExpression, env *enviroment.Enviroment) any {
	caller := node.Caller.(*parser.CallStatement)
	var methods []*parser.FunctionDeclaration
	var className string

	switch t := caller.Caller.(type) {
	case *parser.IdentifierStatement:
		className = t.Value.(string)
		methods = env.GetVariable(className, -1).(*parser.ClassDeclaration).Methods
	case *parser.MemberExpression:
		className = t.Prop.(*parser.IdentifierStatement).Value.(string)
		methods = Interpreter(t, env).(*parser.ClassDeclaration).Methods
	}
	scope := enviroment.NewEnv(env)
	obj := types.NewObject(className, scope)
	EvalFunctionDeclaration(&parser.FunctionDeclaration{
		Name: "init",
		Statement: &parser.Statement{
			Line: -1,
			Kind: parser.FunctionDeclarationNode,
		},
		Body:   []any{},
		Params: []*parser.IdentifierStatement{},
	}, scope, obj)
	for _, method := range methods {
		EvalFunctionDeclaration(method, scope, obj)
	}

	EvalCallStatement(&parser.CallStatement{
		Caller: &parser.IdentifierStatement{
			Value: "init",
		},
		Args: caller.Args,
	}, scope)
	return obj

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

	} else if node.Assign != nil {
		left := Interpreter(node.Left, env).(*types.ObjectValue)
		name := node.Prop.(*parser.IdentifierStatement).Value.(string)
		return left.Enviroment.DeclareVariable(name, node.Assign, -1)
	} else {
		left := Interpreter(node.Left, env)
		switch t := left.(type) {
		case *types.ObjectValue:
			name := node.Prop.(*parser.IdentifierStatement).Value.(string)
			return t.Enviroment.GetVariable(name, -1)
		default:
			return t.(map[string]any)[node.Prop.(*parser.IdentifierStatement).Value.(string)]
		}

		// v := reflect.ValueOf(left)
		// return v.MethodByName(string(strings.ToUpper(name[:1]) + name[1:]))
	}
}

func EvalAssignmentExpression(node *parser.AssignmentExpression, env *enviroment.Enviroment) any {
	switch t := node.Owner.(type) {
	case *parser.MemberExpression:
		t.Assign = Interpreter(node.Value, env)
		Interpreter(t, env)
	default:
		env.AssignmenVariable(node.Owner.(*parser.IdentifierStatement).Value.(string), Interpreter(node.Value, env), node.Statement.Line)
	}
	return nil
}

func EvalForStatement(node *parser.ForStatement, env *enviroment.Enviroment) any {
	// scope := enviroment.NewEnv(env) NOTE: for uchun scope yaratilsa condition xato ishlamoqda to'g'irlash kerak
	for Interpreter(node.Condition, env).(*types.RuntimeValue).Value.(bool) {
		for _, statement := range node.Body {
			result := Interpreter(statement, env)
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

func EvalFunctionDeclaration(node *parser.FunctionDeclaration, env *enviroment.Enviroment, ownerObject any) any {
	fn := &types.FunctionDeclaration{
		Type:        types.Function,
		Body:        node.Body,
		Params:      node.Params,
		OwnerObject: ownerObject,
	}
	return env.AssignmenVariable(node.Name, fn, node.Statement.Line)
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
		scope.DeclareVariable("this", v.OwnerObject, -1)
		for index, name := range v.Params {
			scope.DeclareVariable(name.Value.(string), Interpreter(node.Args[index], env), -1)
		}
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
		// default:

	}
	return nil

}
