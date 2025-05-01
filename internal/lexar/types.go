package lexar

type TokenType int

const (
	// Literal types
	Number TokenType = iota
	Identifier
	String

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
	Let
	Const
	Fn
	Return
	If
	ElseIf
	Else

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
	"var":      Let,
	"const":    Const,
	"endi":     Equals,
	"def":      Fn,
	"return":   Return,
	"if":       If,
	"elseif":   ElseIf,
	"else":     Else,
	"&&":       BinaryOperator,
	"||":       BinaryOperator,
	"for":      For,
	"continue": Continue,
	"break":    Break,
}
