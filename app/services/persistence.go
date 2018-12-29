package services

import (
	"time"

	"github.com/revel/revel"
)

// Restore get model from database
func Restore() {
	revel.AppLog.Info("DBリストア 開始")
	defer revel.AppLog.Info("DBリストア 終了")

	MuStatic.Lock()
	defer MuStatic.Unlock()

	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	time.Sleep(1 * time.Second)
}

// Backup set model to database
func Backup() {
	revel.AppLog.Info("バックアップ 開始")
	defer revel.AppLog.Info("バックアップ 終了")

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	time.Sleep(10 * time.Second)
}
