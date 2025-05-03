package types

type ObjectValue struct {
	Values map[string]any
}

func NewObject(values map[string]any) Object {
	return &ObjectValue{
		Values: values,
	}
}

func (a *ObjectValue) GetValue() any {
	return a.Values
}
