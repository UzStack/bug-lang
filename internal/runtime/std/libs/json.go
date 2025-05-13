package libs

import (
	"encoding/json"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func JsonEncode(value any) any {
	data, err := json.Marshal(utils.EncodeBug(value))
	if err != nil {
		panic(err.Error())
	}
	return types.NewString(string(data))
}
