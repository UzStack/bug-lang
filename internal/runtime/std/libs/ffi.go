package libs

import (
	"plugin"
	"reflect"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func FFILoad(lib *types.StringValue) any {
	plg, err := plugin.Open(lib.Value)
	if err != nil {
		panic(err.Error())
	}
	return plg
}

func FFICall(plg *plugin.Plugin, symbolName *types.StringValue, args *types.ArrayValue) any {
	sym, err := plg.Lookup(symbolName.Value)
	if err != nil {
		panic(err.Error())
	}
	v := reflect.ValueOf(sym)
	callArgs := make([]reflect.Value, len(args.Values))
	for i, arg := range args.Values {
		callArgs[i] = reflect.ValueOf(arg)
	}
	results := v.Call(callArgs)
	result := make([]any, len(results))
	for i, res := range results {
		result[i] = res.Interface()
	}
	return types.NewArray(result)
}
