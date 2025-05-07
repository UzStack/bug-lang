package types

type IntValue struct {
	Value any
}

func NewInt(value any) Object {
	return &IntValue{
		Value: value,
	}
}
func (a *IntValue) GetValue() any {
	return a.Value
}
