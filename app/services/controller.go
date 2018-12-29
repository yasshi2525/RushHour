package services

import (
	"time"
)

// ViewMap immitates user requests view
// TODO remove
func ViewMap() {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	MuStatic.RLock()
	defer MuStatic.RUnlock()
}
