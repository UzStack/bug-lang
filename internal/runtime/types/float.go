package types

type FloatValue struct {
	Value any
}

func NewFloat(value any) Object {
	return &FloatValue{
		Value: value,
	}
}
func (a *FloatValue) GetValue() any {
	return a.Value
}
