package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// MigrateDB migrate database by reference Meta
// POLICY:
//   1. not allow nullable field because go cann't identify nil or zero value
//      so, no foreign key restriction for nullable field
//   2. id must be grater than 0. id: 0 means nil
//   3. not use sql.NullInst64 for performance
func MigrateDB() {
	foreign := make(map[entities.ModelType]string)

	db.AutoMigrate(&OpLog{})

	// create instance corresponding to each record
	for _, key := range entities.TypeList {
		if !key.IsDB() {
			continue
		}

		proto := key.Obj(Model)
		db.AutoMigrate(proto)

		// foreign key for owner
		if proto.B().O != nil {
			owner := fmt.Sprintf("%s(id)", entities.PLAYER.Table())
			db.Model(proto).AddForeignKey("owner_id", owner, "RESTRICT", "RESTRICT")
		}

		foreign[key] = fmt.Sprintf("%s(id)", key.Table())
	}

	// RailEdge connects RailNode
	db.Model(entities.RAILEDGE.Obj(Model)).AddForeignKey("from_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")
	db.Model(entities.RAILEDGE.Obj(Model)).AddForeignKey("to_id", foreign[entities.RAILNODE], "CASCADE", "RESTRICT")

	// Station composes Platforms and Gates
	db.Model(entities.PLATFORM.Obj(Model)).AddForeignKey("rail_node_id", foreign[entities.RAILNODE], "RESTRICT", "RESTRICT")
	db.Model(entities.PLATFORM.Obj(Model)).AddForeignKey("station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")
	db.Model(entities.GATE.Obj(Model)).AddForeignKey("station_id", foreign[entities.STATION], "CASCADE", "RESTRICT")

	// Line composes LineTasks
	db.Model(entities.LINETASK.Obj(Model)).AddForeignKey("rail_line_id", foreign[entities.RAILLINE], "CASCADE", "RESTRICT")

	// Human departs from Residence and destinates to Company
	db.Model(entities.HUMAN.Obj(Model)).AddForeignKey("from_id", foreign[entities.RESIDENCE], "RESTRICT", "RESTRICT")
	db.Model(entities.HUMAN.Obj(Model)).AddForeignKey("to_id", foreign[entities.COMPANY], "RESTRICT", "RESTRICT")
}
