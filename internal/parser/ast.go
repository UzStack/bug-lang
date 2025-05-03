package parser

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
)

type Statement struct {
	Kind NodeType
	Line int
}

type Expression struct {
	*Statement
}

type Program struct {
	*Statement
	Body []any
}

type BinaryExpression struct {
	*Statement
	Left     any
	Right    any
	Operator any
}

type IdentifierStatement struct {
	*Statement
	Value any
}

type BaseStatement struct {
	*Statement
}

// Call Statement
type CallStatement struct {
	*Statement
	Caller *Caller
	Value  any
	Args   []any
}

// Variable Declatation
type VariableDeclaration struct {
	*Statement
	Name  string
	Value any
}

type Caller struct {
	*Statement
	Name string
}

type NumberLiteral struct {
	*Statement
	Value any
}

type StringLiteral struct {
	*Statement
	Value any
}

type FunctionDeclaration struct {
	*Statement
	Name   string
	Params []*IdentifierStatement
	Body   []any
}

type IfStatement struct {
	*Statement
	Condition any
	Body      []any
	Childs    []any
}

type ElseIfStatement struct {
	*Statement
	Condition any
	Body      []any
}

type ElseStatement struct {
	*Statement
	Body []any
}

type ForStatement struct {
	*Statement
	Condition any
	Body      []any
}

type AssignmentExpression struct {
	*Statement
	Owner any
	Value any
}

type ReturnStatement struct {
	*Statement
	Value any
}
