package std

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/UzStack/bug-lang/internal/runtime/std/libs"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

var STDLIBS = map[string]map[string]any{
	"math": {
		"round": libs.Round,
		"pow":   libs.Pow,
	},
	"ffi": {
		"load": libs.FFILoad,
		"call": libs.FFICall,
	},
	"json": {
		"encode": libs.JsonEncode,
	},
}

func QuotationMark(value any) bool {
	switch value.(type) {
	case *types.StringValue:
		return true
	default:
		return false
	}
}

func Header(key any, value any) {
	// Not implement
}

func Pprint(buf *bytes.Buffer, values ...any) {
	for _, val := range values {
		switch v := val.(type) {
		case *types.StringValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.IntValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.FloatValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.BoolValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.ArrayValue:
			fmt.Fprint(buf, "[")
			for index, el := range v.Values {
				Pprint(buf, el)
				if index < len(v.Values)-1 {
					fmt.Fprint(buf, ",")
				}
			}
			fmt.Fprint(buf, "]")
		case *types.MapValue:
			fmt.Fprint(buf, "{")
			i := 0
			values := v.GetValue().(map[string]any)
			size := len(values)
			for key, value := range values {
				i++
				fmt.Fprint(buf, "\"", key, "\"", ":")
				isQuotationMark := QuotationMark(value)
				if isQuotationMark {
					fmt.Fprint(buf, "\"")
				}
				Pprint(buf, value)
				if isQuotationMark {
					fmt.Fprint(buf, "\"")
				}
				if size != i {
					fmt.Fprint(buf, ",")
				}
			}
			fmt.Fprint(buf, "}")
		default:
			refValue := reflect.ValueOf(val)
			if refValue.Kind() == reflect.Slice {
				for i := 0; i < refValue.Len(); i++ {
					Pprint(buf, refValue.Index(i).Interface())
				}
			} else {
				fmt.Print(val)
			}
		}
	}
}
func Println(values ...any) {
	Print(values)
	fmt.Print("\t\n")
}

func Print(values ...any) {
	var buf bytes.Buffer
	Pprint(&buf, values...)
	fmt.Print(buf.String())
}

func Input(values ...any) any {
	fmt.Print(values[0].(*types.RuntimeValue).Value)
	reader := bufio.NewReader(os.Stdin)
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)
	return &types.RuntimeValue{
		Value: data,
	}
}

func Super(values any) any {
	return values
}
