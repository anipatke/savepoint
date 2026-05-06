package board

import (
	"fmt"
	"os"
	"sync/atomic"
)

var debugEnabled atomic.Bool

// SetDebug enables or disables debug logging for the board package.
func SetDebug(enabled bool) {
	debugEnabled.Store(enabled)
}

// DebugEnabled reports whether debug logging is active.
func DebugEnabled() bool {
	return debugEnabled.Load()
}

func debugf(format string, args ...any) {
	if !debugEnabled.Load() {
		return
	}
	fmt.Fprintf(os.Stderr, "[savepoint debug] "+format+"\n", args...)
}
