package types

type StringValue struct {
	Value string
}

func NewString() Object {
	return &StringValue{}
}

func (a *StringValue) GetValue() any {
	return a.Value
}
