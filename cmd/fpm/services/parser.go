package services

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func ParsePostData(request *http.Request) any {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err.Error())
	}
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return types.NewMap(nil)
	}
	return utils.DecodeBug(data)
}

func ParseGetData(request *http.Request) any {
	data := types.NewMap(make(map[string]any)).(*types.MapValue)
	rawQuery := request.URL.RawQuery
	if len(rawQuery) == 0 {
		return data
	}
	for _, p := range strings.Split(rawQuery, "&") {
		param := strings.Split(p, "=")
		if len(param) >= 2 {
			data.Append(types.NewString(param[0]), utils.DecodeBug(param[1]))
		}
	}
	return data
}

func ParseRequest(request *http.Request) any {
	headers := types.NewMap(make(map[string]any)).(*types.MapValue)
	for key, values := range request.Header {
		if len(values) > 1 {
			value := types.NewArray([]any{}).(*types.ArrayValue)
			for _, v := range values {
				value.Append(utils.DecodeBug(v))
			}
			headers.Append(types.NewString(key), value)
		} else if len(values) == 1 {
			headers.Append(types.NewString(key), utils.DecodeBug(values[0]))
		} else {
			headers.Append(types.NewString(key), types.NewNull())
		}
	}

	globals := types.NewMap(make(map[string]any)).(*types.MapValue)
	globals.Append(types.NewString("RequestURI"), types.NewString(request.RequestURI))
	globals.Append(types.NewString("Host"), types.NewString(request.Host))
	globals.Append(types.NewString("Method"), types.NewString(request.Method))
	globals.Append(types.NewString("Headers"), headers)
	globals.Append(types.NewString("Path"), types.NewString(request.URL.Path))
	return globals
}

func ParseGlobals(request *http.Request) any {
	globals := types.NewMap(make(map[string]any))
	return globals
}
