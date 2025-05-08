package types

type StdLibValue struct {
	Name string
	Lib  map[string]any
}

func NewStdLib(name string, lib map[string]any) Object {
	return &StdLibValue{
		Name: name,
		Lib:  lib,
	}
}

func (a *StdLibValue) GetValue() any {
	return a.Lib
}
