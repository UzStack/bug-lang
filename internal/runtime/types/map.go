package types

import "maps"

type MapValue struct {
	Values map[string]any
}

func NewMap(values map[string]any) Object {
	return &MapValue{
		Values: values,
	}
}

func (a *MapValue) GetValue() any {
	return a.Values
}

func (a *MapValue) Contains(value Object) bool {
	_, ok := a.Values[value.GetValue().(string)]
	return ok
}

func (a *MapValue) Remove(value Object) {
	maps.DeleteFunc(a.Values, func(key string, v any) bool {
		return key == value.GetValue()
	})
}

func (a *MapValue) Append(key *StringValue, value any) any {
	a.Values[key.GetValue().(string)] = value
	return nil
}
