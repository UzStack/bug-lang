package libs

import (
	"encoding/json"
	"fmt"
	"strconv"

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

func JsonDecode(value *types.StringValue) any {
	var data any
	str, err := strconv.Unquote("\"" + value.GetValue().(string) + "\"")
	if err != nil {
		panic(fmt.Sprintf("string unquote error: %s", err.Error()))
	}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		panic(err.Error())
	}
	return utils.DecodeBug(data)
}
