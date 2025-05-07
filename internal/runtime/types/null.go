package types

type NullValue struct {
	Value any
}

func NewNull(value any) Object {
	return &NullValue{
		Value: value,
	}
}
func (a *NullValue) GetValue() any {
	return a.Value
}
