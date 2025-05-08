package types

import "strings"

type StringValue struct {
	Value string
}

func NewString(value string) Object {
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
