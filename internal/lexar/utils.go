package lexar

import (
	"regexp"
)

func IsAlpha(value string) bool {
	re := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return re.MatchString(value)
}

func IsNumber(value string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(value)
}
