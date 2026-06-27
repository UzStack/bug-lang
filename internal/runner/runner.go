// Package runner bitta .bug faylni HTTP so'rovi kontekstida ishlatadi.
// Ham FastCGI (bug-fpm), ham oddiy HTTP dev server (bug-serve) shuni ishlatadi.
package runner

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/UzStack/bug-lang/cmd/fpm/services"
	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

type Header struct {
	Key   string
	Value string
}

type Response struct {
	Body    string
	Headers []Header
	Status  int
}

// Run faylni ishlatadi va javobni qaytaradi. panic'ni 500 ga aylantiradi.
func Run(file, baseDir string, r *http.Request) (resp Response) {
	defer func() {
		if rec := recover(); rec != nil {
			resp = Response{Body: fmt.Sprintf("500 Internal Error: %v", rec), Status: 500}
		}
	}()

	code, err := os.ReadFile(file)
	if err != nil {
		return Response{Body: fmt.Sprintf("404 Not Found: %s", file), Status: 404}
	}
	tokenize := lexar.NewTokenize()
	tokens, err := tokenize.Tokenize(string(code))
	if err != nil {
		return Response{Body: err.Error(), Status: 500}
	}
	ast := parser.NewParser(tokens, baseDir).CreateAST()
	env := enviroment.NewGlobalEnv()
	std.Load(env)

	var buf bytes.Buffer
	var headers []Header
	status := 200

	env.AssignmenVariable("print", types.NewNativeFunction(func(values ...any) {
		std.Pprint(&buf, values)
	}), -1)
	env.AssignmenVariable("println", types.NewNativeFunction(func(values ...any) {
		std.Pprint(&buf, values)
		fmt.Fprint(&buf, "\n")
	}), -1)
	env.AssignmenVariable("header", types.NewNativeFunction(func(key any, value any) {
		headers = append(headers, Header{Key: asStr(key), Value: asStr(value)})
	}), -1)
	env.AssignmenVariable("status", types.NewNativeFunction(func(values ...any) any {
		if len(values) > 0 {
			if code, ok := values[0].(*types.IntValue); ok {
				status = code.Value
			}
		}
		return types.NewNull()
	}), -1)
	env.AssignmenVariable("cookie", types.NewNativeFunction(func(values ...any) any {
		if len(values) >= 2 {
			headers = append(headers, Header{
				Key:   "Set-Cookie",
				Value: asStr(values[0]) + "=" + asStr(values[1]) + "; Path=/",
			})
		}
		return types.NewNull()
	}), -1)
	env.AssignmenVariable("redirect", types.NewNativeFunction(func(values ...any) any {
		if len(values) > 0 {
			headers = append(headers, Header{Key: "Location", Value: asStr(values[0])})
			status = 302
		}
		return types.NewNull()
	}), -1)

	env.AssignmenVariable("_POST", services.ParsePostData(r), -1)
	env.AssignmenVariable("_GET", services.ParseGetData(r), -1)
	env.AssignmenVariable("_REQUEST", services.ParseRequest(r), -1)
	env.AssignmenVariable("_GLOBALS", services.ParseGlobals(r), -1)

	if _, err := runtime.Interpreter(ast, env); err != nil {
		return Response{Body: err.Error(), Status: 500}
	}
	return Response{Body: buf.String(), Headers: headers, Status: status}
}

// Apply javobni http.ResponseWriter ga yozadi.
func (resp Response) Apply(w http.ResponseWriter) {
	for _, h := range resp.Headers {
		if h.Key == "Set-Cookie" {
			w.Header().Add(h.Key, h.Value)
		} else {
			w.Header().Set(h.Key, h.Value)
		}
	}
	if resp.Status != 0 && resp.Status != 200 {
		w.WriteHeader(resp.Status)
	}
	w.Write([]byte(resp.Body))
}

func asStr(v any) string {
	if s, ok := v.(*types.StringValue); ok {
		return s.Value
	}
	if obj, ok := v.(types.Object); ok {
		return fmt.Sprintf("%v", obj.GetValue())
	}
	return fmt.Sprintf("%v", v)
}
