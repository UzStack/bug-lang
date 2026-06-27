package std

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/UzStack/bug-lang/internal/runtime/std/libs"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

var STDLIBS = map[string]map[string]any{
	"math": {
		"round": libs.Round,
		"pow":   libs.Pow,
	},
	"ffi": {
		"load": libs.FFILoad,
		"call": libs.FFICall,
	},
	"json": {
		"encode": libs.JsonEncode,
		"decode": libs.JsonDecode,
	},
	"requestsgo": {
		"request": libs.Request,
		"json":    libs.Json,
	},
	"os": {
		"read":   libs.OsRead,
		"write":  libs.OsWrite,
		"append": libs.OsAppend,
		"exists": libs.OsExists,
		"remove": libs.OsRemove,
		"env":    libs.OsEnv,
		"args":   libs.OsArgs,
	},
	"regex": {
		"match":   libs.RegexMatch,
		"find":    libs.RegexFind,
		"replace": libs.RegexReplace,
	},
	"time": {
		"now":    libs.TimeNow,
		"unix":   libs.TimeUnix,
		"format": libs.TimeFormat,
		"sleep":  libs.TimeSleep,
	},
	"sqldriver": {
		"connect": libs.DBConnect,
		"query":   libs.DBQuery,
		"exec":    libs.DBExec,
		"close":   libs.DBClose,
	},
}

func QuotationMark(value any) bool {
	switch value.(type) {
	case *types.StringValue:
		return true
	default:
		return false
	}
}

func Header(key any, value any) {
	// CLI rejimida no-op; FPM worker haqiqiy versiyani o'rnatadi
}

// Web yordamchilari — CLI rejimida no-op, FPM'da qayta belgilanadi.
func Status(values ...any) any { return types.NewNull() }
func Redirect(values ...any) any { return types.NewNull() }
func Cookie(values ...any) any { return types.NewNull() }

func Pprint(buf *bytes.Buffer, values ...any) {
	for _, val := range values {
		switch v := val.(type) {
		case *types.StringValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.IntValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.FloatValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.BoolValue:
			fmt.Fprint(buf, v.GetValue())
		case *types.ArrayValue:
			fmt.Fprint(buf, "[")
			for index, el := range v.Values {
				Pprint(buf, el)
				if index < len(v.Values)-1 {
					fmt.Fprint(buf, ",")
				}
			}
			fmt.Fprint(buf, "]")
		case *types.MapValue:
			fmt.Fprint(buf, "{")
			i := 0
			values := v.GetValue().(map[string]any)
			size := len(values)
			for key, value := range values {
				i++
				fmt.Fprint(buf, "\"", key, "\"", ":")
				isQuotationMark := QuotationMark(value)
				if isQuotationMark {
					fmt.Fprint(buf, "\"")
				}
				Pprint(buf, value)
				if isQuotationMark {
					fmt.Fprint(buf, "\"")
				}
				if size != i {
					fmt.Fprint(buf, ",")
				}
			}
			fmt.Fprint(buf, "}")
		default:
			refValue := reflect.ValueOf(val)
			if refValue.Kind() == reflect.Slice {
				for i := 0; i < refValue.Len(); i++ {
					Pprint(buf, refValue.Index(i).Interface())
				}
			} else {
				fmt.Print(val)
			}
		}
	}
}
func Println(values ...any) {
	Print(values)
	fmt.Print("\t\n")
}

func Print(values ...any) {
	var buf bytes.Buffer
	Pprint(&buf, values...)
	fmt.Print(buf.String())
}

func Input(values ...any) any {
	if len(values) > 0 {
		if rv, ok := values[0].(*types.RuntimeValue); ok {
			fmt.Print(rv.Value)
		} else if obj, ok := values[0].(types.Object); ok {
			fmt.Print(obj.GetValue())
		}
	}
	reader := bufio.NewReader(os.Stdin)
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)
	return &types.RuntimeValue{
		Value: data,
	}
}

func Super(values any) any {
	return values
}

// toString ixtiyoriy qiymatni satrga aylantiradi.
func toString(value any) string {
	obj, ok := value.(types.Object)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	switch v := obj.GetValue().(type) {
	case nil:
		return "null"
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Len: string (rune soni), array yoki map uzunligi.
func Len(values ...any) any {
	if len(values) == 0 {
		return types.NewInt(0)
	}
	switch v := values[0].(type) {
	case *types.StringValue:
		return types.NewInt(len([]rune(v.Value)))
	case *types.ArrayValue:
		return types.NewInt(len(v.Values))
	case *types.MapValue:
		return types.NewInt(len(v.Data))
	}
	return types.NewInt(0)
}

func Str(values ...any) any {
	if len(values) == 0 {
		return types.NewString("")
	}
	return types.NewString(toString(values[0]))
}

func Int(values ...any) any {
	if len(values) == 0 {
		return types.NewInt(0)
	}
	obj, ok := values[0].(types.Object)
	if !ok {
		return types.NewInt(0)
	}
	switch v := obj.GetValue().(type) {
	case int:
		return types.NewInt(v)
	case float64:
		return types.NewInt(int(v))
	case bool:
		if v {
			return types.NewInt(1)
		}
		return types.NewInt(0)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			panic("int(): '" + v + "' butun songa aylanmaydi")
		}
		return types.NewInt(n)
	}
	return types.NewInt(0)
}

func Float(values ...any) any {
	if len(values) == 0 {
		return types.NewFloat(0)
	}
	obj, ok := values[0].(types.Object)
	if !ok {
		return types.NewFloat(0)
	}
	switch v := obj.GetValue().(type) {
	case int:
		return types.NewFloat(float64(v))
	case float64:
		return types.NewFloat(v)
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			panic("float(): '" + v + "' songa aylanmaydi")
		}
		return types.NewFloat(f)
	}
	return types.NewFloat(0)
}

func Bool(values ...any) any {
	if len(values) == 0 {
		return types.NewBool(false)
	}
	obj, ok := values[0].(types.Object)
	if !ok {
		return types.NewBool(false)
	}
	switch v := obj.GetValue().(type) {
	case nil:
		return types.NewBool(false)
	case bool:
		return types.NewBool(v)
	case int:
		return types.NewBool(v != 0)
	case float64:
		return types.NewBool(v != 0)
	case string:
		return types.NewBool(v != "")
	}
	return types.NewBool(true)
}

func TypeOf(values ...any) any {
	if len(values) == 0 {
		return types.NewString("null")
	}
	switch values[0].(type) {
	case *types.StringValue:
		return types.NewString("string")
	case *types.IntValue:
		return types.NewString("int")
	case *types.FloatValue:
		return types.NewString("float")
	case *types.BoolValue:
		return types.NewString("bool")
	case *types.ArrayValue:
		return types.NewString("list")
	case *types.MapValue:
		return types.NewString("dict")
	case *types.NullValue:
		return types.NewString("null")
	case *types.ObjectValue:
		return types.NewString("object")
	case *types.FunctionDeclaration, *types.NativeFunctionValue:
		return types.NewString("function")
	}
	return types.NewString("unknown")
}

// Format: "{}" o'rinbosarlarini ketma-ket argumentlar bilan almashtiradi.
func Format(values ...any) any {
	if len(values) == 0 {
		return types.NewString("")
	}
	tmpl := toString(values[0])
	for _, arg := range values[1:] {
		tmpl = strings.Replace(tmpl, "{}", toString(arg), 1)
	}
	return types.NewString(tmpl)
}

func Assert(values ...any) any {
	if len(values) == 0 {
		panic("assert(): shart kerak")
	}
	obj, ok := values[0].(types.Object)
	cond := ok && obj.GetValue() == true
	if !cond {
		msg := "assertion failed"
		if len(values) > 1 {
			msg = toString(values[1])
		}
		panic(msg)
	}
	return types.NewNull()
}

func Throw(values ...any) any {
	msg := "error"
	if len(values) > 0 {
		msg = toString(values[0])
	}
	panic(msg)
}
