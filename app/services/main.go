package services

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
	"github.com/yasshi2525/RushHour/app/entities"
)

var db *gorm.DB

var isInOperation bool

// ServiceConfig represents service settings
type ServiceConfig struct {
	// IsPersist is whether storing user data to database
	IsPersist bool
	// AppConf is whether using conf file
	AppConf *config.Config
	// Auther has encrypt keys
	Auther *auth.Auther
}

var serviceConf *ServiceConfig

// Init prepares for starting game
func Init(sc *ServiceConfig) {
	log.Println("start preparation for game")
	defer log.Println("end preparation for game")

	serviceConf = sc

	start := time.Now()

	InitLock()
	defer WarnLongExec(start, start, serviceConf.AppConf.Game.Service.Perf.Init.D, "initialization", true)
	InitRepository()
	if serviceConf.IsPersist {
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
	if serviceConf.IsPersist {
		PurgeDB(o)
	}
	InitRepository()
	if serviceConf.IsPersist {
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
	if serviceConf.IsPersist {
		StartBackupTicker()
	}
	StartProcedure()
	isInOperation = true
	if serviceConf.AppConf.Game.Service.Procedure.Simulation {
		StartModelWatching()
		StartSimulation()
	}
}

// Stop stop game
func Stop() {
	log.Println("start stopping game procedure")
	defer log.Println("end stopping game procedure")
	if serviceConf.AppConf.Game.Service.Procedure.Simulation {
		StopSimulation()
		StopModelWatching()
	}
	isInOperation = false
	CancelRouting()
	StopProcedure()
	if serviceConf.IsPersist {
		StopBackupTicker()
		Backup(false)
	}
}

func connectDB() *gorm.DB {
	var (
		database *gorm.DB
		err      error
	)
	driver := serviceConf.AppConf.Secret.DB.Driver
	spec := serviceConf.AppConf.Secret.DB.Spec

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
