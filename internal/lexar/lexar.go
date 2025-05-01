package lexar

import (
	"strings"

	"github.com/k0kubun/pp"
)

type tokenize struct {
	Tokens   []*Token
	Line     int
	TempWord string
}

func NewTokenize() *tokenize {
	return &tokenize{
		Line: 1,
	}
}

func (tokenize *tokenize) token(value any, tokenType TokenType) {
	tokenize.Tokens = append(tokenize.Tokens, &Token{
		Value: value,
		Type:  tokenType,
		Line:  tokenize.Line,
	})
}

func (t *tokenize) Handle() {
	if IsAlpha(t.TempWord) {
		t.token(t.TempWord, String)
	} else if IsNumber(t.TempWord) {
		t.token(t.TempWord, Number)
	}
	t.TempWord = ""
}

func (tokenize tokenize) Get() []*Token {
	return tokenize.Tokens
}

func (t *tokenize) Tokenize(code string) int {
	chars := strings.Split(code, "")
	for _, char := range chars {
		if char == "\n" {
			t.Line++
		}
		if char == ":" {
			t.token(char, Colon)
		}
		if char == "=" {
			t.Handle()
			t.token(char, Equals)
		}
		if char == "(" {
			t.Handle()
			t.token(char, OpenParen)
		}
		if char == ")" {
			t.Handle()
			t.token(char, CloseParen)
		}
		if char == "{" {
			t.Handle()
			t.token(char, OpenBrace)
		}
		if char == "}" {
			t.Handle()
			t.token(char, CloseBrace)
		}
		if char == "[" {
			t.Handle()
			t.token(char, CloseBracket)
		}
		if char == "]" {
			t.Handle()
			t.token(char, CloseBracket)
		}
		if IsAlpha(char) {
			t.TempWord += char
			continue
		}
		if IsNumber(char) {
			t.TempWord += char
			continue
		}
		t.Handle()
	}
	t.Handle()
	t.token("EndOfLine", EOF)
	pp.Println(t.Get())
	return 0
}
