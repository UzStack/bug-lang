package parser

import "github.com/UzStack/bug-lang/internal/runtime/enviroment"

var STDLIBS = []any{
	"math",
	"ffi",
	"json",
	"requestsgo",
}

type Statement struct {
	Line int
}

type Expression struct {
	Line int
}

type ProgramNode struct {
	Body []any
	Line int
}
type ModuleNode struct {
	Line int
	Name string
	Body []any
	Path string
}

type StdModuleNode struct {
	Line int
	Name string
	Path string
}

type BinaryNode struct {
	Line     int
	Left     any
	Right    any
	Operator any
}

type IdentifierNode struct {
	Line  int
	Value any
}

type BaseNode struct {
	Line int
}

// Call Statement
type CallNode struct {
	Line   int
	Caller any
	Value  any
	Args   []any
}

// Variable Declatation
type VariableDeclarationNode struct {
	Line  int
	Name  string
	Value any
}

type CallerNode struct {
	Line int
	Name string
}

type NumberLiteralNode struct {
	Line  int
	Value any
}
type FloatLiteralNode struct {
	Line  int
	Value any
}

type StringLiteralNode struct {
	Line  int
	Value any
}

type FunctionDeclarationNode struct {
	Line   int
	Name   string
	Params []any
	Body   []any
}

type IfNode struct {
	Line      int
	Condition any
	Body      []any
	Childs    []any
}

type ElseIfNode struct {
	Line      int
	Condition any
	Body      []any
}

type ElseNode struct {
	Line int
	Body []any
}

type ForNode struct {
	Line      int
	Condition any
	Body      []any
}

type AssignmentNode struct {
	Line  int
	Owner any
	Value any
}

type ReturnNode struct {
	Line  int
	Value any
}

type ArrayNode struct {
	Line   int
	Values []any
	Left   any
}

type MapNode struct {
	Line   int
	Values map[string]any
}

type MemberNode struct {
	Line     int
	Left     any
	Prop     any
	Computed bool
	Assign   any
}

type ClassDeclarationNode struct {
	Line       int
	Name       any
	Body       []any
	Methods    []*FunctionDeclarationNode
	Extends    []any
	Enviroment *enviroment.Enviroment
}

type ObjectNode struct {
	Line   int
	Caller any
}
