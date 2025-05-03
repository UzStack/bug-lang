package types

type ArrayValue struct {
	Values []any
}

func NewArray() Object {
	return &ArrayValue{}
}

func (a *ArrayValue) GetValue() any {
	return a.Values
}

func (a *ArrayValue) Add(value any) {
	a.Values = append(a.Values, value)
}

func (a *ArrayValue) Pop() {
	a.Values = a.Values[:len(a.Values)-1]
}
