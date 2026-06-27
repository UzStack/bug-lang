package lexar

import (
	"regexp"
)

// ponytail: regexlar bir marta compile qilinadi (lexer hot path)
var (
	alphaRe       = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	numberRe      = regexp.MustCompile(`^[0-9]+$`)
	signedFloatRe = regexp.MustCompile(`^[0-9]+(\.[0-9]*)?$`)
)

func IsAlpha(value string) bool {
	return alphaRe.MatchString(value)
}

func IsNumber(value string) bool {
	return numberRe.MatchString(value)
}

func IsSignedFloat(value string) bool {
	return signedFloatRe.MatchString(value)
}
