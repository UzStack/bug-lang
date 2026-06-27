package types

import "maps"

type MapValue struct {
	Data map[string]any
}

func NewMap(values map[string]any) Object {
	return &MapValue{
		Data: values,
	}
}

func (a *MapValue) GetValue() any {
	return a.Data
}

func (a *MapValue) Contains(value Object) bool {
	key, ok := value.GetValue().(string)
	if !ok {
		return false
	}
	_, ok = a.Data[key]
	return ok
}

func (a *MapValue) Remove(value Object) {
	maps.DeleteFunc(a.Data, func(key string, v any) bool {
		return key == value.GetValue()
	})
}

func (a *MapValue) Append(key *StringValue, value any) any {
	a.Data[key.GetValue().(string)] = value
	return nil
}

func (a *MapValue) Size() any {
	return NewInt(len(a.Data))
}

func (a *MapValue) Keys() any {
	values := make([]any, 0, len(a.Data))
	for k := range a.Data {
		values = append(values, NewString(k))
	}
	return NewArray(values)
}

func (a *MapValue) Values() any {
	values := make([]any, 0, len(a.Data))
	for _, v := range a.Data {
		values = append(values, v)
	}
	return NewArray(values)
}
