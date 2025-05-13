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
	if len(results) >= 1 {
		switch v := results[0].Interface().(type) {
		case int:
			return types.NewInt(v)
		case uint:
			return types.NewInt(int(v))
		case string:
			return types.NewString(v)
		default:
			return v
		}
	}
	return types.NewNull()
}
