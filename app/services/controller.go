package services

import (
	"time"

	"github.com/yasshi2525/RushHour/app/models"
)

func ViewMap() {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	models.MuStatic.RLock()
	defer models.MuStatic.RUnlock()
}

func ChangeMap(msg string) {
	UpdateModel(msg)
}
