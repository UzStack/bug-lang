package std

import (
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func Load(env *enviroment.Enviroment) {
	env.DeclareVariable("print", types.NewNativeFunction(Print), -1)
	env.DeclareVariable("println", types.NewNativeFunction(Println), -1)
	env.DeclareVariable("console", types.NewNativeFunction(Print), -1)
	env.DeclareVariable("consoleln", types.NewNativeFunction(Println), -1)

	env.DeclareVariable("input", types.NewNativeFunction(Input), -1)
	env.DeclareVariable("header", types.NewNativeFunction(Header), -1)
	env.DeclareVariable("true", types.NewBool(true), -1)
	env.DeclareVariable("false", types.NewBool(false), -1)
	env.DeclareVariable("null", types.NewNull(), -1)
}
