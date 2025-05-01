package parser

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/k0kubun/pp/v3"
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
	left := p.ParserCallExpression()
	return left
}

func (p *parser) ParserCallExpression() any {
	left := p.ParsePrimaryExpression()
	if p.Next().Type == lexar.OpenParen {
		identifier := p.Tokens[p.Index-2].Value
		return NewCallStatement(
			p.At().Line,
			identifier,
			p.ParseArgs())
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
		p.Next()
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
	declatation := NewVariableDeclaration(
		p.At().Line, identifier.Value,
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
		return map[string]any{
			"kind":  NumberLiteral,
			"value": p.Next().Value,
		}
	case lexar.String:
		return map[string]any{
			"kind":  StringLiteral,
			"value": p.Next().Value,
		}
	case lexar.Identifier:
		return map[string]any{
			"kind":  Identifier,
			"value": p.Next().Value,
		}
	default:
		return 0
	}
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
	pp.Println(program)
	return 0
}
