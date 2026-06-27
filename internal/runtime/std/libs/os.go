package libs

import (
	"os"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func OsRead(path *types.StringValue) any {
	data, err := os.ReadFile(path.Value)
	if err != nil {
		panic("os.read: " + err.Error())
	}
	return types.NewString(string(data))
}

func OsWrite(path *types.StringValue, content *types.StringValue) any {
	if err := os.WriteFile(path.Value, []byte(content.Value), 0644); err != nil {
		panic("os.write: " + err.Error())
	}
	return types.NewNull()
}

func OsAppend(path *types.StringValue, content *types.StringValue) any {
	f, err := os.OpenFile(path.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("os.append: " + err.Error())
	}
	defer f.Close()
	if _, err := f.WriteString(content.Value); err != nil {
		panic("os.append: " + err.Error())
	}
	return types.NewNull()
}

func OsExists(path *types.StringValue) any {
	_, err := os.Stat(path.Value)
	return types.NewBool(err == nil)
}

func OsRemove(path *types.StringValue) any {
	if err := os.Remove(path.Value); err != nil {
		panic("os.remove: " + err.Error())
	}
	return types.NewNull()
}

func OsEnv(key *types.StringValue) any {
	return types.NewString(os.Getenv(key.Value))
}

func OsArgs() any {
	values := make([]any, len(os.Args))
	for i, a := range os.Args {
		values[i] = types.NewString(a)
	}
	return types.NewArray(values)
}
