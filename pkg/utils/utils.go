package utils

func InArray(value any, arr []any) bool {
	for _, target := range arr {
		if target == value {
			return true
		}
	}
	return false
}
