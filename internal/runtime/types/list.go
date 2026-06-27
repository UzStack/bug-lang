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

func eq(a any, value Object) bool {
	obj, ok := a.(Object)
	if !ok {
		return false
	}
	return obj.GetValue() == value.GetValue()
}

func (a *ArrayValue) Remove(value Object) {
	a.Values = slices.DeleteFunc(a.Values, func(a any) bool {
		return eq(a, value)
	})
}

func (a *ArrayValue) Index(value Object) int {
	return slices.IndexFunc(a.Values, func(a any) bool {
		return eq(a, value)
	})
}

func (a *ArrayValue) Contains(value Object) bool {
	return slices.ContainsFunc(a.Values, func(a any) bool {
		return eq(a, value)
	})
}

func (a *ArrayValue) Pop() {
	if len(a.Values) == 0 {
		return
	}
	a.Values = a.Values[:len(a.Values)-1]
}

func (a *ArrayValue) Size() any {
	return NewInt(len(a.Values))
}
