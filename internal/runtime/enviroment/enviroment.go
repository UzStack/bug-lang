package enviroment

import (
	"fmt"
)

type Enviroment struct {
	Variables map[string]any
	Owner     *Enviroment
	Modules   *Module
}

func NewEnv(owner *Enviroment) *Enviroment {
	return &Enviroment{
		Variables: map[string]any{},
		Owner:     owner,
	}
}

func NewGlobalEnv() *Enviroment {
	env := NewEnv(nil)
	env.Modules = NewModule()
	return env
}

// ponytail: qayta e'lon qilishga ruxsat (qattiq panic o'rniga).
// Aks holda for tanasidagi `var i = ...` ikkinchi iteratsiyada interpreterni o'ldirardi.
func (e *Enviroment) DeclareVariable(name string, value any, line int) any {
	e.Variables[name] = value
	return value
}

func (e *Enviroment) GetModules() *Module {
	if e.Modules != nil {
		return e.Modules
	}
	if e.Owner != nil {
		return e.Owner.GetModules()
	}
	return NewModule()
}

func (e *Enviroment) AssignmenVariable(name string, value any, line int) any {
	e.Variables[name] = value
	return value
}

func (e *Enviroment) GetVariable(name string, line int) (any, error) {
	res, ok := e.Variables[name]
	if !ok {
		if e.Owner != nil {
			return e.Owner.GetVariable(name, line)
		}
		if !ok {
			return nil, fmt.Errorf("o'zgaruvchi topilmadi %s Line: %d", name, line)
		}
	}
	return res, nil
}
