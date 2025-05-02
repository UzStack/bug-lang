package types

type RuntimeTypes string

const (
	String RuntimeTypes = "String"
	Number RuntimeTypes = "Number"
)

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

type RuntimeValue struct {
	Type  RuntimeTypes
	Value any
}
