package std

import (
	"fmt"

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
		}
	}
	fmt.Print("\t\n")
}
