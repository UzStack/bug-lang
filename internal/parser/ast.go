package parser

type NodeType string

const (
	VariableDeclarationNode NodeType = "VariableDeclaration"
	CallStatementNode       NodeType = "CallStatement"
	ProgramNode             NodeType = "Program"
	StringLiteral           NodeType = "StringLiteral"
	Identifier              NodeType = "Identifier"
	NumberLiteral           NodeType = "NumberLiteral"
	BinaryOperator          NodeType = "BinaryOperator"
)

type Statement struct {
	Kind NodeType
	Line int
}

type Expression struct {
	Statement
}

type Program struct {
	*Statement
	Body []any
}

func NewProgram(line int) *Program {
	return &Program{
		Statement: &Statement{
			Kind: ProgramNode,
			Line: line,
		},
	}
}

type BinaryExpression struct {
	Kind     NodeType
	Left     any
	Right    any
	Operator any
}

func NewBinaryExpression(left, right any, operator any) *BinaryExpression {
	return &BinaryExpression{
		Kind:     BinaryOperator,
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

type IdentifierStatement struct {
	Kind  NodeType
	Value any
}

func NewIdentifier(value any) *IdentifierStatement {
	return &IdentifierStatement{
		Kind:  Identifier,
		Value: value,
	}
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

func NewCallStatement(line int, caller *Caller, args []any) *CallStatement {
	return &CallStatement{
		Statement: &Statement{
			Kind: CallStatementNode,
			Line: line,
		},
		Caller: caller,
		Args:   args,
	}
}

// Variable Declatation
type VariableDeclaration struct {
	*Statement
	Name  string
	Value any
}

func NewVariableDeclaration(line int, name string, value any) *VariableDeclaration {
	return &VariableDeclaration{
		Statement: &Statement{
			Kind: VariableDeclarationNode,
			Line: line,
		},
		Name:  name,
		Value: value,
	}
}

type Caller struct {
	Kind any
	Name string
}

func NewCaller(name string) *Caller {
	return &Caller{
		Kind: Identifier,
		Name: name,
	}
}

type NumberLiteralNode struct {
	Kind  NodeType
	Value any
}

type StringLiteralNode struct {
	Kind  NodeType
	Value any
}
