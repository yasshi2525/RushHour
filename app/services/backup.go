package services

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var backupTicker *time.Ticker

// StartBackupTicker start tikcer
func StartBackupTicker() {
	backupTicker = time.NewTicker(Config.Backup.Interval.Duration)

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
	defer WarnLongExec(start, 2, "バックアップ", true)

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
	// set id from reference
	for _, lt := range Static.LineTasks {
		lt.ResolveRef()
	}
	for _, t := range Static.Trains {
		t.ResolveRef()
	}
	for _, h := range Static.Humans {
		h.ResolveRef()
	}
}

func persistStatic(tx *gorm.DB) {
	// upsert
	for _, res := range Meta.StaticList {
		mapdata := Meta.StaticMap[res]
		for _, key := range mapdata.MapKeys() {
			obj := mapdata.MapIndex(key).Interface()
			tx.Save(obj)
			//revel.AppLog.Debugf("persist %T(%d): %+v", obj, key.Uint(), obj)
		}
	}

	// remove old resources
	for _, key := range Meta.StaticList {
		for _, id := range Static.WillRemove[key] {
			sql := fmt.Sprintf("UPDATE %s SET updated_at = ?, deleted_at = ? WHERE id = ?", key.Table())
			tx.Exec(sql, time.Now(), time.Now(), id)
			revel.AppLog.Debugf("delete %s(%d)", key.Table(), id)
		}
		Static.WillRemove[key] = Static.WillRemove[key][:0]
	}
}
