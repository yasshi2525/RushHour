package services

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var db *gorm.DB

// Init prepares for starting game
func Init() {
	revel.AppLog.Info("start preparation for game")
	defer revel.AppLog.Info("end preparation for game")

	start := time.Now()
	InitLock()
	LoadConf()
	defer WarnLongExec(start, start, Const.Perf.Init.D, "initialization", true)
	InitRepository()
	db = connectDB()
	//db.LogMode(true)
	MigrateDB()
	Restore()
	StartRouting()
}

// Terminate finalizes after stopping game
func Terminate() {
	if db != nil {
		closeDB()
	}
}

// Start start game
func Start() {
	StartBackupTicker()
	StartModelWatching()
	StartProcedure()
}

// Stop stop game
func Stop() {
	CancelRouting()
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

	for i := 1; i <= 12; i++ {
		database, err = gorm.Open(driver, spec)
		if err != nil {
			revel.AppLog.Warn("failed to connect database. retry after 5 seconds.")
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		panic("failed to connect database")
	}

	revel.AppLog.Info("connect database successfully")
	return database
}

func closeDB() {
	if err := db.Close(); err != nil {
		panic(err)
	}
	revel.AppLog.Info("disconnect database successfully")
}
