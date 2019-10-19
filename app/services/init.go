package services

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/yasshi2525/RushHour/app/services/auth"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var db *gorm.DB

var isInOperation bool

// Init prepares for starting game
func Init() {
	revel.AppLog.Info("start preparation for game")
	defer revel.AppLog.Info("end preparation for game")

	start := time.Now()

	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	InitLock()
	LoadSecret()
	LoadConf()
	auth.Init(Secret.Auth)
	defer WarnLongExec(start, start, Const.Perf.Init.D, "initialization", true)
	InitRepository()
	db = connectDB()
	//db.LogMode(true)
	MigrateDB()
	Restore(true)
	CreateIfAdmin()
	StartRouting()
}

// Purge deletes all user data
func Purge() error {
	if IsInOperation() {
		return fmt.Errorf("couldn't purge during under operation")
	}
	PurgeDB()
	InitRepository()
	Restore(false)
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
	StartBackupTicker()
	StartModelWatching()
	StartProcedure()
	isInOperation = true
}

// Stop stop game
func Stop() {
	isInOperation = false
	CancelRouting()
	StopProcedure()
	StopModelWatching()
	StopBackupTicker()
	Backup(false)
}

func connectDB() *gorm.DB {
	var (
		database *gorm.DB
		err      error
	)
	driver := getConfig("db.driver")
	spec := getConfig("db.spec")

	for i := 1; i <= 60; i++ {
		database, err = gorm.Open(driver, spec)
		if err != nil {
			revel.AppLog.Warnf("failed to connect database(%v). retry after 10 seconds.", err)
			time.Sleep(10 * time.Second)
		}
	}

	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
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

func getConfig(key string) string {
	if value, found := revel.Config.String(key); found {
		return value
	}
	panic(fmt.Errorf("%s is not defined", key))
}

// IsInOperation returns true after game started.
func IsInOperation() bool {
	return isInOperation
}
