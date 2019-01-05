package services

import (
	"fmt"
	"reflect"
	"time"

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
	defer WarnLongExec(start, Config.Perf.Restore.D, "restore", true)

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
	for _, key := range Meta.List {
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
	for _, key := range Meta.List {
		if !key.IsDB() {
			continue
		}
		// select文組み立て
		if rows, err := db.Table(key.Table()).Where("deleted_at is null").Rows(); err == nil {
			for rows.Next() {
				// 対応する Struct を作成
				base := key.Obj()
				if err := db.ScanRows(rows, base); err == nil {
					// Model に登録
					if obj, ok := base.(entities.Indexable); ok {
						Meta.Map[key].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.ValueOf(obj))
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
}

// resolveStatic set pointer from id for Restore()
func resolveStatic() {
	for _, rn := range Model.RailNodes {
		rn.Resolve(Model.Players[rn.OwnerID])
	}
	for _, re := range Model.RailEdges {
		re.Resolve(
			Model.RailNodes[re.FromID],
			Model.RailNodes[re.ToID],
			Model.RailEdges[re.ReverseID])
	}
	for _, st := range Model.Stations {
		st.Resolve(Model.Players[st.OwnerID])
	}
	for _, g := range Model.Gates {
		g.Resolve(Model.Stations[g.StationID])
	}
	for _, p := range Model.Platforms {
		p.Resolve(
			Model.RailNodes[p.RailNodeID],
			Model.Stations[p.StationID])
	}
	for _, l := range Model.RailLines {
		l.Resolve(Model.Players[l.OwnerID])
	}
	for _, lt := range Model.LineTasks {
		lt.Resolve(Model.RailLines[lt.RailLineID])
		// nullable fields
		if lt.NextID != ZERO {
			lt.Resolve(Model.LineTasks[lt.NextID])
		}
		if lt.StayID != ZERO {
			lt.Resolve(Model.Platforms[lt.StayID])
		}
		if lt.MovingID != ZERO {
			lt.Resolve(Model.RailEdges[lt.MovingID])
		}
	}
	for _, t := range Model.Trains {
		// nullable fields
		if t.TaskID != ZERO {
			t.Resolve(Model.LineTasks[t.TaskID])
		}
	}
	for _, h := range Model.Humans {
		h.Resolve(Model.Residences[h.FromID], Model.Companies[h.ToID])
		// nullable fields
		if h.PlatformID != ZERO {
			h.Resolve(Model.Platforms[h.PlatformID])
		}
		if h.TrainID != ZERO {
			h.Resolve(Model.Platforms[h.TrainID])
		}
	}
}

// genDynamics create Dynamic instances
func genDynamics() {
	for _, r := range Model.Residences {
		// R -> C, G
		GenStepResidence(r)
	}
	for _, c := range Model.Companies {
		// G -> C
		for _, g := range Model.Gates {
			GenStep(g, c)
		}
	}
	for _, p := range Model.Platforms {
		// G <-> P
		g := p.InStation.Gate
		GenStep(p, g)
		GenStep(g, p)
	}
	for _, l := range Model.RailLines {
		if l.IsRing() {
			genStepRailLine(l)
		}
	}
	for _, h := range Model.Humans {
		GenStepHuman(h)
		Model.Agents[h.ID] = entities.NewAgent(h)
	}
}
