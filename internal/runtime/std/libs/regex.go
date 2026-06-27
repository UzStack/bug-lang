package libs

import (
	"regexp"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func RegexMatch(pattern *types.StringValue, s *types.StringValue) any {
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		panic("regex.match: " + err.Error())
	}
	return types.NewBool(re.MatchString(s.Value))
}

func RegexFind(pattern *types.StringValue, s *types.StringValue) any {
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		panic("regex.find: " + err.Error())
	}
	matches := re.FindAllString(s.Value, -1)
	values := make([]any, len(matches))
	for i, m := range matches {
		values[i] = types.NewString(m)
	}
	return types.NewArray(values)
}

func RegexReplace(pattern *types.StringValue, s *types.StringValue, repl *types.StringValue) any {
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		panic("regex.replace: " + err.Error())
	}
	return types.NewString(re.ReplaceAllString(s.Value, repl.Value))
}
