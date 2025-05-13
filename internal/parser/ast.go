package parser

import "github.com/UzStack/bug-lang/internal/runtime/enviroment"

type NodeType string

const (
	VariableDeclarationNode  NodeType = "VariableDeclaration"
	FunctionDeclarationNode  NodeType = "FunctionDeclaration"
	CallStatementNode        NodeType = "CallStatement"
	CallerNode               NodeType = "Caller"
	ProgramNode              NodeType = "Program"
	StringLiteralNode        NodeType = "StringLiteral"
	IdentifierNode           NodeType = "Identifier"
	NumberLiteralNode        NodeType = "NumberLiteral"
	BinaryOperatorNode       NodeType = "BinaryOperator"
	IfStatementNode          NodeType = "IfStatement"
	ElseIfStatementNode      NodeType = "ElseIfStatement"
	ElseStatementNode        NodeType = "ElseStatement"
	ForNode                  NodeType = "For"
	AssignmentExpressionNode NodeType = "AssignmentExpression"
	ReturnNode               NodeType = "Return"
	ArrayNode                NodeType = "Array"
	MemberNode               NodeType = "Member"
	MapNode                  NodeType = "Map"
	ObjectNode               NodeType = "ObjectNode"
)

var STDLIBS = []any{
	"math",
	"ffi",
	"json",
}

type Statement struct {
	Kind NodeType
	Line int
}

type Expression struct {
	Line int
}

type Program struct {
	Body []any
	Line int
}
type Module struct {
	Line int
	Name string
	Body []any
	Path string
}

type StdModule struct {
	Line int
	Name string
	Path string
}

type BinaryExpression struct {
	Line     int
	Left     any
	Right    any
	Operator any
}

type IdentifierStatement struct {
	Line  int
	Value any
}

type BaseStatement struct {
	Line int
}

// Call Statement
type CallStatement struct {
	Line   int
	Caller any
	Value  any
	Args   []any
}

// Variable Declatation
type VariableDeclaration struct {
	Line  int
	Name  string
	Value any
}

type Caller struct {
	Line int
	Name string
}

type NumberLiteral struct {
	Line  int
	Value any
}
type FloatLiteral struct {
	Line  int
	Value any
}

type StringLiteral struct {
	Line  int
	Value any
}

type FunctionDeclaration struct {
	Line   int
	Name   string
	Params []any
	Body   []any
}

type IfStatement struct {
	Line      int
	Condition any
	Body      []any
	Childs    []any
}

type ElseIfStatement struct {
	Line      int
	Condition any
	Body      []any
}

type ElseStatement struct {
	Line int
	Body []any
}

type ForStatement struct {
	Line      int
	Condition any
	Body      []any
}

type AssignmentExpression struct {
	Line  int
	Owner any
	Value any
}

type ReturnStatement struct {
	Line  int
	Value any
}

type ArrayExpression struct {
	Line   int
	Values []any
	Left   any
}

type MapExpression struct {
	Line   int
	Values map[string]any
}

type MemberExpression struct {
	Line     int
	Left     any
	Prop     any
	Computed bool
	Assign   any
}

type ClassDeclaration struct {
	Line       int
	Name       any
	Body       []any
	Methods    []*FunctionDeclaration
	Extends    []any
	Enviroment *enviroment.Enviroment
}

type ObjectExpression struct {
	Line   int
	Caller any
}
