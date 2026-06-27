package libs

import (
	"time"

	"github.com/UzStack/bug-lang/internal/runtime/types"
)

func TimeNow() any {
	return types.NewString(time.Now().Format("2006-01-02 15:04:05"))
}

func TimeUnix() any {
	return types.NewInt(int(time.Now().Unix()))
}

// TimeFormat(layout) — Go layout (masalan "2006-01-02") bilan hozirgi vaqt.
func TimeFormat(layout *types.StringValue) any {
	return types.NewString(time.Now().Format(layout.Value))
}

func TimeSleep(ms *types.IntValue) any {
	time.Sleep(time.Duration(ms.Value) * time.Millisecond)
	return types.NewNull()
}
