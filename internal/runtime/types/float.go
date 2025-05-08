package types

type FloatValue struct {
	Value float64
}

func NewFloat(value float64) Object {
	return &FloatValue{
		Value: value,
	}
}
func (a *FloatValue) GetValue() any {
	return a.Value
}
