package services

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var (
	db       *gorm.DB
	excludes *[]string
)

type eachCallback func(v reflect.Value)

// InitPersistence prepares database connection and migrate
func InitPersistence() {
	excludes = &[]string{"Steps"}

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
		&entities.Station{},
		&entities.Platform{},
		&entities.Gate{},
		&entities.LineTask{},
		&entities.Line{},
		//&entities.Step{}, Step is out of target for persistence because it can intruduce by other resources
		&entities.Train{},
		&entities.Human{},
	)

	// Player has private resources
	for _, t := range []interface{}{
		&entities.RailNode{},
		&entities.RailEdge{},
		&entities.Station{},
		&entities.Platform{},
		&entities.Gate{},
		&entities.LineTask{},
		&entities.Line{},
		&entities.Train{},
	} {
		db.Model(t).AddForeignKey("owner_id", "player(id)", "RESTRICT", "RESTRICT")
	}

	// RailEdge connects RailNode
	db.Model(&entities.RailEdge{}).AddForeignKey("from_id", "rail_node(id)", "CASCADE", "RESTRICT")
	db.Model(&entities.RailEdge{}).AddForeignKey("to_id", "rail_node(id)", "CASCADE", "RESTRICT")

	// Station composes Platforms and Gates
	db.Model(&entities.Platform{}).AddForeignKey("on_id", "rail_node(id)", "RESTRICT", "RESTRICT")
	db.Model(&entities.Platform{}).AddForeignKey("in_id", "station(id)", "CASCADE", "RESTRICT")
	db.Model(&entities.Gate{}).AddForeignKey("in_id", "station(id)", "CASCADE", "RESTRICT")

	// Line composes LineTasks
	db.Model(&entities.LineTask{}).AddForeignKey("line_id", "line(id)", "CASCADE", "RESTRICT")
	// LineTask is chainable
	db.Model(&entities.LineTask{}).AddForeignKey("next_id", "line_task(id)", "SET NULL", "RESTRICT")

	// Train runs on a chain of Line
	db.Model(&entities.Train{}).AddForeignKey("task_id", "line_task(id)", "RESTRICT", "RESTRICT")

	// Human departs from Residence and destinates to Company
	db.Model(&entities.Human{}).AddForeignKey("from_id", "residence(id)", "RESTRICT", "RESTRICT")
	db.Model(&entities.Human{}).AddForeignKey("to_id", "company(id)", "RESTRICT", "RESTRICT")
	// Human is sometimes on Platform or on Train
	db.Model(&entities.Human{}).AddForeignKey("on_platform_id", "platform(id)", "RESTRICT", "RESTRICT")
	db.Model(&entities.Human{}).AddForeignKey("on_train_id", "train(id)", "RESTRICT", "RESTRICT")

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

func eachEntity(skips *[]string, callback eachCallback) {
	rt, rv := reflect.TypeOf(Static), reflect.ValueOf(Static)
	for i := 0; i < rt.NumField(); i++ {
		if f := rv.Field(i); f.Kind() == reflect.Map {
			// skip specific field
			shouldSkip := false
			for _, skip := range *skips {
				if strings.Compare(rt.Field(i).Name, skip) == 0 {
					shouldSkip = true
					break
				}
			}
			if !shouldSkip {
				for _, e := range f.MapKeys() {
					callback(f.MapIndex(e))
				}
			}
		} else {
			revel.AppLog.Warnf("%s is not map", f.Kind().String())
		}
	}
}

// Backup set model to database
func Backup() {
	revel.AppLog.Info("バックアップ 開始")
	defer revel.AppLog.Info("バックアップ 終了")

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	tx := db.Begin()

	// resolve refer
	eachEntity(excludes, func(val reflect.Value) {
		if v, ok := val.Interface().(entities.Resolvable); ok {
			v.ResolveRef()
		} else {
			revel.AppLog.Warnf("%s is not resolvable", val.String())
		}
	})

	// upsert
	eachEntity(excludes, func(val reflect.Value) {
		db.Save(val.Elem().Interface())
	})

	// remove old resources
	for _, resource := range EntityTypes {
		for _, id := range WillRemove[resource] {
			sql := fmt.Sprintf("UPDATE %s SET updated_at = ?, deleted_at = ? WHERE id = ?", resource)
			db.Exec(sql, time.Now(), time.Now(), id)
		}
	}

	tx.Commit()
}
