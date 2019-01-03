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
	backupTicker = time.NewTicker(Config.Backup.Interval.D)

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
	revel.AppLog.Info("start backup")
	defer revel.AppLog.Info("end backup")

	start := time.Now()
	defer WarnLongExec(start, Config.Perf.Backup.D, "backup")

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	updateForeignKey()

	tx := db.Begin()

	persistStatic(tx)

	tx.Commit()
}

func updateForeignKey() {
	// set mutable id from reference
	for _, res := range []entities.ModelType{
		entities.LINETASK, entities.TRAIN, entities.HUMAN} {
		mapdata := Meta.Map[res]
		for _, key := range mapdata.MapKeys() {
			obj := mapdata.MapIndex(key).Interface()
			obj.(entities.Resolvable).ResolveRef()
		}
	}
}

func persistStatic(tx *gorm.DB) {
	// upsert
	for _, res := range Meta.List {
		if res.IsDB() {
			ForeachModel(res, func(obj interface{}) {
				tx.Save(obj)
			})
		}
	}

	// remove old resources
	for i := len(Meta.List) - 1; i >= 0; i-- {
		key := Meta.List[i]
		if key.IsDB() {
			for _, id := range Model.Remove[key] {
				sql := fmt.Sprintf("UPDATE %s SET updated_at = ?, deleted_at = ? WHERE id = ?", key.Table())
				tx.Exec(sql, time.Now(), time.Now(), id)
				revel.AppLog.Debugf("delete %v(%d)", key, id)
			}
			Model.Remove[key] = Model.Remove[key][:0]
		}
	}
}
