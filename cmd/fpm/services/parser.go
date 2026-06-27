package services

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

const maxBodySize = 10 << 20 // 10MB — cheksiz body DoS oldini olish

func ParsePostData(request *http.Request) any {
	body, err := io.ReadAll(io.LimitReader(request.Body, maxBodySize))
	if err != nil {
		return types.NewMap(make(map[string]any))
	}
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return types.NewMap(make(map[string]any))
	}
	return utils.DecodeBug(data)
}

func ParseGetData(request *http.Request) any {
	data := types.NewMap(make(map[string]any)).(*types.MapValue)
	// URL.Query() percent-decode va '=' belgisini to'g'ri ishlaydi
	for key, values := range request.URL.Query() {
		if len(values) > 0 {
			data.Append(types.NewString(key), utils.DecodeBug(values[0]))
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
