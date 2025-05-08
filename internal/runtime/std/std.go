package std

import (
	"bufio"
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

func Pprint(values ...any) {
	for _, val := range values {
		switch v := val.(type) {
		case *types.StringValue:
			fmt.Print(v.GetValue())
		case *types.IntValue:
			fmt.Print(v.GetValue())
		case *types.FloatValue:
			fmt.Print(v.GetValue())
		case *types.ArrayValue:
			fmt.Print("[")
			for index, el := range v.Values {
				Pprint(el)
				if index < len(v.Values)-1 {
					fmt.Print(",")
				}
			}
			fmt.Print("]")
		case *types.MapValue:
			fmt.Print("{")
			i := 0
			values := v.GetValue().(map[string]any)
			size := len(values)
			for key, value := range values {
				i++
				fmt.Print("\"", key, "\"", ":")
				isQuotationMark := QuotationMark(value)
				if isQuotationMark {
					fmt.Print("\"")
				}
				Pprint(value)
				if isQuotationMark {
					fmt.Print("\"")
				}
				if size != i {
					fmt.Print(",")
				}
			}
			fmt.Print("}")
		default:
			refValue := reflect.ValueOf(val)
			if refValue.Kind() == reflect.Slice {
				for i := 0; i < refValue.Len(); i++ {
					Print(refValue.Index(i).Interface())
				}
			}
		}
	}
}
func Print(values ...any) {
	Pprint(values...)
	fmt.Print("\t\n")
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
