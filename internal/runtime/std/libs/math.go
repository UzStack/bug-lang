package libs

import (
	"math"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func Round(values ...any) any {
	if len(values) == 0 {
		panic("math.round: argument kerak")
	}
	value, err := utils.Int2Float(values[0])
	if err != nil {
		panic("math.round: " + err.Error())
	}
	decimals := 1.0
	if len(values) > 1 {
		v, _ := utils.Int2Float(values[1])
		decimals = math.Pow(10.0, v)
	}
	return types.NewFloat(math.Round(value*decimals) / decimals)
}

func Pow(base any, exponet any) any {
	b, _ := utils.Int2Float(base)
	e, _ := utils.Int2Float(exponet)
	return types.NewFloat(math.Pow(b, e))
}
