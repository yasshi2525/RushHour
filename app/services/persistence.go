package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var (
	db      *gorm.DB
	logMode = false
)

// InitPersistence prepares database connection and migrate
func InitPersistence() {
	revel.AppLog.Info("start init for persistence")
	defer revel.AppLog.Info("end init for persistence")

	db = connectDB()
	configureDB(db)
	migrateDB(db)
}

// TerminatePersistence defines the end task before application shutdown
func TerminatePersistence() {
	closeDB()
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

func configureDB(database *gorm.DB) {
	database.LogMode(logMode)
}

func migrateDB(database *gorm.DB) {
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

// Restore get model from database
func Restore() {
	revel.AppLog.Info("start restore from database")
	defer revel.AppLog.Info("end restore from database")

	start := time.Now()
	defer WarnLongExec(start, 5, "DBリストア", true)

	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	setNextID()
	fetchStatic()
	resolveStatic()
	genDynamics()
}

// setNextID set max id as NextID from database for Restore()
func setNextID() {
	for _, key := range Meta.StaticList {
		var maxID struct {
			V uint64
		}
		sql := fmt.Sprintf("SELECT max(id) as v FROM %s", key.Table())
		if err := db.Raw(sql).Scan(&maxID).Error; err == nil {
			Static.NextIDs[key] = &maxID.V
			revel.AppLog.Debugf("set NextID[%s] = %d", key, *Static.NextIDs[key])
		} else {
			panic(err)
		}
	}
}

// fetchStatic selects records for Restore()
func fetchStatic() {
	for _, key := range Meta.StaticList {
		// select文組み立て
		if rows, err := db.Table(key.Table()).Where("deleted_at is null").Rows(); err == nil {
			for rows.Next() {
				// 対応する Struct を作成
				base := key.Obj()
				if err := db.ScanRows(rows, base); err == nil {
					// Static に登録
					if obj, ok := base.(entities.Indexable); ok {
						Meta.StaticMap[key].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.ValueOf(obj))
						//revel.AppLog.Debugf("set Static[%s][%d] = %+v", key, obj.Idx(), obj)
					} else {
						panic(fmt.Errorf("invalid type %T: %+v", base, base))
					}
				} else {
					panic(err)
				}
			}
		} else {
			panic(err)
		}

	}
}

// resolveStatic set pointer from id for Restore()
func resolveStatic() {
	for _, rn := range Static.RailNodes {
		rn.Resolve(Static.Players[rn.OwnerID])
	}
	for _, re := range Static.RailEdges {
		re.Resolve(Static.RailNodes[re.FromID], Static.RailNodes[re.ToID])
	}
	for _, st := range Static.Stations {
		st.Resolve(Static.Players[st.OwnerID])
	}
	for _, g := range Static.Gates {
		g.Resolve(Static.Stations[g.InStationID])
	}
	for _, p := range Static.Platforms {
		p.Resolve(Static.RailNodes[p.OnRailNodeID], Static.Stations[p.InStationID])
	}
	for _, l := range Static.RailLines {
		l.Resolve(Static.Players[l.OwnerID])
	}
	for _, lt := range Static.LineTasks {
		lt.Resolve(Static.RailLines[lt.RailLineID])
		// nullable fields
		if lt.NextID != 0 {
			lt.Resolve(Static.LineTasks[lt.NextID])
		}
		if lt.StayID != 0 {
			lt.Resolve(Static.Platforms[lt.StayID])
		}
		if lt.MovingID != 0 {
			lt.Resolve(Static.RailEdges[lt.MovingID])
		}
	}
	for _, t := range Static.Trains {
		t.Resolve(Static.LineTasks[t.TaskID])
	}
	for _, h := range Static.Humans {
		h.Resolve(Static.Residences[h.FromID], Static.Companies[h.ToID])
		// nullable fields
		if h.OnPlatformID != 0 {
			h.Resolve(Static.Platforms[h.OnPlatformID])
		}
		if h.OnTrainID != 0 {
			h.Resolve(Static.Platforms[h.OnTrainID])
		}
	}
}

// genDynamics create Dynamic instances
func genDynamics() {
	walk, train := Config.Human.Weight, Config.Train.Weight
	for _, r := range Static.Residences {
		// R -> C, G
		GenStepResidence(r)
	}
	for _, c := range Static.Companies {
		// G -> C
		for _, g := range Static.Gates {
			GenStep(g, c, walk)
		}
	}
	for _, p := range Static.Platforms {
		// G <-> P
		g := p.InStation.Gate
		GenStep(p, g, walk)
		GenStep(g, p, walk)

		// P <-> P
		for _, p2 := range Static.Platforms {
			if p != p2 {
				GenStep(p, p2, train)
			}
		}
	}
	for _, h := range Static.Humans {
		GenStepHuman(h)
		Dynamic.Agents[h.ID] = entities.NewAgent(h)
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
