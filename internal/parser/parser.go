package parser

import (
	"fmt"
	"log"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/pkg/utils"
)

type parser struct {
	Tokens []*lexar.Token
	Index  int
}

func NewParser(tokens []*lexar.Token) *parser {
	return &parser{
		Tokens: tokens,
		Index:  0,
	}
}
func (p parser) At() *lexar.Token {
	return p.Tokens[p.Index]
}
func (p parser) Prev() *lexar.Token {
	return p.Tokens[p.Index-1]
}
func (p *parser) Next() *lexar.Token {
	p.Index++
	return p.Tokens[p.Index-1]
}

func (p parser) IsEOF() bool {
	return p.At().Type != lexar.EOF
}

func (p *parser) ParseAssignmentExpression() any {
	left := p.ParseLogicalExpression()
	if p.At().Type == lexar.Equals {
		p.Next()
		return &AssignmentExpression{
			Statement: &Statement{
				Line: p.At().Line,
				Kind: AssignmentExpressionNode,
			},
			Owner: left,
			Value: p.ParseAdditiveExpression(),
		}
	}
	return left
}

func (p *parser) ParseReturnStatement() any {
	p.Next()
	return &ReturnStatement{
		Statement: &Statement{
			Kind: ReturnNode,
			Line: p.At().Line,
		},
		Value: p.ParseAssignmentExpression(),
	}
}

func (p *parser) ParseCallExpression() any {
	caller := p.ParsePrimaryExpression()

	if p.At().Type == lexar.OpenParen {
		return &CallStatement{
			Statement: &Statement{
				Line: p.At().Line,
				Kind: CallStatementNode,
			},
			Caller: &Caller{
				Statement: &Statement{
					Line: p.At().Line,
					Kind: CallerNode,
				},
				Name: caller.(*IdentifierStatement).Value.(string),
			},
			Args: p.ParseArgs(),
		}
	}
	return caller
}

func (p *parser) ParseIfStatement() any {
	p.Next()
	p.Except(lexar.OpenParen, "Except open paren IF Statement")
	condition := p.ParseLogicalExpression()
	p.Except(lexar.CloseParen, "Except close paren IF Statement")
	body := p.ParseBody()
	var chields []any
	for p.At().Type == lexar.ElseIf {
		chields = append(chields, p.ParseElseIfStatement())
	}
	if p.At().Type == lexar.Else {
		chields = append(chields, p.ParseElseStatement())
	}
	return &IfStatement{
		Statement: &Statement{
			Kind: IfStatementNode,
			Line: p.At().Line,
		},
		Condition: condition,
		Body:      body,
		Childs:    chields,
	}
}

func (p *parser) ParseElseIfStatement() any {
	p.Next()
	p.Except(lexar.OpenParen, "Except open paren IF Statement")
	condition := p.ParseLogicalExpression()
	p.Except(lexar.CloseParen, "Except close paren IF Statement")
	body := p.ParseBody()
	return &ElseIfStatement{
		Statement: &Statement{
			Kind: ElseIfStatementNode,
			Line: p.At().Line,
		},
		Condition: condition,
		Body:      body,
	}
}

func (p *parser) ParseElseStatement() any {
	p.Next()
	body := p.ParseBody()
	return &ElseStatement{
		Statement: &Statement{
			Kind: ElseStatementNode,
			Line: p.At().Line,
		},
		Body: body,
	}
}

func (p *parser) ParseLogicalExpression() any {
	left := p.ParseRelationalExpression()
	for utils.InArray(p.At().Value, []any{"&&", "||"}) {
		operator := p.Next().Value
		right := p.ParseRelationalExpression()

		left = &BinaryExpression{
			Statement: &Statement{
				Kind: BinaryOperatorNode,
				Line: p.At().Line,
			},
			Right:    right,
			Left:     left,
			Operator: operator,
		}
	}
	return left
}

func (p *parser) ParseRelationalExpression() any {
	left := p.ParseAdditiveExpression()

	for utils.InArray(p.At().Value, []any{"==", ">=", "<=", "<", ">", "!="}) {
		operator := p.Next().Value
		right := p.ParseAdditiveExpression()

		left = &BinaryExpression{
			Statement: &Statement{
				Kind: BinaryOperatorNode,
				Line: p.At().Line,
			},
			Right:    right,
			Left:     left,
			Operator: operator,
		}
	}

	return left
}

func (p *parser) ParseAdditiveExpression() any {
	left := p.ParseMultiplicativeExpression()

	for p.At().Value == "+" || p.At().Value == "-" {
		operator := p.Next().Value
		right := p.ParseMultiplicativeExpression()
		left = &BinaryExpression{
			Statement: &Statement{
				Kind: BinaryOperatorNode,
				Line: p.At().Line,
			},
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left
}

func (p *parser) ParseMultiplicativeExpression() any {
	left := p.ParseCallExpression()

	for p.At().Value == "*" || p.At().Value == "/" || p.At().Value == "%" {
		operator := p.Next().Value
		right := p.ParseCallExpression()
		left = &BinaryExpression{
			Statement: &Statement{
				Kind: BinaryOperatorNode,
				Line: p.At().Line,
			},
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left
}

func (p *parser) ParseArgs() []any {
	p.Except(lexar.OpenParen, "Except open paren ARGS")
	defer p.Except(lexar.CloseParen, "Except close paren ARGS")
	var args []any
	if p.At().Type == lexar.CloseParen {
		return args
	}
	return p.ParserArgList()
}

func (p *parser) ParserArgList() []any {
	args := []any{p.ParseAssignmentExpression()}
	for p.At().Type == lexar.Comma {
		p.Next()
		args = append(args, p.ParseAssignmentExpression())
	}
	return args
}

func (p *parser) Except(tokenType lexar.TokenType, errMsg string) *lexar.Token {
	if p.Next().Type != tokenType {
		panic(fmt.Sprintf("Parser error: %s -> %s - Expected: %d Line: %d", errMsg, p.Tokens[p.Index-1].Value, tokenType, p.Index))
	}
	return p.Prev()
}

func (p *parser) ParseVariableDeclaration() any {
	p.Next()
	identifier := p.Except(lexar.Identifier, "O'zgaruvchi nomi nato'g'ri")
	p.Except(lexar.Equals, "O'zgaruvchi yaratishda xatolik yuz berdi")
	value, ok := identifier.Value.(string)
	if !ok {
		panic(fmt.Sprintf("Error: %d", p.Index))
	}
	declatation := &VariableDeclaration{
		Statement: &Statement{
			Kind: VariableDeclarationNode,
			Line: p.At().Line,
		},
		Name:  value,
		Value: p.ParseAdditiveExpression(),
	}
	if p.At().Type == lexar.Semicolon {
		p.Next()
	}
	return declatation
}

func (p *parser) ParseForStatement() any {
	p.Next()
	p.Except(lexar.OpenParen, "Except open paren FOR")
	condition := p.ParseLogicalExpression()
	p.Except(lexar.CloseParen, "Except close paren FOR")
	return &ForStatement{
		Statement: &Statement{
			Kind: ForNode,
			Line: p.At().Line,
		},
		Condition: condition,
		Body:      p.ParseBody(),
	}
}

func (p *parser) ParseBody() []any {
	var body []any
	p.Except(lexar.OpenBrace, "Except open brace")
	for p.At().Type != lexar.CloseBrace {
		body = append(body, p.ParseStatement())
	}
	p.Except(lexar.CloseBrace, "Except close brace")
	return body
}

func (p *parser) ParsePrimaryExpression() any {
	switch p.At().Type {
	case lexar.Number:
		return &NumberLiteral{
			Statement: &Statement{
				Kind: NumberLiteralNode,
				Line: p.At().Line,
			},
			Value: p.Next().Value,
		}
	case lexar.String:
		return &StringLiteral{
			Statement: &Statement{
				Kind: NumberLiteralNode,
				Line: p.At().Line,
			},
			Value: p.Next().Value,
		}
	case lexar.Identifier:
		return &IdentifierStatement{
			Statement: &Statement{
				Kind: IdentifierNode,
				Line: p.At().Line,
			},
			Value: p.Next().Value,
		}
	default:
		p.Next()
	}
	return 0
}

func (p *parser) ParseFnDeclaration() any {
	p.Next()
	identifier := p.Except(lexar.Identifier, "Funcsiya nomi nato'g'ri")
	args := p.ParseArgs()
	var params []*IdentifierStatement
	for _, arg := range args {
		param, ok := arg.(*IdentifierStatement)
		if !ok {
			log.Fatal("Funcsiya parametri nato'g'ri")
		}
		params = append(params, param)
	}
	return &FunctionDeclaration{
		Statement: &Statement{
			Kind: FunctionDeclarationNode,
			Line: p.At().Line,
		},
		Name:   identifier.Value.(string),
		Params: params,
		Body:   p.ParseBody(),
	}
}

func (p *parser) ParseStatement() any {
	switch p.At().Type {
	case lexar.Var:
		return p.ParseVariableDeclaration()
	case lexar.Fn:
		return p.ParseFnDeclaration()
	case lexar.If:
		return p.ParseIfStatement()
	case lexar.For:
		return p.ParseForStatement()
	case lexar.Return:
		return p.ParseReturnStatement()
	default:
		return p.ParseAssignmentExpression()
	}
}
func (p *parser) CreateAST() any {
	program := &Program{
		Statement: &Statement{
			Kind: ProgramNode,
			Line: -1,
		},
		Body: []any{},
	}
	for p.IsEOF() {
		stmt := p.ParseStatement()
		program.Body = append(program.Body, stmt)
	}
	return program
}
