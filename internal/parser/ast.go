package parser

type NodeType string

const (
	VariableDeclarationNode NodeType = "VariableDeclaration"
	CallStatementNode       NodeType = "CallStatement"
	ProgramNode             NodeType = "Program"
	StringLiteral           NodeType = "StringLiteral"
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

type TokenStatement struct {
	Value any
	Kind  NodeType
}

func NewTokenStatement(kind NodeType, value any) *TokenStatement {
	return &TokenStatement{
		Kind:  kind,
		Value: value,
	}
}

// Call Statement
type CallStatement struct {
	*Statement
	Name  any
	Value any
	Args  []any
}

func NewCallStatement(line int, name any, args []any) *CallStatement {
	return &CallStatement{
		Statement: &Statement{
			Kind: CallStatementNode,
			Line: line,
		},
		Name: name,
		Args: args,
	}
}

// Variable Declatation
type VariableDeclaration struct {
	*Statement
	Name  any
	Value any
}

func NewVariableDeclaration(line int, name any, value any) *VariableDeclaration {
	return &VariableDeclaration{
		Statement: &Statement{
			Kind: VariableDeclarationNode,
			Line: line,
		},
		Name:  name,
		Value: value,
	}
}
