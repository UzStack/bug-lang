package types

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

func (a *MapValue) Add(key *StringValue, value any) any {
	a.Values[key.GetValue().(string)] = value
	return nil
}
