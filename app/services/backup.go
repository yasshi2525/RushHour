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
		Backup()
	}
}

// Backup set model to database
func Backup() {
	start := time.Now()

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	updateForeignKey()

	tx := db.Begin()
	new, up, del, skip := persistStatic(tx)
	tx.Commit()

	WarnLongExec(start, Const.Perf.Backup.D, "backup")
	revel.AppLog.Infof("backup was successfully ended (new %d, up %d, del %d, skip %d)", new, up, del, skip)
}

func updateForeignKey() {
	// set mutable id from reference
	for _, res := range []entities.ModelType{
		entities.LINETASK, entities.TRAIN, entities.HUMAN} {
		Model.ForEach(res, func(obj entities.Indexable) {
			obj.(entities.Resolvable).ResolveRef()
		})
	}
}

func persistStatic(tx *gorm.DB) (int, int, int, int) {
	// upsert
	var createCnt, updateCnt, skipCnt int
	for _, res := range entities.TypeList {
		if res.IsDB() {
			Model.ForEach(res, func(raw entities.Indexable) {
				obj := raw.(entities.Persistable)
				if obj.IsNew() {
					db.Create(obj)
					obj.Reset()
					createCnt++
				} else if obj.IsChanged() {
					tx.Save(obj)
					obj.Reset()
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
			for _, id := range Model.Remove[key] {
				sql := fmt.Sprintf("UPDATE %s SET updated_at = ?, deleted_at = ? WHERE id = ?", key.Table())
				tx.Exec(sql, time.Now(), time.Now(), id)
				removeCnt++
			}
			Model.Remove[key] = Model.Remove[key][:0]
		}
	}
	return createCnt, updateCnt, removeCnt, skipCnt
}
