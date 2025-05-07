package types

type BoolValue struct {
	Value any
}

func NewBool(value any) Object {
	return &BoolValue{
		Value: value,
	}
}
func (a *BoolValue) GetValue() any {
	return a.Value
}
