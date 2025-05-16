package types

import (
	"slices"
)

type ArrayValue struct {
	Values []any
}

func NewArray(value []any) Object {
	return &ArrayValue{
		Values: value,
	}
}

func (a *ArrayValue) GetValue() any {
	return a.Values
}

func (a *ArrayValue) Append(value any) {
	a.Values = append(a.Values, value)
}

func (a *ArrayValue) Extend(value *ArrayValue) {
	a.Values = append(a.Values, value.GetValue().([]any)...)
}

func (a *ArrayValue) Remove(value Object) {
	a.Values = slices.DeleteFunc(a.Values, func(a any) bool {
		return a.(Object).GetValue() == value.GetValue()
	})
}

func (a *ArrayValue) Index(value Object) int {
	return slices.IndexFunc(a.Values, func(a any) bool {
		return a.(Object).GetValue() == value.GetValue()
	})
}

func (a *ArrayValue) Contains(value Object) bool {
	return slices.ContainsFunc(a.Values, func(a any) bool {
		return a.(Object).GetValue() == value.GetValue()
	})
}

func (a *ArrayValue) Pop() {
	a.Values = a.Values[:len(a.Values)-1]
}

func (a *ArrayValue) Size() any {
	return NewInt(len(a.Values))
}
