package std

import "fmt"

func Print(values ...any) {
	for _, val := range values {
		fmt.Print(val)
	}
	fmt.Print("\t\n")
}
