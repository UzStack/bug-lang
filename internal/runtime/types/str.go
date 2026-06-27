package types

import "strings"

type StringValue struct {
	Value string
}

func NewString(value string) *StringValue {
	return &StringValue{
		Value: value,
	}
}

func (a *StringValue) GetValue() any {
	return a.Value
}

func (a *StringValue) Upper() any {
	return NewString(strings.ToUpper(a.Value))
}

func (a *StringValue) Lower() any {
	return NewString(strings.ToLower(a.Value))
}

func (a *StringValue) Size() any {
	return NewInt(len([]rune(a.Value)))
}

func (a *StringValue) Trim() any {
	return NewString(strings.TrimSpace(a.Value))
}

func (a *StringValue) Split(sep *StringValue) any {
	parts := strings.Split(a.Value, sep.Value)
	values := make([]any, len(parts))
	for i, p := range parts {
		values[i] = NewString(p)
	}
	return NewArray(values)
}

func (a *StringValue) Replace(old *StringValue, replacement *StringValue) any {
	return NewString(strings.ReplaceAll(a.Value, old.Value, replacement.Value))
}

func (a *StringValue) Contains(sub *StringValue) any {
	return NewBool(strings.Contains(a.Value, sub.Value))
}

func (a *StringValue) IndexOf(sub *StringValue) any {
	return NewInt(strings.Index(a.Value, sub.Value))
}

func (a *StringValue) StartsWith(p *StringValue) any {
	return NewBool(strings.HasPrefix(a.Value, p.Value))
}

func (a *StringValue) EndsWith(s *StringValue) any {
	return NewBool(strings.HasSuffix(a.Value, s.Value))
}

// Substring(start, end) — rune asosida, chegaradan tashqari xavfsiz.
func (a *StringValue) Substring(start *IntValue, end *IntValue) any {
	runes := []rune(a.Value)
	s, e := start.Value, end.Value
	if s < 0 {
		s = 0
	}
	if e > len(runes) {
		e = len(runes)
	}
	if s > e {
		return NewString("")
	}
	return NewString(string(runes[s:e]))
}
