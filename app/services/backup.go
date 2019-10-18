package services

import (
	"fmt"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var backupTicker *time.Ticker

// StartBackupTicker start tikcer
func StartBackupTicker() {
	backupTicker = time.NewTicker(Const.Backup.Interval.D)

	go watchBackup()
	revel.AppLog.Info("backup ticker was successfully started.")
}

// StopBackupTicker stop ticker
func StopBackupTicker() {
	if backupTicker != nil {
		backupTicker.Stop()
		revel.AppLog.Info("backup ticker was successfully stopped.")
	}
}

func watchBackup() {
	for range backupTicker.C {
		Backup(true)
	}
}

// Backup set model to database
func Backup(withLock bool) {
	start := time.Now()
	if withLock {
		MuModel.RLock()
		defer MuModel.RUnlock()
	}
	lock := time.Now()

	tx := db.Begin()
	new, up, del, skip := persistStatic(tx)
	log := logOperation(tx)
	tx.Commit()

	WarnLongExec(start, lock, Const.Perf.Backup.D, "backup")
	revel.AppLog.Infof("backup was successfully ended (new %d, up %d, del %d, skip %d, log %d)", new, up, del, skip, log)
}

func persistStatic(tx *gorm.DB) (int, int, int, int) {
	// upsert
	var createCnt, updateCnt, skipCnt int
	for _, res := range entities.TypeList {
		if res.IsDB() {
			Model.ForEach(res, func(raw entities.Entity) {
				obj := raw.(entities.Persistable)
				if obj.P().IsNew() {
					db.Create(obj)
					obj.P().Reset()
					createCnt++
				} else if obj.P().IsChanged() {
					tx.Save(obj)
					obj.P().Reset()
					updateCnt++
				} else {
					skipCnt++
				}
			})
		}
	}

	// remove old resources
	removeCnt := 0
	for i := len(entities.TypeList) - 1; i >= 0; i-- {
		key := entities.TypeList[i]
		if key.IsDB() {
			for _, id := range Model.Deletes[key] {
				sql := fmt.Sprintf("UPDATE %s SET updated_at = ?, deleted_at = ? WHERE id = ?", key.Table())
				tx.Exec(sql, time.Now(), time.Now(), id)
				removeCnt++
			}
			Model.Deletes[key] = Model.Deletes[key][:0]
		}
	}
	return createCnt, updateCnt, removeCnt, skipCnt
}

func logOperation(tx *gorm.DB) int {
	logCnt := 0
	for _, op := range OpCache {
		tx.Create(op)
		logCnt++
	}
	OpCache = OpCache[:0]
	return logCnt
}
