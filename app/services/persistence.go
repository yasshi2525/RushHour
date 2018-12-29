package services

import (
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"
)

func Restore() {
	revel.AppLog.Info("DBリストア 開始")
	defer revel.AppLog.Info("DBリストア 終了")

	entities.MuStatic.Lock()
	defer entities.MuStatic.Unlock()

	entities.MuAgent.Lock()
	defer entities.MuAgent.Unlock()

	time.Sleep(1 * time.Second)
}

func Backup() {
	revel.AppLog.Info("バックアップ 開始")
	defer revel.AppLog.Info("バックアップ 終了")

	entities.MuStatic.RLock()
	defer entities.MuStatic.RUnlock()

	entities.MuAgent.RLock()
	defer entities.MuAgent.RUnlock()

	time.Sleep(10 * time.Second)
}
