package enviroment

type Module struct {
	Modules map[string]any
	Owner   *Module
}

func NewModule() *Module {
	return &Module{
		Modules: map[string]any{},
		Owner:   nil,
	}
}

func (m *Module) Add(path string, module any) {
	m.Modules[path] = module
}
func (m *Module) Get(path string) any {
	return m.Modules[path]
}
