package types

type IntValue struct {
	Value int
}

func NewInt(value int) Object {
	return &IntValue{
		Value: value,
	}
}
func (a *IntValue) GetValue() any {
	return a.Value
}
