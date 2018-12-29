package services

import (
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
)

func ViewMap() {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	entities.MuStatic.RLock()
	defer entities.MuStatic.RUnlock()
}

func ChangeMap(msg string) {
	UpdateModel(msg)
}
