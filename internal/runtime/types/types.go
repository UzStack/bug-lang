package types

import "github.com/UzStack/bug-lang/internal/runtime/enviroment"

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
	Enviroment  *enviroment.Enviroment
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
