package types

type FunctionDeclaration struct {
	Name   any
	Params []any
	Body   map[string]any
	Type   string
}

type NativeFunctionDeclaration struct {
	Type string
	Call any
}
