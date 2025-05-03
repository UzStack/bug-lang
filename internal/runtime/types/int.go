package types

type IntValue struct {
	Value int
}

func NewInt() Object {
	return &IntValue{}
}
func (a *IntValue) GetValue() any {
	return a.Value
}
