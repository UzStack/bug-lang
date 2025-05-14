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
	env.DeclareVariable("println", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Println,
	}, -1)
	env.DeclareVariable("console", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Print,
	}, -1)
	env.DeclareVariable("consoleln", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Println,
	}, -1)

	env.DeclareVariable("input", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Input,
	}, -1)
	env.DeclareVariable("header", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: Header,
	}, -1)
	env.DeclareVariable("true", types.NewBool(true), -1)
	env.DeclareVariable("false", types.NewBool(false), -1)
	env.DeclareVariable("null", types.NewNull(), -1)
}
