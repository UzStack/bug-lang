package types

import "github.com/UzStack/bug-lang/internal/parser"

type RuntimeTypes string

const (
	String   RuntimeTypes = "String"
	Number   RuntimeTypes = "Number"
	Function RuntimeTypes = "Function"
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
