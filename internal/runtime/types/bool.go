package types

type BoolValue struct {
	Value bool
}

func NewBool(value bool) Object {
	return &BoolValue{
		Value: value,
	}
}
func (a *BoolValue) GetValue() any {
	return a.Value
}
