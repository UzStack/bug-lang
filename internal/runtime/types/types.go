package types

import "github.com/UzStack/bug-lang/internal/parser"

type RuntimeTypes string

const (
	String   RuntimeTypes = "String"
	Number   RuntimeTypes = "Number"
	Function RuntimeTypes = "Function"
	Flow     RuntimeTypes = "Flow"
)

type FunctionDeclaration struct {
	Name   any
	Params []*parser.IdentifierStatement
	Body   []any
	Type   RuntimeTypes
}

type NativeFunctionDeclaration struct {
	Type string
	Call any
}

type RuntimeValue struct {
	Type  RuntimeTypes
	Value any
}

type NullValue struct {
	Type  RuntimeTypes
	Value any
}

type FlowValue struct {
	Type    RuntimeTypes
	Catched bool
}
