package runtime

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
	"github.com/k0kubun/pp"
)

func Interpreter(astBody any, env *enviroment.Enviroment) (any, error) {
	switch node := astBody.(type) {
	case *parser.NumberLiteral:
		value, err := strconv.Atoi(node.Value.(string))
		if err != nil {
			fmt.Println("Type error: ", node.Value, "not integer")
			os.Exit(1)
		}
		return types.NewInt(value), nil
	case *parser.StringLiteral:
		value, err := node.Value.(string)
		if !err {
			fmt.Println("Type not string", value)
			os.Exit(1)
		}
		return types.NewString(value), nil
	case *parser.FloatLiteral:
		value, err := strconv.ParseFloat(node.Value.(string), 64)
		if err != nil {
			fmt.Println("Type error: ", node.Value, "not float")
			os.Exit(1)
		}
		return types.NewFloat(value), nil
	case *parser.Module:
		return EvalModuleStatement(node, env)
	case *parser.StdModule:
		return EvalStdModuleStatement(node, env)
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
		return node, nil
		// fmt.Printf("Tip: %T", astBody)
	}
}

func EvalBody(statements []any, env *enviroment.Enviroment) (any, error, bool) {
	for _, statement := range statements {
		res, err := Interpreter(statement, env)
		if err != nil {
			return nil, err, false
		}
		if isReturn, response := IsReturn(res); isReturn {
			return response, nil, true
		}
	}
	return nil, nil, false
}

func EvalStdModuleStatement(node *parser.StdModule, env *enviroment.Enviroment) (any, error) {
	env.DeclareVariable(node.Name, types.NewStdLib(node.Name, std.STDLIBS[node.Path]), node.Line)
	return nil, nil
}
func EvalModuleStatement(node *parser.Module, env *enviroment.Enviroment) (any, error) {
	scope := enviroment.NewEnv(nil)
	std.Load(scope)
	if module := enviroment.Modules.Get(node.Path); module != nil {
		return env.DeclareVariable(node.Name, module, node.Line), nil
	}
	for _, stmt := range node.Body {
		if _, err := Interpreter(stmt, scope); err != nil {
			return nil, err
		}
	}
	env.DeclareVariable(node.Name, types.NewModule(scope), node.Line)
	enviroment.Modules.Add(node.Path, types.NewModule(scope))
	return nil, nil
}

func DeclareExtends(class *parser.ClassDeclaration, env *enviroment.Enviroment, scope *enviroment.Enviroment, obj types.Object) (map[string]*enviroment.Enviroment, error) {
	envs := make(map[string]*enviroment.Enviroment)
	for _, extend := range class.Extends {
		extScope := enviroment.NewEnv(scope)
		res, err := Interpreter(extend, env)
		if err != nil {
			return nil, err
		}
		exd := res.(*parser.ClassDeclaration)
		DeclareExtends(exd, env, scope, obj)
		for _, method := range exd.Methods {
			EvalFunctionDeclaration(method, extScope, obj)
		}
		for key, value := range extScope.Variables {
			scope.Variables[key] = value
		}
		envs[exd.Name.(*parser.IdentifierStatement).Value.(string)] = extScope
	}
	return envs, nil
}

func EvalObjectExpression(node *parser.ObjectExpression, env *enviroment.Enviroment) (any, error) {
	caller := node.Caller.(*parser.CallStatement)
	var methods []*parser.FunctionDeclaration
	var className string
	var class *parser.ClassDeclaration
	switch t := caller.Caller.(type) {
	case *parser.IdentifierStatement:
		className = t.Value.(string)
		res, err := env.GetVariable(className, node.Line)
		if err != nil {
			return nil, err
		}
		class = res.(*parser.ClassDeclaration)
		methods = class.Methods
	case *parser.MemberExpression:
		className = t.Prop.(*parser.IdentifierStatement).Value.(string)
		res, err := Interpreter(t, env)
		if err != nil {
			return nil, err
		}
		class = res.(*parser.ClassDeclaration)
		methods = class.Methods
	}
	scope := enviroment.NewEnv(class.Enviroment)
	obj := types.NewObject(className, scope).(*types.ObjectValue)
	EvalFunctionDeclaration(&parser.FunctionDeclaration{
		Name:   "init",
		Line:   node.Line,
		Body:   []any{},
		Params: []any{},
	}, scope, obj)
	extends, err := DeclareExtends(class, env, scope, obj)
	if err != nil {
		return nil, err
	}
	obj.Extends = extends
	for _, method := range methods {
		EvalFunctionDeclaration(method, scope, obj)
	}
	_, err = EvalCallStatement(&parser.CallStatement{
		Caller: &parser.IdentifierStatement{
			Value: "init",
		},
		Args: caller.Args,
	}, scope)

	if err != nil {
		return nil, err
	}
	return obj, nil

}

func EvalClassDeclaration(node *parser.ClassDeclaration, env *enviroment.Enviroment) (any, error) {
	node.Enviroment = env
	return env.DeclareVariable(node.Name.(*parser.IdentifierStatement).Value.(string), node, node.Line), nil
}

func VariableDeclaration(node *parser.VariableDeclaration, env *enviroment.Enviroment) (any, error) {
	value, err := Interpreter(node.Value, env)
	if err != nil {
		return nil, err
	}
	return env.DeclareVariable(node.Name, value, node.Line), nil
}

func EvalMapDeclaration(node *parser.MapExpression, env *enviroment.Enviroment) (any, error) {
	values := make(map[string]any)

	for key, item := range node.Values {
		value, err := Interpreter(item, env)
		if err != nil {
			return nil, err
		}
		values[key] = value
	}
	return types.NewMap(values), nil
}
func EvalArrayExpression(node *parser.ArrayExpression, env *enviroment.Enviroment) (any, error) {
	var values []any
	for _, item := range node.Values {
		value, err := Interpreter(item, env)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return types.NewArray(values), nil
}

func EvalReturnStatement(node *parser.ReturnStatement, env *enviroment.Enviroment) (any, error) {
	value, err := Interpreter(node.Value, env)
	if err != nil {
		return nil, err
	}
	return &types.ReturnValue{
		Value: value,
	}, nil
}
func EvalMemberExpression(node *parser.MemberExpression, env *enviroment.Enviroment) (any, error) {
	if node.Computed {
		left, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		switch t := left.(type) {
		case *types.ArrayValue:
			prop, err := Interpreter(node.Prop, env)
			if err != nil {
				return nil, err
			}
			index, _ := utils.Str2Int(prop.(types.Object).GetValue())
			return t.Values[index], nil
		// case *types.ObjectValue:
		// return t.Values[Interpreter(node.Prop, env).(types.Object).GetValue().(string)]
		case *types.MapValue:
			prop, err := Interpreter(node.Prop, env)
			if err != nil {
				return nil, err
			}
			name := prop.(types.Object)
			if node.Assign != nil {
				t.Values[name.GetValue().(string)] = node.Assign
			}
			prop, err = Interpreter(node.Prop, env)
			if err != nil {
				return nil, err
			}
			index, _ := prop.(types.Object).GetValue().(string)
			response, ok := t.Values[index]
			if !ok {
				fmt.Printf("Map item not found: %s", index)
			}
			return response, nil
		default:
			return nil, nil
		}

	} else if node.Assign != nil {
		res, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		left := res.(*types.ObjectValue)
		name := node.Prop.(*parser.IdentifierStatement).Value.(string)
		return left.Enviroment.AssignmenVariable(name, node.Assign, node.Line), nil
	} else {
		left, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		prop := node.Prop.(*parser.IdentifierStatement).Value.(string)
		switch t := left.(type) {
		case *types.ObjectValue:
			return t.Enviroment.GetVariable(prop, node.Line)
		case *types.StdLibValue:
			return reflect.ValueOf(t.Lib[prop]), nil
		case *types.ModuleValue:
			return t.Enviroment.GetVariable(prop, node.Line)
		case types.Object:
			v := reflect.ValueOf(left)
			return v.MethodByName(string(strings.ToUpper(prop[:1]) + prop[1:])), nil
		case *enviroment.Enviroment:
			return t.GetVariable(prop, node.Line)
		default:
			return t.(map[string]any)[prop], nil
		}

	}
}

func EvalAssignmentExpression(node *parser.AssignmentExpression, env *enviroment.Enviroment) (any, error) {
	switch t := node.Owner.(type) {
	case *parser.MemberExpression:
		value, err := Interpreter(node.Value, env)
		if err != nil {
			return nil, err
		}
		t.Assign = value
		res, err := Interpreter(t, env)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		value, err := Interpreter(node.Value, env)
		if err != nil {
			return nil, err
		}
		return env.AssignmenVariable(node.Owner.(*parser.IdentifierStatement).Value.(string), value, node.Line), nil
	}
}

func EvalForStatement(node *parser.ForStatement, env *enviroment.Enviroment) (any, error) {
	// scope := enviroment.NewEnv(env) NOTE: for uchun scope yaratilsa condition xato ishlamoqda to'g'irlash kerak
	for {
		res, err := Interpreter(node.Condition, env)
		if err != nil {
			return nil, err
		}
		condition := res.(types.Object).GetValue().(bool)
		if !condition {
			break
		}
		res, err, isReturn := EvalBody(node.Body, env)
		if err != nil {
			return nil, err
		}
		if isReturn {
			return res, nil
		}
	}
	return nil, nil
}

func EvalIfStatement(node *parser.IfStatement, env *enviroment.Enviroment) (any, error) {
	condition, err := Interpreter(node.Condition, env)
	if err != nil {
		return nil, err
	}
	if condition.(types.Object).GetValue().(bool) {
		res, err, isReturn := EvalBody(node.Body, env)
		if err != nil {
			return nil, err
		}
		if isReturn {
			return res, nil
		}
		return &types.FlowValue{
			Catched: true,
			Type:    types.Flow,
		}, nil
	}
	for _, child := range node.Childs {
		child, err := Interpreter(child, env)
		if err != nil {
			return nil, err
		}
		result := child
		if result.(*types.FlowValue).Catched {
			return &types.FlowValue{
				Catched: true,
				Type:    types.Flow,
			}, nil
		}
	}
	return &types.FlowValue{
		Catched: false,
		Type:    types.Flow,
	}, nil
}

func EvalElseIfStatement(node *parser.ElseIfStatement, env *enviroment.Enviroment) (any, error) {
	condition, err := Interpreter(node.Condition, env)
	if err != nil {
		return nil, err
	}
	if condition.(types.Object).GetValue().(bool) {
		res, err, isReturn := EvalBody(node.Body, env)
		if err != nil {
			return nil, err
		}
		if isReturn {
			return res, nil
		}
		return &types.FlowValue{
			Type:    types.Flow,
			Catched: true,
		}, nil
	}
	return &types.FlowValue{
		Type:    types.Flow,
		Catched: false,
	}, nil
}
func EvalElseStatement(node *parser.ElseStatement, env *enviroment.Enviroment) (any, error) {
	res, err, isReturn := EvalBody(node.Body, env)
	if err != nil {
		return nil, err
	}
	if isReturn {
		return res, nil
	}
	return &types.FlowValue{
		Catched: true,
		Type:    types.Flow,
	}, nil
}

func EvalIdentifier(node *parser.IdentifierStatement, env *enviroment.Enviroment) (any, error) {
	name, _ := node.Value.(string)
	return env.GetVariable(name, node.Line)
}

func EvalFunctionDeclaration(node *parser.FunctionDeclaration, env *enviroment.Enviroment, ownerObject any) (any, error) {
	fn := &types.FunctionDeclaration{
		Name:        node.Name,
		Type:        types.Function,
		Body:        node.Body,
		Params:      node.Params,
		OwnerObject: ownerObject,
		Enviroment:  env,
	}
	return env.AssignmenVariable(node.Name, fn, node.Line), nil
}

func EvalProgram(node *parser.Program, env *enviroment.Enviroment) (any, error) {
	for _, statement := range node.Body {
		if _, err := Interpreter(statement, env); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func EvalBinaryExpression(node *parser.BinaryExpression, env *enviroment.Enviroment) (any, error) {
	if utils.InArray(node.Operator, []any{"&&", "||"}) {
		var value bool
		res, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		left := res.(types.Object).GetValue().(bool)
		res, err = Interpreter(node.Right, env)
		if err != nil {
			return nil, err
		}
		right := res.(types.Object).GetValue().(bool)
		switch node.Operator {
		case "&&":
			value = left && right
		case "||":
			value = left || right
		}
		return types.NewBool(value), nil
	} else {
		var left, right float64
		res, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		leftValue := res
		res, err = Interpreter(node.Right, env)
		if err != nil {
			return nil, err
		}
		rightValue := res

		if leftValue != 0 {
			left, _ = utils.Int2Float(leftValue.(types.Object).GetValue())
		}
		if rightValue != 0 {
			right, _ = utils.Int2Float(rightValue.(types.Object).GetValue())
		}
		if utils.InArray(node.Operator, []any{"+", "-", "*", "/"}) {
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
				value = int(left) % int(right)
			}
			_, isRightFloat := rightValue.(*types.FloatValue)
			_, isLeftFloat := leftValue.(*types.FloatValue)
			if isLeftFloat || isRightFloat {
				return types.NewFloat(value.(float64)), nil
			}
			v, _ := utils.Float2Int(value)
			return types.NewInt(v), nil
		} else if utils.InArray(node.Operator, []any{"==", ">=", "<=", "!=", ">", "<"}) {
			var value bool
			switch node.Operator {
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
			return types.NewBool(value), nil
		}

	}
	return nil, nil
}

func IsReturn(result any) (bool, any) {
	switch val := result.(type) {
	case *types.ReturnValue:
		return true, val
	}
	return false, nil
}

func EvalCallStatement(node *parser.CallStatement, env *enviroment.Enviroment) (any, error) {
	scope := enviroment.NewEnv(env)
	var args []any
	for _, arg := range node.Args {
		a, err := Interpreter(arg, scope)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	caller, err := Interpreter(node.Caller, scope)
	if err != nil {
		return nil, err
	}
	switch v := caller.(type) {
	case *types.NativeFunctionValue:
		fun := reflect.ValueOf(v.Call)
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
		}
		out := fun.Call(callArgs)
		if len(out) >= 1 {
			return out[0].Interface(), nil
		}
		return types.NewNull(), nil
	case *types.FunctionDeclaration:
		if len(v.Params) != len(node.Args) {
			return nil, fmt.Errorf("%s funcsiyasi parametri xato berildi funcsiya kutmoqda: %d berildi: %d line: %d", v.Name, len(v.Params), len(node.Args), node.Line)
		}
		var result any
		scope = v.Enviroment
		scope.AssignmenVariable("this", v.OwnerObject, node.Line)
		scope.AssignmenVariable("super", &types.NativeFunctionValue{
			Call: func(value *parser.ClassDeclaration) any {
				return v.OwnerObject.(*types.ObjectValue).Extends[value.Name.(*parser.IdentifierStatement).Value.(string)]
			},
		}, node.Line)
		for index, name := range v.Params {
			param, err := Interpreter(node.Args[index], env)
			if err != nil {
				return nil, err
			}
			scope.AssignmenVariable(name.(*parser.IdentifierStatement).Value.(string), param, node.Line)
		}
		res, err, isReturn := EvalBody(v.Body, scope)
		if err != nil {
			return nil, err
		}
		if isReturn {
			return res.(*types.ReturnValue).Value, nil
		}

		return result, nil
	case reflect.Value:
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
		}
		if res := v.Call(callArgs); len(res) >= 1 {
			return res[0].Interface(), nil
		}
		return nil, nil
	default:
		pp.Print(v)
	}
	return nil, nil
}
