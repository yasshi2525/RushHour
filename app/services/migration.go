package services

import (
	"fmt"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// MigrateDB migrate database by reference Meta
func MigrateDB() {
	foreign := make(map[entities.ModelType]string)

	// create instance corresponding to each record
	for _, key := range Meta.List {
		if !key.IsDB() {
			continue
		}

		proto := key.Obj()
		db.AutoMigrate(proto)

		revel.AppLog.Debugf("migrated for %T", proto)

		// foreign key for owner
		if _, ok := proto.(entities.Ownable); ok {
			owner := fmt.Sprintf("%s(id)", entities.PLAYER.Table())
			db.Model(proto).AddForeignKey("owner_id", owner, "RESTRICT", "RESTRICT")

			//revel.AppLog.Debugf("added owner foreign key for %s table", owner)
		}

		foreign[key] = fmt.Sprintf("%s(id)", key.Table())
	}

	// RailEdge connects RailNode
	db.Model(entities.RAILEDGE.Obj()).AddForeignKey("from_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")
	db.Model(entities.RAILEDGE.Obj()).AddForeignKey("to_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")

	// Station composes Platforms and Gates
	db.Model(entities.PLATFORM.Obj()).AddForeignKey("rail_node_id", foreign[entities.RAILNODE], "RESTRICT", "RESTRICT")
	db.Model(entities.PLATFORM.Obj()).AddForeignKey("station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")
	db.Model(entities.GATE.Obj()).AddForeignKey("station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")

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
	db.Model(entities.HUMAN.Obj()).AddForeignKey("platform_id", foreign[entities.PLATFORM], "RESTRICT", "RESTRICT")
	db.Model(entities.HUMAN.Obj()).AddForeignKey("train_id", foreign[entities.TRAIN], "RESTRICT", "RESTRICT")
}
