package types

type ArrayValue struct {
	Values []any
}

func NewArray(value []any) *ArrayValue {
	return &ArrayValue{
		Values: value,
	}
}

func (a *ArrayValue) GetValue() []any {
	return a.Values
}

func (a *ArrayValue) Add(value any) {
	a.Values = append(a.Values, value)
}

func (a *ArrayValue) Pop() {
	a.Values = a.Values[:len(a.Values)-1]
}

func (a *ArrayValue) Size() any {
	return NewInt(len(a.Values))
}
