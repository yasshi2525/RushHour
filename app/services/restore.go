package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yasshi2525/RushHour/app/services/route"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// ZERO means nil
const ZERO = 0

// Restore get model from database
func Restore() {
	revel.AppLog.Info("start restore from database")
	defer revel.AppLog.Info("end restore from database")
	start := time.Now()
	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()
	lock := time.Now()
	defer WarnLongExec(start, lock, Const.Perf.Restore.D, "restore", true)

	setNextID()
	fetchStatic()
	resolveStatic()
	for _, l := range Model.RailLines {
		lineValidation(l) // [DEBUG]
	}
	genDynamics()
}

// setNextID set max id as NextID from database for Restore()
func setNextID() {
	for _, key := range entities.TypeList {
		if !key.IsDB() {
			continue
		}
		var maxID struct {
			V uint64
		}
		sql := fmt.Sprintf("SELECT max(id) as v FROM %s", key.Table())
		if err := db.Raw(sql).Scan(&maxID).Error; err == nil {
			Model.NextIDs[key] = &maxID.V
		} else {
			panic(err)
		}
	}
}

// fetchStatic selects records for Restore()
func fetchStatic() {
	var cnt int
	for _, key := range entities.TypeList {
		if !key.IsDB() {
			continue
		}
		// select文組み立て
		if rows, err := db.Table(key.Table()).Where("deleted_at is null").Rows(); err == nil {
			for rows.Next() {
				// 対応する Struct を作成
				base := key.Obj(Model)
				if err := db.ScanRows(rows, base); err == nil {
					if obj, ok := base.(entities.Persistable); ok {
						obj.Reset() // DBNew -> DBMerged
					} else {
						panic(fmt.Errorf("invalid type %T: %+v", base, base))
					}

					// Model に登録
					if obj, ok := base.(entities.Indexable); ok {
						Model.Values[key].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.ValueOf(obj))
						cnt++
						//revel.AppLog.Debugf("set Model[%s][%d] = %v", key, obj.Idx(), obj)
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
	revel.AppLog.Infof("restored %d entities", cnt)
}

// resolveStatic set pointer from id for Restore()
func resolveStatic() {
	for _, key := range entities.TypeList {
		if !key.IsDB() {
			continue
		}
		Model.ForEach(key, func(obj entities.Indexable) {
			obj.(entities.Migratable).UnMarshal()
		})
	}
}

// genDynamics create Dynamic instances
func genDynamics() {
	for _, o := range Model.Players {
		route.RefreshTracks(o, Const.Routing.Worker)
	}
	for _, r := range Model.Residences {
		r.GenOutSteps()
	}
	for _, g := range Model.Gates {
		g.GenOutSteps()
	}
	for _, p := range Model.Platforms {
		p.GenOutSteps()
	}
	for _, l := range Model.RailLines {
		route.RefreshTransports(l, Const.Routing.Worker)
	}
	for _, h := range Model.Humans {
		h.GenOutSteps()
		Model.Agents[h.ID] = Model.NewAgent(h)
	}
}
