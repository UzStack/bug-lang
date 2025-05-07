package types

import "strings"

type StringValue struct {
	Value any
}

func NewString(value any) Object {
	return &StringValue{
		Value: value,
	}
}

func (a *StringValue) GetValue() any {
	return a.Value
}

func (a *StringValue) Upper() any {
	return NewString(strings.ToUpper(a.Value.(string)))
}
