package parser

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/UzStack/bug-lang/assets"
	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/pkg/utils"
)

type parser struct {
	Tokens  []*lexar.Token
	Index   int
	BaseDir string
}

func NewParser(tokens []*lexar.Token, baseDir string) *parser {
	return &parser{
		Tokens:  tokens,
		Index:   0,
		BaseDir: baseDir,
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

// - Orders Of Prescidence -
// Assignment
// Object
// Logical
// Relational
// AdditiveExpr
// MultiplicitaveExpr
// CallMember
// Member
// PrimaryExpr

func (p *parser) ParseAssignmentExpression() any {
	left := p.ParseNewObjectExpression()
	if p.At().Type == lexar.Equals {
		p.Next()
		return &AssignmentNode{
			Line:  p.At().Line,
			Owner: left,
			Value: p.ParseAssignmentExpression(),
		}
	}
	return left
}

func (p *parser) ParseNewObjectExpression() any {
	if p.At().Type != lexar.New {
		return p.ParseArrayExpression()
	}
	p.Next()
	return &ObjectNode{
		Caller: p.ParseAssignmentExpression(),
	}
}

func (p *parser) ParseReturnStatement() any {
	p.Next()
	return &ReturnNode{
		Line:  p.At().Line,
		Value: p.ParseAssignmentExpression(),
	}
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
	return &IfNode{
		Line:      p.At().Line,
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
	return &ElseIfNode{
		Line:      p.At().Line,
		Condition: condition,
		Body:      body,
	}
}

func (p *parser) ParseElseStatement() any {
	p.Next()
	body := p.ParseBody()
	return &ElseNode{
		Line: p.At().Line,
		Body: body,
	}
}

func (p *parser) ParseMapExpression() any {
	if p.At().Type != lexar.OpenBrace {
		return p.ParseLogicalExpression()
	}
	return &MapNode{
		Line:   p.At().Line,
		Values: p.ParseMapItems(),
	}
}

func (p *parser) ParseMapItems() map[string]any {
	items := make(map[string]any)
	p.Except(lexar.OpenBrace, "Except open brace Object")
	if p.At().Type == lexar.CloseBrace {
		p.Next()
		return items
	}
	for p.At().Type != lexar.CloseBrace {
		key := p.Next().Value.(string)
		if p.At().Type != lexar.Colon {
			panic("Syntax error object: " + p.At().Value.(string))
		}
		p.Except(lexar.Colon, "Except colon Object")
		items[key] = p.ParseAssignmentExpression()
		// p.Except(lexar.Comma, "Except comma Object")
		if p.At().Type == lexar.Comma {
			p.Next()
		}
	}
	p.Except(lexar.CloseBrace, "Except close brace Object")
	return items
}

func (p *parser) ParseArrayExpression() any {
	if p.At().Type != lexar.OpenBracket {
		return p.ParseMapExpression()
	}
	return &ArrayNode{
		Line:   p.At().Line,
		Values: p.ParseArrayItems(),
	}
}

func (p *parser) ParseArrayItems() []any {
	var params []any
	p.Except(lexar.OpenBracket, "Except open bracket Array")
	if p.At().Type == lexar.CloseBracket {
		p.Next()
		return params
	}

	params = append(params, p.ParseAssignmentExpression())
	for p.At().Type != lexar.CloseBracket {
		if p.At().Type == lexar.Comma {
			p.Next()
			continue
		}
		params = append(params, p.ParseAssignmentExpression())
	}
	p.Except(lexar.CloseBracket, "Except close bracket Array")
	return params
}

func (p *parser) ParseLogicalExpression() any {
	left := p.ParseRelationalExpression()

	for utils.InArray(p.At().Value, []any{"&&", "||"}) {
		operator := p.Next().Value
		right := p.ParseRelationalExpression()

		left = &BinaryNode{
			Line:     p.At().Line,
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

		left = &BinaryNode{
			Line:     p.At().Line,
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
		left = &BinaryNode{
			Line:     p.At().Line,
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left
}

func (p *parser) ParseMultiplicativeExpression() any {
	left := p.ParseComputedMember()

	for p.At().Value == "*" || p.At().Value == "/" || p.At().Value == "%" {
		operator := p.Next().Value
		right := p.ParseComputedMember()
		left = &BinaryNode{
			Line:     p.At().Line,
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
	declatation := &VariableDeclarationNode{
		Line:  p.At().Line,
		Name:  value,
		Value: p.ParseAssignmentExpression(),
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
	return &ForNode{
		Line:      p.At().Line,
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

func (p *parser) ParseComputedMember() any {
	left := p.ParseCallMemberExpression(nil)
	for p.At().Type == lexar.OpenBracket {
		left = p.ParseCallMemberExpression(left)
	}
	return left
}

func (p *parser) ParseCallMemberExpression(left any) any {
	member := p.ParseMemberExpression(left)
	if p.At().Type == lexar.OpenParen {
		return p.ParseCallExpression(member)
	}
	return member
}

func (p *parser) ParseMemberExpression(left any) any {
	if left == nil {
		left = p.ParsePrimaryExpression()
	}
	if p.At().Type == lexar.OpenParen {
		left = p.ParseCallExpression(left)
	}
	for p.At().Type == lexar.OpenBracket {
		p.Next()
		left = &MemberNode{
			Line:     p.At().Line,
			Left:     left,
			Prop:     p.ParseAssignmentExpression(),
			Computed: true,
		}
		p.Except(lexar.CloseBracket, "Except close bracket Member")
	}
	for p.At().Type == lexar.Dot {
		p.Next()
		left = &MemberNode{
			Line:     p.At().Line,
			Left:     left,
			Prop:     p.ParsePrimaryExpression(),
			Computed: false,
		}
	}
	return left
}

func (p *parser) ParseCallExpression(caller any) any {
	return &CallNode{
		Line:   p.At().Line,
		Caller: caller,
		Args:   p.ParseArgs(),
	}
}

func (p *parser) ParsePrimaryExpression() any {
	switch p.At().Type {
	case lexar.Number:
		return &NumberLiteralNode{
			Line:  p.At().Line,
			Value: p.Next().Value,
		}
	case lexar.Float:
		return &FloatLiteralNode{
			Value: p.Next().Value,
		}
	case lexar.String:
		return &StringLiteralNode{
			Line:  p.At().Line,
			Value: p.Next().Value,
		}
	case lexar.Identifier:
		return &IdentifierNode{
			Line:  p.At().Line,
			Value: p.Next().Value,
		}
	case lexar.Semicolon:
		p.Next()
	default:
		p.At()
	}
	return 0
}
func (p *parser) ParseClassStatement() any {
	p.Next()
	var methods []*FunctionDeclarationNode
	identifier := p.Except(lexar.Identifier, "Except Identifier Class")
	args := p.ParseArgs()
	p.Except(lexar.OpenBrace, "Except open brace class")
	for p.At().Type != lexar.CloseBrace {
		methods = append(methods, p.ParseFnDeclaration())
	}
	p.Except(lexar.CloseBrace, "Except close brace class")

	return &ClassDeclarationNode{
		Line:    p.At().Line,
		Methods: methods,
		Name: &IdentifierNode{
			Value: identifier.Value,
		},
		Extends: args,
	}
}

func (p *parser) ParseFnDeclaration() *FunctionDeclarationNode {
	p.Next()
	identifier := p.Except(lexar.Identifier, "Funcsiya nomi nato'g'ri")
	args := p.ParseArgs()
	var params []any
	for _, arg := range args {
		switch arg.(type) {
		case *IdentifierNode, *AssignmentNode:
			params = append(params, arg)
		default:
			log.Fatal("Funcsiya parametri nato'g'ri")
		}
	}
	return &FunctionDeclarationNode{
		Line:   p.At().Line,
		Name:   identifier.Value.(string),
		Params: params,
		Body:   p.ParseBody(),
	}
}

func (p *parser) ParseImportStatement() any {
	p.Next()
	var name string
	module := p.ParseAssignmentExpression().(*StringLiteralNode).Value.(string)
	path := strings.Replace(module, ".", "/", -1)
	nameSegments := strings.Split(module, ".")
	if p.At().Type == lexar.As {
		p.Next()
		name = p.Next().Value.(string)
	} else {
		name = nameSegments[len(nameSegments)-1]
	}
	if utils.InArray(module, STDLIBS) {
		return &StdModuleNode{
			Line: p.At().Line,
			Name: name,
			Path: module,
		}
	}
	readCode := func(path string) []byte {
		readFile := func(filePath string) ([]byte, error) {
			data, err := assets.LibsFS.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
		if !utils.FileExists(p.BaseDir+path+".bug") && !utils.IsDirectory(p.BaseDir+path) {
			packagePath := "./packages/" + path
			if utils.FileExists(packagePath) {
				path = packagePath
			}
			data, err := readFile("libs/" + path + ".bug")
			if err == nil {
				return data
			}
			if _, err := assets.LibsFS.ReadDir("libs/" + path); err == nil {
				data, err := readFile("libs/" + path + "/init.bug")
				if err == nil {
					return data
				}
			}
		}
		if utils.IsDirectory(p.BaseDir + path) {
			path = path + "/init"
		}
		code, err := os.ReadFile(p.BaseDir + path + ".bug")
		if err != nil {
			panic("Error reading code: " + err.Error())
		}

		return code
	}
	tokenizer := lexar.NewTokenize()
	tokens, err := tokenizer.Tokenize(string(readCode(path)))
	if err != nil {
		return nil
	}
	parser := NewParser(tokens, p.BaseDir)
	program := &ModuleNode{
		Line: p.At().Line,
		Name: name,
		Path: path,
		Body: []any{},
	}
	for parser.IsEOF() {
		stmt := parser.ParseStatement()
		program.Body = append(program.Body, stmt)
	}
	return program
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
	case lexar.Class:
		return p.ParseClassStatement()
	case lexar.Import:
		return p.ParseImportStatement()
	default:
		return p.ParseAssignmentExpression()
	}
}

func (p *parser) CreateAST() any {
	program := &ProgramNode{
		Line: p.At().Line,
		Body: []any{},
	}
	for p.IsEOF() {
		stmt := p.ParseStatement()
		program.Body = append(program.Body, stmt)
	}
	return program
}
