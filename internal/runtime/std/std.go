package std

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func Print(values ...any) {
	for _, val := range values {
		switch v := val.(type) {
		case *types.RuntimeValue:
			fmt.Print(v.Value, " ")
		case *types.ArrayValue:
			fmt.Print("[")
			for index, el := range v.Values {
				fmt.Print(el.(*types.RuntimeValue).Value)
				if index < len(v.Values)-1 {
					fmt.Print(",")
				}
			}
			fmt.Print("]\n")
		default:
			refValue := reflect.ValueOf(val)
			if refValue.Kind() == reflect.Slice {
				for i := 0; i < refValue.Len(); i++ {
					Print(refValue.Index(i).Interface())
				}
			}
		}
	}
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
