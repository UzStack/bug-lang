package enviroment

import (
	"fmt"
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
	return env
}

func (e *Enviroment) DeclareVariable(name string, value any, line int) any {
	_, ok := e.Variables[name]
	if ok {
		panic(fmt.Sprintf("O'zvaruvchi mavjud %s Line: %d", name, line))
	}
	e.Variables[name] = value
	return value
}

func (e *Enviroment) AssignmenVariable(name string, value any, line int) any {
	e.Variables[name] = value
	return value
}

func (e *Enviroment) GetVariable(name string, line int) any {
	res, ok := e.Variables[name]
	if !ok {
		if e.Owner != nil {
			return e.Owner.GetVariable(name, line)
		}
		if !ok {
			panic(fmt.Sprintf("O'zgaruvchi topilmadi %s Line: %d", name, line))
		}
	}
	return res
}
