package types

type RuntimeTypes string

const (
	String   RuntimeTypes = "String"
	Number   RuntimeTypes = "Number"
	Function RuntimeTypes = "Function"
	Flow     RuntimeTypes = "Flow"
)

type FunctionDeclaration struct {
	Name        any
	Params      []any
	Body        []any
	Type        RuntimeTypes
	OwnerObject any
}

type NativeFunctionDeclaration struct {
	Type string
	Call any
}

type RuntimeValue struct {
	Type  RuntimeTypes
	Value any
}

type FlowValue struct {
	Type    RuntimeTypes
	Catched bool
}

type ReturnValue struct {
	Value any
}

type Object interface {
	GetValue() any
}
