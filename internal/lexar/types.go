package lexar

type TokenType int

const (
	// Literal types
	Number TokenType = iota
	Float
	Identifier
	String
	Import

	// Grouping operators
	BinaryOperator
	Equals
	OpenParen    // (
	CloseParen   // )
	OpenBrace    // {
	CloseBrace   // }
	OpenBracket  // [
	CloseBracket // ]
	Colon
	Semicolon
	Comma
	Dot
	EOF // end of file

	// Keywords
	Var
	Const
	Class
	Fn
	Return
	If
	ElseIf
	Else
	New
	As

	// Conditions / loops
	For
	Continue
	Break
)

type Token struct {
	Value any
	Type  TokenType
	Line  int
}

var Keywords = map[string]TokenType{
	"var":      Var,
	"const":    Const,
	"=":        Equals,
	"func":     Fn,
	"return":   Return,
	"if":       If,
	"elseif":   ElseIf,
	"else":     Else,
	"&&":       BinaryOperator,
	"||":       BinaryOperator,
	"for":      For,
	"continue": Continue,
	"break":    Break,
	"class":    Class,
	"new":      New,
	"import":   Import,
	"as":       As,
}
