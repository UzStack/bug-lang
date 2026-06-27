package libs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/UzStack/bug-lang/internal/runtime/types"
	"github.com/UzStack/bug-lang/pkg/utils"
)

// ponytail: SSRF/ichki-tarmoq filtri yo'q — timeout bilan cheklandi.
// Ishonchsiz skript ishlatilsa, IP allowlist qo'shish kerak.
func Request(method *types.StringValue, url *types.StringValue, p types.Object) any {
	client := http.Client{Timeout: 30 * time.Second}
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
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	return types.NewMap(map[string]any{
		"body":  body,
		"error": nil,
	})
}

func Json(response types.Object) any {
	m, ok := response.GetValue().(map[string]any)
	if !ok {
		panic("requests.json: response map emas")
	}
	raw, ok := m["body"].([]byte)
	if !ok {
		panic("requests.json: body o'qib bo'lmadi")
	}
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		panic("requests.json: " + err.Error())
	}
	return utils.DecodeBug(data)
}
