package lexar

import (
	"strings"

	"github.com/UzStack/bug-lang/pkg/utils"
)

type tokenize struct {
	Tokens   []*Token
	Line     int
	TempWord string
	Chars    []string
	Index    int
}

func NewTokenize() *tokenize {
	return &tokenize{
		Line:  1,
		Chars: []string{},
		Index: 0,
	}
}

func (tokenize *tokenize) SaveAlpha() {
	keyword, ok := Keywords[tokenize.TempWord]
	if ok {
		tokenize.token(tokenize.TempWord, keyword)
		return
	}
	tokenize.token(tokenize.TempWord, Identifier)
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
		t.SaveAlpha()
	} else if IsNumber(t.TempWord) {
		t.token(t.TempWord, Number)
	}
	t.TempWord = ""
}

func (tokenize tokenize) Get() []*Token {
	return tokenize.Tokens
}
func (t *tokenize) Next() string {
	t.Index++
	return t.Chars[t.Index-1]
}
func (t *tokenize) At() string {
	return t.Chars[t.Index]
}

func (t *tokenize) FindString() {
	var str string
	for t.At() != "\"" {
		str += t.Next()
	}
	t.token(str, String)
	t.Index++
}

func (t *tokenize) Tokenize(code string) []*Token {
	t.Chars = strings.Split(code, "")
	for t.Index < len(t.Chars) {
		char := t.Next()
		if char == "\n" {
			t.Line++
		}
		if char == ":" {
			t.token(char, Colon)
		}
		if char == "\"" {
			t.Handle()
			t.FindString()
			continue
		}
		if utils.InArray(char, []any{"+", "-", "/", "*", "%"}) {
			t.Handle()
			t.token(char, BinaryOperator)
			continue
		}
		if char == ";" {
			t.Handle()
			t.token(char, Semicolon)
		}
		if char == "," {
			t.Handle()
			t.token(char, Comma)
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
	return t.Get()
}
