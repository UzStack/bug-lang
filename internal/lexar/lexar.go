package lexar

import (
	"fmt"
	"strconv"
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
	} else if IsSignedFloat(t.TempWord) {
		t.token(t.TempWord, Float)
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

func (t *tokenize) FindString() error {
	var str string
	for t.At() != "\"" || t.Chars[t.Index-1] == "\\" {
		str += t.Next()
	}
	str, err := strconv.Unquote("\"" + str + "\"")
	if err != nil {
		return fmt.Errorf("string unquote error: %s", err.Error())
	}
	t.token(str, String)
	t.Index++
	return nil
}

func (t *tokenize) FindEnd() {
	for len(t.Chars) != t.Index+1 && t.At() != ";" {
		t.Next()
	}
	t.Next()
}

func (t *tokenize) Tokenize(code string) ([]*Token, error) {
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
			if err := t.FindString(); err != nil {
				return nil, err
			}
			continue
		}
		if char == "/" && t.Chars[t.Index] == "/" {
			t.Handle()
			t.FindEnd()
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
			t.token(char, OpenBracket)
		}
		if char == "]" {
			t.Handle()
			t.token(char, CloseBracket)
		}

		if char == "." && !IsSignedFloat(t.TempWord) {
			t.Handle()
			t.token(char, Dot)
		}

		if IsNumber(char) && strings.Count(t.TempWord, ".") == 0 {
			t.TempWord += char
			continue
		}

		if (IsNumber(char) || char == ".") && IsSignedFloat(t.TempWord) {
			t.TempWord += char
			continue
		}

		if IsAlpha(char) {
			t.TempWord += char
			continue
		}

		if utils.InArray(char, []any{"=", "<", ">", "!", "&", "|"}) {
			t.Handle()
			if utils.InArray(t.At(), []any{"&", "|", "="}) {
				t.token(char+t.Next(), BinaryOperator)
			} else if char == "=" {
				t.token(char, Equals)
			} else {
				t.token(char, BinaryOperator)
			}
		}

		t.Handle()
	}
	t.Handle()
	t.token("EndOfLine", EOF)
	return t.Get(), nil
}
