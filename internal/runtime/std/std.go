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
		}
	}
	fmt.Print("\t\n")
}
