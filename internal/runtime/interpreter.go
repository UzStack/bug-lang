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
	case *parser.NumberLiteralNode:
		value, err := strconv.Atoi(node.Value.(string))
		if err != nil {
			fmt.Println("Type error: ", node.Value, "not integer")
			os.Exit(1)
		}
		return types.NewInt(value), nil
	case *parser.StringLiteralNode:
		value, err := node.Value.(string)
		if !err {
			fmt.Println("Type not string", value)
			os.Exit(1)
		}
		return types.NewString(value), nil
	case *parser.FloatLiteralNode:
		value, err := strconv.ParseFloat(node.Value.(string), 64)
		if err != nil {
			fmt.Println("Type error: ", node.Value, "not float")
			os.Exit(1)
		}
		return types.NewFloat(value), nil
	case *parser.ModuleNode:
		return EvalModuleStatement(node, env)
	case *parser.StdModuleNode:
		return EvalStdModuleStatement(node, env)
	case *parser.ProgramNode:
		return EvalProgram(node, env)
	case *parser.IdentifierNode:
		return EvalIdentifier(node, env)
	case *parser.CallNode:
		return EvalCallStatement(node, env)
	case *parser.VariableDeclarationNode:
		return VariableDeclaration(node, env)
	case *parser.BinaryNode:
		return EvalBinaryExpression(node, env)
	case *parser.FunctionDeclarationNode:
		return EvalFunctionDeclaration(node, env, nil)
	case *parser.IfNode:
		return EvalIfStatement(node, env)
	case *parser.ElseIfNode:
		return EvalElseIfStatement(node, env)
	case *parser.ElseNode:
		return EvalElseStatement(node, env)
	case *parser.ForNode:
		return EvalForStatement(node, env)
	case *parser.AssignmentNode:
		return EvalAssignmentExpression(node, env)
	case *parser.ReturnNode:
		return EvalReturnStatement(node, env)
	case *parser.MemberNode:
		return EvalMemberExpression(node, env)
	case *parser.ArrayNode:
		return EvalArrayExpression(node, env)
	case *parser.ClassDeclarationNode:
		return EvalClassDeclaration(node, env)
	case *parser.MapNode:
		return EvalMapDeclaration(node, env)
	case *parser.ObjectNode:
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

func EvalStdModuleStatement(node *parser.StdModuleNode, env *enviroment.Enviroment) (any, error) {
	env.DeclareVariable(node.Name, types.NewStdLib(node.Name, std.STDLIBS[node.Path]), node.Line)
	return nil, nil
}
func EvalModuleStatement(node *parser.ModuleNode, env *enviroment.Enviroment) (any, error) {
	if module := env.GetModules().Get(node.Path); module != nil {
		return env.DeclareVariable(node.Name, module, node.Line), nil
	}
	// WARNING: buni to'g'irlash kerak module ichida global env ishlatilyapti bu xato
	scope := enviroment.NewEnv(env)
	// std.Load(scope)
	for _, stmt := range node.Body {
		if _, err := Interpreter(stmt, scope); err != nil {
			return nil, err
		}
	}
	env.DeclareVariable(node.Name, types.NewModule(scope), node.Line)
	env.GetModules().Add(node.Path, types.NewModule(scope))
	return nil, nil
}

func DeclareExtends(class *parser.ClassDeclarationNode, env *enviroment.Enviroment, scope *enviroment.Enviroment, obj types.Object) (map[string]*enviroment.Enviroment, error) {
	envs := make(map[string]*enviroment.Enviroment)
	for _, extend := range class.Extends {
		extScope := enviroment.NewEnv(scope)
		res, err := Interpreter(extend, scope)
		if err != nil {
			return nil, err
		}
		exd := res.(*parser.ClassDeclarationNode)
		DeclareExtends(exd, env, scope, obj)
		for _, method := range exd.Methods {
			EvalFunctionDeclaration(method, extScope, obj)
		}
		for key, value := range extScope.Variables {
			scope.Variables[key] = value
		}
		envs[exd.Name.(*parser.IdentifierNode).Value.(string)] = extScope
	}
	return envs, nil
}

func EvalObjectExpression(node *parser.ObjectNode, env *enviroment.Enviroment) (any, error) {
	caller := node.Caller.(*parser.CallNode)
	var methods []*parser.FunctionDeclarationNode
	var className string
	var class *parser.ClassDeclarationNode
	switch t := caller.Caller.(type) {
	case *parser.IdentifierNode:
		className = t.Value.(string)
		res, err := env.GetVariable(className, node.Line)
		if err != nil {
			return nil, err
		}
		class = res.(*parser.ClassDeclarationNode)
		methods = class.Methods
	case *parser.MemberNode:
		className = t.Prop.(*parser.IdentifierNode).Value.(string)
		res, err := Interpreter(t, env)
		if err != nil {
			return nil, err
		}
		class = res.(*parser.ClassDeclarationNode)
		methods = class.Methods
	}
	scope := enviroment.NewEnv(class.Enviroment)
	obj := types.NewObject(className, scope).(*types.ObjectValue)
	EvalFunctionDeclaration(&parser.FunctionDeclarationNode{
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
	_, err = EvalCallStatement(&parser.CallNode{
		Caller: &parser.IdentifierNode{
			Value: "init",
		},
		Args: caller.Args,
	}, scope)

	if err != nil {
		return nil, err
	}
	return obj, nil

}

func EvalClassDeclaration(node *parser.ClassDeclarationNode, env *enviroment.Enviroment) (any, error) {
	node.Enviroment = env
	return env.DeclareVariable(node.Name.(*parser.IdentifierNode).Value.(string), node, node.Line), nil
}

func VariableDeclaration(node *parser.VariableDeclarationNode, env *enviroment.Enviroment) (any, error) {
	value, err := Interpreter(node.Value, env)
	if err != nil {
		return nil, err
	}
	return env.DeclareVariable(node.Name, value, node.Line), nil
}

func EvalMapDeclaration(node *parser.MapNode, env *enviroment.Enviroment) (any, error) {
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
func EvalArrayExpression(node *parser.ArrayNode, env *enviroment.Enviroment) (any, error) {
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

func EvalReturnStatement(node *parser.ReturnNode, env *enviroment.Enviroment) (any, error) {
	value, err := Interpreter(node.Value, env)
	if err != nil {
		return nil, err
	}
	return &types.ReturnValue{
		Value: value,
	}, nil
}
func EvalMemberExpression(node *parser.MemberNode, env *enviroment.Enviroment) (any, error) {
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
				return nil, fmt.Errorf("map item not found: %s", index)
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
		name := node.Prop.(*parser.IdentifierNode).Value.(string)
		return left.Enviroment.AssignmenVariable(name, node.Assign, node.Line), nil
	} else {
		left, err := Interpreter(node.Left, env)
		if err != nil {
			return nil, err
		}
		prop := node.Prop.(*parser.IdentifierNode).Value.(string)
		switch t := left.(type) {
		case *types.ObjectValue:
			return t.Enviroment.GetVariable(prop, node.Line)
		case *types.StdLibValue:
			return reflect.ValueOf(t.Lib[prop]), nil
		case *types.ModuleValue:
			return t.Enviroment.GetVariable(prop, node.Line)
		case types.Object:
			v := reflect.ValueOf(left)
			method := v.MethodByName(string(strings.ToUpper(prop[:1]) + prop[1:]))
			if !method.IsValid() {
				return nil, fmt.Errorf("method %s not found in %T line: %d", prop, t, node.Line)
			}
			return method, nil
		case *enviroment.Enviroment:
			return t.GetVariable(prop, node.Line)
		default:
			return t.(map[string]any)[prop], nil
		}

	}
}

func EvalAssignmentExpression(node *parser.AssignmentNode, env *enviroment.Enviroment) (any, error) {
	switch t := node.Owner.(type) {
	case *parser.MemberNode:
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
		return env.AssignmenVariable(node.Owner.(*parser.IdentifierNode).Value.(string), value, node.Line), nil
	}
}

func EvalForStatement(node *parser.ForNode, env *enviroment.Enviroment) (any, error) {
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

func EvalIfStatement(node *parser.IfNode, env *enviroment.Enviroment) (any, error) {
	condition, err := Interpreter(node.Condition, env)
	if err != nil {
		return nil, err
	}
	// left, _ := Interpreter(node.Condition.(*parser.BinaryNode).Left, env)

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
		switch v := child.(type) {
		case *types.FlowValue:
			if v.Catched {
				return v, nil
			}
		case *types.ReturnValue:
			return v, nil
		default:
			pp.Print(v)
		}
	}
	return &types.FlowValue{
		Catched: false,
		Type:    types.Flow,
	}, nil
}

func EvalElseIfStatement(node *parser.ElseIfNode, env *enviroment.Enviroment) (any, error) {
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
func EvalElseStatement(node *parser.ElseNode, env *enviroment.Enviroment) (any, error) {
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

func EvalIdentifier(node *parser.IdentifierNode, env *enviroment.Enviroment) (any, error) {
	name, _ := node.Value.(string)
	return env.GetVariable(name, node.Line)
}

func EvalFunctionDeclaration(node *parser.FunctionDeclarationNode, env *enviroment.Enviroment, ownerObject any) (any, error) {
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

func EvalProgram(node *parser.ProgramNode, env *enviroment.Enviroment) (any, error) {
	for _, statement := range node.Body {
		if _, err := Interpreter(statement, env); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func EvalBinaryExpression(node *parser.BinaryNode, env *enviroment.Enviroment) (any, error) {
	leftRes, err := Interpreter(node.Left, env)
	if err != nil {
		return nil, err
	}
	rightRes, err := Interpreter(node.Right, env)
	if err != nil {
		return nil, err
	}

	leftVal := leftRes.(types.Object).GetValue()
	rightVal := rightRes.(types.Object).GetValue()

	leftKind := reflect.TypeOf(leftVal).Kind()
	rightKind := reflect.TypeOf(rightVal).Kind()

	switch node.Operator {
	case "&&", "||":
		leftBool, ok1 := leftVal.(bool)
		rightBool, ok2 := rightVal.(bool)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid types for logical operation: %T, %T", leftVal, rightVal)
		}
		var result bool
		if node.Operator == "&&" {
			result = leftBool && rightBool
		} else {
			result = leftBool || rightBool
		}
		return types.NewBool(result), nil

	case "+", "-", "*", "/", "%":
		// String qo‘shish
		if node.Operator == "+" && (leftKind == reflect.String || rightKind == reflect.String) {
			return types.NewString(fmt.Sprintf("%v%v", leftVal, rightVal)), nil
		}

		// Arifmetik amallar
		leftFloat, err := utils.Int2Float(leftVal)
		if err != nil {
			return nil, err
		}
		rightFloat, err := utils.Int2Float(rightVal)
		if err != nil {
			return nil, err
		}

		var result any
		switch node.Operator {
		case "+":
			result = leftFloat + rightFloat
		case "-":
			result = leftFloat - rightFloat
		case "*":
			result = leftFloat * rightFloat
		case "/":
			if rightFloat == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result = leftFloat / rightFloat
		case "%":
			result = int(leftFloat) % int(rightFloat)
		}

		// Float yoki int aniqlash
		if leftKind == reflect.Float64 || rightKind == reflect.Float64 || node.Operator == "/" {
			return types.NewFloat(result.(float64)), nil
		}
		intVal, _ := utils.Float2Int(result)
		return types.NewInt(intVal), nil

	case "==", "!=", ">", "<", ">=", "<=":
		// Turli tiplardagi solishtirishlar
		switch left := leftVal.(type) {
		case string:
			right, ok := rightVal.(string)
			if !ok {
				return nil, fmt.Errorf("type mismatch in comparison: %T and %T", leftVal, rightVal)
			}
			switch node.Operator {
			case "==":
				return types.NewBool(left == right), nil
			case "!=":
				return types.NewBool(left != right), nil
			case ">":
				return types.NewBool(left > right), nil
			case "<":
				return types.NewBool(left < right), nil
			case ">=":
				return types.NewBool(left >= right), nil
			case "<=":
				return types.NewBool(left <= right), nil
			}
		case float64, int:
			lv, _ := utils.Int2Float(leftVal)
			rv, _ := utils.Int2Float(rightVal)
			switch node.Operator {
			case "==":
				return types.NewBool(lv == rv), nil
			case "!=":
				return types.NewBool(lv != rv), nil
			case ">":
				return types.NewBool(lv > rv), nil
			case "<":
				return types.NewBool(lv < rv), nil
			case ">=":
				return types.NewBool(lv >= rv), nil
			case "<=":
				return types.NewBool(lv <= rv), nil
			}
		default:
			return nil, fmt.Errorf("unsupported types for comparison: %T", leftVal)
		}
	}

	return nil, fmt.Errorf("unsupported operator: %s", node.Operator)
}

func IsReturn(result any) (bool, any) {
	switch val := result.(type) {
	case *types.ReturnValue:
		return true, val
	}
	return false, nil
}

func EvalCallStatement(node *parser.CallNode, env *enviroment.Enviroment) (any, error) {
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
		// if len(v.Params) != len(node.Args) {
		// 	return nil, fmt.Errorf("%s funcsiyasi parametri xato berildi funcsiya kutmoqda: %d berildi: %d line: %d", v.Name, len(v.Params), len(node.Args), node.Line)
		// }
		var result any
		scope = enviroment.NewEnv(v.Enviroment)
		scope.AssignmenVariable("this", v.OwnerObject, node.Line)
		scope.AssignmenVariable("super", &types.NativeFunctionValue{
			Call: func(value *parser.ClassDeclarationNode) any {
				return v.OwnerObject.(*types.ObjectValue).Extends[value.Name.(*parser.IdentifierNode).Value.(string)]
			},
		}, node.Line)
		for index, arg := range v.Params {
			var name string
			switch v := arg.(type) {
			case *parser.AssignmentNode:
				if _, err := Interpreter(arg, scope); err != nil {
					return nil, err
				}
				name = v.Owner.(*parser.IdentifierNode).Value.(string)
			case *parser.IdentifierNode:
				name = v.Value.(string)
			}
			if len(node.Args) > index {
				param, err := Interpreter(node.Args[index], env)
				if err != nil {
					return nil, err
				}
				scope.AssignmenVariable(name, param, node.Line)
			}

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
			return utils.DecodeBug(res[0].Interface()), nil
		}
		return nil, nil
	default:
		pp.Print(v)
	}
	return nil, nil
}
