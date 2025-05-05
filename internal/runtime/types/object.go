package types

import (
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
)

type ObjectValue struct {
	Name       string
	Enviroment *enviroment.Enviroment
}

func NewObject(name string, env *enviroment.Enviroment) Object {
	return &ObjectValue{
		Name:       name,
		Enviroment: env,
	}
}

func (o *ObjectValue) GetEnviroment() *enviroment.Enviroment {
	return o.Enviroment
}

func (a *ObjectValue) GetValue() any {
	return a.Name
}
