package libs

import (
	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/k0kubun/pp/v3"
)

func Round(value *types.FloatValue) any {
	pp.Print(value.GetValue())
	return nil
}
