package parser

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/lexar"
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
	left := p.ParseAdditiveExpression()
	return left
}

func (p *parser) ParseCallExpression() any {
	caller := p.ParsePrimaryExpression()
	if p.At().Type == lexar.OpenParen {
		p.Next()
		return NewCallStatement(
			p.At().Line,
			NewCaller(caller.(*IdentifierStatement).Value.(string)),
			p.ParseArgs())
	}
	return caller
}

func (p *parser) ParseAdditiveExpression() any {
	left := p.ParseMultiplicativeExpression()

	for p.At().Value == "+" || p.At().Value == "-" {
		operator := p.Next().Value
		right := p.ParseMultiplicativeExpression()

		left = NewBinaryExpression(
			left,
			right,
			operator,
		)
	}
	return left
}

func (p *parser) ParseMultiplicativeExpression() any {
	left := p.ParseCallExpression()
	for p.At().Value == "*" || p.At().Value == "/" || p.At().Value == "%" {
		operator := p.Next().Value
		right := p.ParseCallExpression()

		left = NewBinaryExpression(
			left,
			right,
			operator,
		)
	}
	return left
}

func (p *parser) ParseArgs() []any {
	var args []any
	if p.At().Type == lexar.CloseParen {
		return args
	}
	return p.ParserArgList()
}

func (p *parser) ParserArgList() []any {
	args := []any{p.ParseAssignmentExpression()}
	for p.Next().Type == lexar.Comma {
		args = append(args, p.ParseAssignmentExpression())
	}
	if p.At().Type == lexar.CloseParen {
		p.Next()
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
	declatation := NewVariableDeclaration(
		p.At().Line, value,
		p.ParseAssignmentExpression(),
	)
	if p.At().Type == lexar.Semicolon {
		p.Next()
	}
	return declatation
}

func (p *parser) ParsePrimaryExpression() any {
	switch p.At().Type {
	case lexar.Number:
		return &NumberLiteralNode{
			Kind:  NumberLiteral,
			Value: p.Next().Value,
		}
	case lexar.String:
		return &StringLiteralNode{
			Kind:  StringLiteral,
			Value: p.Next().Value,
		}
	case lexar.Identifier:
		return NewIdentifier(p.Next().Value)
	default:
		p.Next()
	}
	return 0
}

func (p *parser) ParseStatement() any {
	switch p.At().Type {
	case lexar.Var:
		return p.ParseVariableDeclaration()
	default:
		return p.ParseAssignmentExpression()
	}
}
func (p *parser) CreateAST() any {
	program := NewProgram(1)
	for p.IsEOF() {
		stmt := p.ParseStatement()
		program.Body = append(program.Body, stmt)
	}
	return program
}
