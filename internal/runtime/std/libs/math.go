package libs

import (
	"math"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func Round(values ...any) any {
	value := values[0].(*types.FloatValue)
	decimals := 1.0
	if len(values) > 1 {
		v, _ := utils.Int2Float(values[1].(*types.IntValue).Value)
		decimals = math.Pow(10.0, v)
	}
	return types.NewFloat(math.Round(value.Value*decimals) / decimals)
}

func Pow(base any, exponet any) any {
	b, _ := utils.Int2Float(base)
	e, _ := utils.Int2Float(exponet)
	return types.NewFloat(math.Pow(b, e))
}
