package types

type NullValue struct {
	Value any
}

func NewNull() Object {
	return &NullValue{
		Value: nil,
	}
}
func (a *NullValue) GetValue() any {
	return a.Value
}
