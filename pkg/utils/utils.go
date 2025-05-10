package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func InArray(value any, arr []any) bool {
	for _, target := range arr {
		if target == value {
			return true
		}
	}
	return false
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
