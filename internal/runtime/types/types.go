package types

type FunctionDeclaration struct {
	Name   any
	Params []any
	Body   map[string]any
	Type   string
}

func NewFunctionDeclaration(line int, name string, body map[string]any) *FunctionDeclaration {
	return &FunctionDeclaration{
		Type:   "function",
		Name:   name,
		Body:   body,
		Params: []any{},
	}
}

type NativeFunctionDeclaration struct {
	Type string
	Call any
}

func NewNativeFunctionDeclaration(call any) *NativeFunctionDeclaration {
	return &NativeFunctionDeclaration{
		Type: "native-function",
		Call: call,
	}
}
