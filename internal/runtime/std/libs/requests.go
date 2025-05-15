package libs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

func Request(method *types.StringValue, url *types.StringValue, p types.Object) any {
	client := http.Client{}
	var payload *bytes.Reader
	if p.GetValue() != nil {
		payload = bytes.NewReader([]byte(utils.EncodeBug(p.GetValue()).(string)))
	} else {
		payload = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method.GetValue().(string), url.GetValue().(string), payload)
	if err != nil {
		panic(err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	return types.NewMap(map[string]any{
		"body":  body,
		"error": err,
	})
}

func Json(response types.Object) any {
	var data any
	json.Unmarshal(response.GetValue().(map[string]any)["body"].([]byte), &data)
	return utils.DecodeBug(data)
}
