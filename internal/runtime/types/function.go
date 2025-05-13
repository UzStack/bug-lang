package types

type NativeFunctionValue struct {
	Call any
}

func NewNativeFunction(call any) Object {
	return &NativeFunctionValue{
		Call: call,
	}
}

func (a *NativeFunctionValue) GetValue() any {
	return a.Call
}
