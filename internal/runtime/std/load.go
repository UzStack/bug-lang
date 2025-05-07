package std

import (
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func Load(env *enviroment.Enviroment) {
	env.DeclareVariable("print", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Print,
	}, -1)
	env.DeclareVariable("input", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Input,
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
}
