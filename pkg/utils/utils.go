package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func InArray(value any, arr []any) bool {
	return slices.Contains(arr, value)
}

func Str2Int(value any) (int, error) {
	switch val := value.(type) {
	case int:
		return val, nil
	case string:
		return strconv.Atoi(val)
	}
	panic(fmt.Sprintf("Str2Int error: %s", value))
}
func Int2String(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case string:
		return v
	}
	panic(fmt.Sprintf("Int2String error: %s", value))
}

func Int2Float(value any) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case *types.IntValue:
		return float64(v.Value), nil
	case *types.FloatValue:
		return float64(v.Value), nil
	}
	return -1, errors.New("value not integer")
}

func Float2Int(value any) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case float32:
		return int(v), nil
	case *types.FloatValue:
		return int(v.Value), nil
	}
	return -1, errors.New("value not integer")
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func DecodeBug(data any) any {
	switch reflect.ValueOf(data).Kind() {
	case reflect.String:
		return types.NewString(data.(string))
	case reflect.Int:
		return types.NewInt(data.(int))
	case reflect.Int16:
		return types.NewInt(int(data.(int16)))
	case reflect.Int32:
		return types.NewInt(int(data.(int32)))
	case reflect.Int64:
		return types.NewInt(int(data.(int64)))
	case reflect.Bool:
		return types.NewBool(data.(bool))
	case reflect.Map:
		response := types.NewMap(make(map[string]any)).(*types.MapValue)
		for key, value := range data.(map[string]any) {
			response.Add(key, DecodeBug(value))
		}
		return response
	case reflect.Slice:
		var response []any
		for _, value := range data.([]any) {
			response = append(response, DecodeBug(value))
		}
		return response
	}
	return data
}

func EncodeBug(data any) any {
	switch v := data.(type) {
	case *types.ArrayValue:
		var response []any
		for _, value := range v.GetValue().([]any) {
			response = append(response, EncodeBug(value))
		}
		return response
	case *types.MapValue:
		response := make(map[string]any)
		for key, value := range v.GetValue().(map[string]any) {
			response[Int2String(EncodeBug(key))] = EncodeBug(value)
		}
		return response
	case types.Object:
		return v.GetValue()
	}
	return data
}
