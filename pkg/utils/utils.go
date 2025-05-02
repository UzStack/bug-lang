package utils

import (
	"fmt"
	"strconv"
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
