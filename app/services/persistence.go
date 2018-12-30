package services

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var db *gorm.DB

// InitPersistence prepares database connection and migrate
func InitPersistence() {

	db = connectDB()
	configureDB(db)
	migrateDB(db)
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

func configureDB(database *gorm.DB) *gorm.DB {
	database.LogMode(true)
	db.SingularTable(true)
	return database
}

func migrateDB(database *gorm.DB) *gorm.DB {
	db.AutoMigrate(
		&entities.Player{},
		&entities.Residence{},
		&entities.Company{},
		&entities.RailNode{},
		&entities.RailEdge{},
		&entities.Platform{},
		&entities.Gate{},
		&entities.Station{},
		&entities.LineTask{},
		&entities.Line{},
		&entities.Step{},
		&entities.Train{},
		&entities.Human{},
	)

	return db
}

// TerminatePersistence defines the end task before application shutdown
func TerminatePersistence() {
	closeDB()
}

func closeDB() {
	if err := db.Close(); err != nil {
		revel.AppLog.Error("failed to close the database", "error", err)
	}
	revel.AppLog.Info("disconnect database successfully")
}

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

	db.Begin()

	db.Commit()
}
