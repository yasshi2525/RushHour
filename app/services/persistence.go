package services

import (
	"time"

	"github.com/yasshi2525/RushHour/app/models"

	"github.com/revel/revel"
)

func Restore() {
	revel.AppLog.Info("DBリストア 開始")
	defer revel.AppLog.Info("DBリストア 終了")

	models.MuStatic.Lock()
	defer models.MuStatic.Unlock()

	models.MuAgent.Lock()
	defer models.MuAgent.Unlock()

	time.Sleep(1 * time.Second)
}

func Backup() {
	revel.AppLog.Info("バックアップ 開始")
	defer revel.AppLog.Info("バックアップ 終了")

	models.MuStatic.RLock()
	defer models.MuStatic.RUnlock()

	models.MuAgent.RLock()
	defer models.MuAgent.RUnlock()

	time.Sleep(10 * time.Second)
}
