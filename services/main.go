package services

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
	"github.com/yasshi2525/RushHour/entities"
)

var db *gorm.DB

var isInOperation bool

var conf *config.Config
var auther *auth.Auther

// Init prepares for starting game
func Init(c *config.Config, a *auth.Auther) {
	log.Println("start preparation for game")
	defer log.Println("end preparation for game")

	conf, auther = c, a

	start := time.Now()

	InitLock()
	defer WarnLongExec(start, start, conf.Game.Service.Perf.Init.D, "initialization", true)
	InitRepository()
	if conf.Game.Service.Backup.Enabled {
		db = connectDB()
		//db.LogMode(true)
		MigrateDB()
		Restore(true)
	}
	CreateIfAdmin()
	StartRouting()
}

// Purge deletes all user data
func Purge(o *entities.Player) error {
	if IsInOperation() {
		return fmt.Errorf("couldn't purge during under operation")
	}
	log.Println("start purging user data")
	defer log.Println("end purging user data")
	if conf.Game.Service.Backup.Enabled {
		PurgeDB(o)
	}
	InitRepository()
	if conf.Game.Service.Backup.Enabled {
		Restore(false)
	}
	CreateIfAdmin()
	StartRouting()
	return nil
}

// Terminate finalizes after stopping game
func Terminate() {
	if db != nil {
		closeDB()
	}
}

// Start start game
func Start() {
	log.Println("start starting game procedure")
	defer log.Println("end starting game procedure")
	if conf.Game.Service.Backup.Enabled {
		StartBackupTicker()
	}
	StartProcedure()
	isInOperation = true
	if conf.Game.Service.Procedure.Simulation {
		StartModelWatching()
		StartSimulation()
	}
}

// Stop stop game
func Stop() {
	log.Println("start stopping game procedure")
	defer log.Println("end stopping game procedure")
	if conf.Game.Service.Procedure.Simulation {
		StopSimulation()
		StopModelWatching()
	}
	isInOperation = false
	CancelRouting()
	StopProcedure()
	if conf.Game.Service.Backup.Enabled {
		StopBackupTicker()
		Backup(false)
	}
}

func connectDB() *gorm.DB {
	var (
		database *gorm.DB
		err      error
	)
	driver := conf.Secret.DB.Driver
	spec := conf.Secret.DB.Spec

	for i := 1; i <= 60; i++ {
		database, err = gorm.Open(driver, spec)
		if err != nil {
			log.Printf("failed to connect database(%v). retry after 10 seconds.", err)
			time.Sleep(10 * time.Second)
		}
	}

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}

	log.Println("connect database successfully")
	return database
}

func closeDB() {
	if err := db.Close(); err != nil {
		panic(err)
	}
	log.Println("disconnect database successfully")
}

// IsInOperation returns true after game started.
func IsInOperation() bool {
	return isInOperation
}
