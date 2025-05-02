package enviroment

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

type Enviroment struct {
	Variables map[string]any
	Owner     *Enviroment
}

func NewEnv(owner *Enviroment) *Enviroment {
	return &Enviroment{
		Variables: map[string]any{},
		Owner:     owner,
	}
}

func NewGlobalEnv() *Enviroment {
	env := NewEnv(nil)
	env.DeclareVariable("print", &types.NativeFunctionDeclaration{
		Type: "native-function",
		Call: std.Print,
	}, -1)
	return env
}

func (e *Enviroment) DeclareVariable(name string, value any, line int) {
	_, ok := e.Variables[name]
	if ok {
		panic(fmt.Sprintf("O'zvaruvchi mavjud %s Line: %d", name, line))
	}
	e.Variables[name] = value
}

func (e *Enviroment) GetVariable(name string, line int) any {
	res, ok := e.Variables[name]
	if !ok {
		panic(fmt.Sprintf("O'zgaruvchi topilmadi %s Line: %d", name, line))
	}
	return res
}
