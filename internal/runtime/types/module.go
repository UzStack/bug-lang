package types

import "github.com/UzStack/bug-lang/internal/runtime/enviroment"

type ModuleValue struct {
	Enviroment *enviroment.Enviroment
}

func NewModule(env *enviroment.Enviroment) Object {
	return &ModuleValue{
		Enviroment: env,
	}
}

func (o *ModuleValue) GetEnviroment() *enviroment.Enviroment {
	return o.Enviroment
}

func (a *ModuleValue) GetValue() any {
	return a.Enviroment
}
