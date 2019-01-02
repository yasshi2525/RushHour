package services

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var db *gorm.DB

// Init prepares for starting game
func Init() {
	revel.AppLog.Info("start preparation for game")
	defer revel.AppLog.Info("end preparation for game")

	start := time.Now()
	defer WarnLongExec(start, 10, "初期化", true)

	InitLock()
	LoadConf()
	InitRepository()
	db = connectDB()
	MigrateDB()
	Restore()
}

// Terminate finalizes after stopping game
func Terminate() {
	if db != nil {
		closeDB()
	}
}

// Start start game
func Start() {
	StartModelWatching()
	StartProcedure()
}

// Stop stop game
func Stop() {
	StopProcedure()
	StopModelWatching()
	StopBackupTicker()
	Backup()
}

func connectDB() *gorm.DB {
	var (
		database     *gorm.DB
		driver, spec string
		found        bool
		err          error
	)
	if driver, found = revel.Config.String("db.driver"); !found {
		panic("db.drvier is not defined")
	}
	if spec, found = revel.Config.String("db.spec"); !found {
		panic("db.spec is not defined")
	}

	if database, err = gorm.Open(driver, spec); err != nil {
		panic("failed to connect database")
	}

	revel.AppLog.Info("connect database successfully")
	return database
}

// MigrateDB migrate database by reference Meta
func MigrateDB() {
	foreign := make(map[entities.StaticRes]string)

	// create instance corresponding to each record
	for _, key := range Meta.StaticList {
		proto := key.Obj()
		db.AutoMigrate(proto)

		revel.AppLog.Debugf("migrated for %T", proto)

		// foreign key for owner
		if _, ok := proto.(entities.Ownable); ok {
			owner := fmt.Sprintf("%s(id)", entities.PLAYER.Table())
			db.Model(proto).AddForeignKey("owner_id", owner, "RESTRICT", "RESTRICT")

			revel.AppLog.Debugf("added owner foreign key for %s table", owner)
		}

		foreign[key] = fmt.Sprintf("%s(id)", key.Table())
	}

	// RailEdge connects RailNode
	db.Model(entities.RAILEDGE.Obj()).AddForeignKey("from_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")
	db.Model(entities.RAILEDGE.Obj()).AddForeignKey("to_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")

	// Station composes Platforms and Gates
	db.Model(entities.PLATFORM.Obj()).AddForeignKey("on_rail_node_id", foreign[entities.RAILNODE], "RESTRICT", "RESTRICT")
	db.Model(entities.PLATFORM.Obj()).AddForeignKey("in_station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")
	db.Model(entities.GATE.Obj()).AddForeignKey("in_station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")

	// Line composes LineTasks
	db.Model(entities.LINETASK.Obj()).AddForeignKey("rail_line_id", foreign[entities.LINE], "CASCADE", "RESTRICT")
	// LineTask is chainable
	db.Model(entities.LINETASK.Obj()).AddForeignKey("next_id", foreign[entities.LINETASK], "SET NULL", "RESTRICT")
	// LineTask is sometimes on rail or platform
	db.Model(entities.LINETASK.Obj()).AddForeignKey("moving_id", foreign[entities.RAILEDGE], "RESTRICT", "RESTRICT")
	db.Model(entities.LINETASK.Obj()).AddForeignKey("stay_id", foreign[entities.PLATFORM], "RESTRICT", "RESTRICT")

	// Train runs on a chain of Line
	db.Model(entities.TRAIN.Obj()).AddForeignKey("task_id", foreign[entities.LINETASK], "RESTRICT", "RESTRICT")

	// Human departs from Residence and destinates to Company
	db.Model(entities.HUMAN.Obj()).AddForeignKey("from_id", foreign[entities.RESIDENCE], "RESTRICT", "RESTRICT")
	db.Model(entities.HUMAN.Obj()).AddForeignKey("to_id", foreign[entities.COMPANY], "RESTRICT", "RESTRICT")
	// Human is sometimes on Platform or on Train
	db.Model(entities.HUMAN.Obj()).AddForeignKey("on_platform_id", foreign[entities.PLATFORM], "RESTRICT", "RESTRICT")
	db.Model(entities.HUMAN.Obj()).AddForeignKey("on_train_id", foreign[entities.TRAIN], "RESTRICT", "RESTRICT")
}

func closeDB() {
	if err := db.Close(); err != nil {
		panic(err)
	}
	revel.AppLog.Info("disconnect database successfully")
}
