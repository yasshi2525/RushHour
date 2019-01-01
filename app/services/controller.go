package services

import (
	"time"
)

// ViewMap immitates user requests view
func ViewMap() interface{} {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	return Repo.Static
}
