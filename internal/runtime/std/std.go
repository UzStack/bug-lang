package std

import (
	"fmt"

	"github.com/UzStack/bug-lang/internal/parser"
)

func Print(values ...any) {
	for _, val := range values {
		switch v := val.(type) {
		case *parser.NumberLiteralNode:
			fmt.Print(v.Value)
		case *parser.StringLiteralNode:
			fmt.Print(v.Value)
		}
	}
	fmt.Print("\t\n")
}
